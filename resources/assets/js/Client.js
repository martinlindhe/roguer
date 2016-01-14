/**
 * Websocket client
 */
export class Client
{
    constructor(gameState)
    {
        var loc = window.location, wsEntry;
        if (loc.protocol === "https:") {
            wsEntry = "wss:";
        } else {
            wsEntry = "ws:";
        }
        wsEntry += "//" + loc.host;
        wsEntry += loc.pathname + "ws";

        this.gameState = gameState;
        this.socket = new WebSocket(wsEntry);

        this.sessionToken = window.sessionStorage.getItem('_token');

        var parent = this;

        /**
         * @param msg MessageEvent
         */
        this.socket.onmessage = function(msg)
        {
            var cmd = JSON.parse(msg.data);

            switch (cmd.Type) {
            case 'xy':
                parent.handleXyMessage(cmd);
                break;
            case 'move_res':
                parent.handleMoveResMessage(cmd);
                break;

            case 'ok':
                console.log("server OK: " + msg.data);
                break;

            case 'tick':
                parent.gameState.setServerTime(cmd.Time);
                break;

            case 'msg':
                //console.log(cmd);
                parent.gameState.messageToLog({time: cmd.Time, text: cmd.Message});
                break;

            default:
                console.log("<-recv- " + msg.data);
                console.log("unknown command from server: " + cmd.Type);
            }
        };

        this.socket.onopen = function()
        {
            console.log('Websocket connected');

            if (parent.sessionToken) {
                console.log("Resuming session");
                this.send("continue " + parent.sessionToken);
            } else {
                this.send("new_player " + parent.gameState.playerName);
            }
        };
    }

    sendCommand(cmd)
    {
        this.socket.send(cmd);
    }

    sendMove()
    {
        var newX = Math.floor((this.gameState.playerSprite.x / this.gameState.tileWidth) / this.gameState.worldScale.x);
        var newY = Math.floor((this.gameState.playerSprite.y / this.gameState.tileHeight) / this.gameState.worldScale.y);

        if (this.prevX == newX && this.prevY == newY) {
            // dont spam server when coords havent changed
            return;
        }

        this.sendCommand("move " + newX + " " + newY + " " + this.sessionToken);
        this.prevX = newX;
        this.prevY = newY;
    }

    handleXyMessage(cmd)
    {
        this.gameState.spawnPlayer(cmd);

        this.sessionToken = cmd.Token;
        // console.log("got token " + this.sessionToken);
        window.sessionStorage.setItem('_token', cmd.Token);
    }

    handleMoveResMessage(cmd)
    {
        // console.log("Rendering " + cmd.LocalSpawns.length + " spawns at " + cmd.X + ", " + cmd.Y);
        this.gameState.renderLocalSpawns(cmd.LocalSpawns);
    }
}
