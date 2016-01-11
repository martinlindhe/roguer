/**
 * Websocket client
 */
class Client
{
    constructor(gameState)
    {
        this.gameState = gameState;
        this.socket = new WebSocket('ws://localhost:3322/ws');

        this.sessionToken = window.sessionStorage.getItem('_token');
        if (this.sessionToken) {
            console.log("Re-used previous token " + this.sessionToken);
        }

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

            default:
                console.log("<-recv- " + msg.data);
                console.log("unknown command from server: " + cmd.Type);
            }
        };

        this.socket.onopen = function()
        {
            console.log('Websocket connected');
            this.send("new_player " + parent.gameState.playerName);
        };
    }

    sendMove()
    {
        var newX = Math.floor(this.gameState.playerGroup.x / this.gameState.tileWidth);
        var newY = Math.floor(this.gameState.playerGroup.y / this.gameState.tileHeight);

        if (this.prevX == newX && this.prevY == newY) {
            // dont spam server when coords havent changed
            return;
        }

        this.socket.send("move " + newX + " " + newY + " " + this.sessionToken);
        this.prevX = newX;
        this.prevY = newY;
    }

    handleXyMessage(cmd)
    {
        this.gameState.spawnPlayer(cmd);

        this.sessionToken = cmd.Token;
        window.sessionStorage.setItem('_token', cmd.Token);
    }

    handleMoveResMessage(cmd)
    {
        // console.log("Rendering " + cmd.LocalSpawns.length + " spawns at " + cmd.X + ", " + cmd.Y);
        this.gameState.renderLocalSpawns(cmd.LocalSpawns);
    }
}

export default Client;
