
class GameState extends Phaser.State
{
    preload()
    {
        this.game.time.advancedTiming = true; // required for fps counter

        this.game.stage.backgroundColor = '#262f71';  // deep water

        this.game.load.tilemap('islandMap', '/island/full', null, Phaser.Tilemap.TILED_JSON);
        this.game.load.image('ground', 'img/tileset/oddball/ground.png', 4, 8);

        this.game.load.image('minimap', 'img/islands/current.png');

        this.game.load.atlas('characterAtlas', 'img/tileset/oddball/characters.png', 'sprite/character');
        this.game.load.atlas('itemAtlas', 'img/tileset/oddball/items.png', 'sprite/item');
        this.game.load.atlas('ground2Atlas', 'img/tileset/oddball/ground2.png', 'sprite/ground2');

        this.game.load.spritesheet('button', 'img/tileset/ui/buttons.png', 27, 24);

        this.game.load.image('oddballFont', 'img/tileset/oddball/font.png');

        this.game.load.audio('bgSound', ['audio/dead_feelings.mp3']);
    }

    create()
    {
        this.worldScale = 1.0;
        this.tileWidth = 8;
        this.tileHeight = 4;
        this.playerName = "Jimpson";

        this.logMessages = [
            {text: "hello there", time: 1222},
            {text: "later on", time: 1234},
        ];


        // scale to whole window
        this.game.scale.setGameSize(window.innerWidth, window.innerHeight);



        // world (except UI) is in this group, so it can be scaled
        this.stageGroup = this.game.add.group();


        this.groundMap = this.game.add.tilemap('islandMap');
        this.groundMap.addTilesetImage('island_tiles', 'ground');
        this.groundLayer = this.groundMap.createLayer(0);
        this.groundLayer.resizeWorld();

        this.stageGroup.add(this.groundLayer);





        this.music = this.game.add.audio('bgSound');
        this.music.volume = 0.20; // 20%


        this.spawnLayer = this.game.add.group();
        this.spawnLayer.z = 5;
        this.stageGroup.add(this.spawnLayer);


        this.cursors = this.game.input.keyboard.createCursorKeys();


        var minimapScale = 3;
        var minimapX = this.game.width - this.game.cache.getImage('minimap').width / minimapScale;
        this.minimap = this.game.add.sprite(minimapX, 0, 'minimap');
        this.minimap.fixedToCamera = true;
        this.minimap.scale.set(1.0 / minimapScale);
        this.minimap.alpha = 0.8;
        this.minimap.setScaleMinMax(1.0 / minimapScale, 1.0 / minimapScale);



        var button = this.game.add.button(
            this.game.width - 102,
            2,
            'button',
            function() { // onClick
                if (this.music.isPlaying) {
                    this.music.stop();
                } else {
                    this.music.play();
                }
            },
            this,
            0,
            0,
            0
        );

        button.fixedToCamera = true;




        // fog of war

        // The radius of the circle of light
        this.LIGHT_RADIUS = 100;

        // Create the shadow texture
        this.shadowTexture = this.game.add.bitmapData(this.game.width, this.game.height);

        // Create an object that will use the bitmap as a texture
        var lightSprite = this.game.add.image(0, 0, this.shadowTexture);

        // Set the blend mode to MULTIPLY. This will darken the colors of
        // everything below this sprite.
        lightSprite.blendMode = Phaser.blendModes.MULTIPLY;

        this.redrawLogMessages();

        this.initWebsockets();
    }


    redrawLogMessages()
    {
        this.logMessageWindow = this.game.add.group();
        this.logMessageWindow.x = 100;
        this.logMessageWindow.y = 100;
        this.logMessageWindow.z = 20;
        this.logMessageWindow.fixedToCamera = true;




        // XXX make some pretty log window with scroll
        for (let msg of this.logMessages) {
            console.log(msg);

            var txt = this.makeText(msg.text);

            // floating name over head of player
            var img = this.game.add.image(0, -10, txt);
            this.logMessageWindow.add(img);
        }
    }


    update()
    {
        if (!this.playerGroup) {
            return;
        }

        // Update the shadow texture each frame
        //this.updateShadowTexture();

        this.game.physics.arcade.collide(this.player, this.groundLayer);

        var steppingVert = 2;
        var steppingHoriz = 4;

        // flip horizontally
        if (this.player.body.velocity.x == this.cursors.left.isDown) {
            this.player.scale.x = -1;
        } else if (this.player.body.velocity.x == this.cursors.right.isDown) {
            this.player.scale.x = 1;
        }

        if (this.cursors.up.isDown) {
            this.playerGroup.y -= steppingVert;
            this.sendMove();
        } else if (this.cursors.down.isDown) {
            this.playerGroup.y += steppingVert;
            this.sendMove();
        }

        if (this.cursors.left.isDown) {
            this.playerGroup.x -= steppingHoriz;
            this.sendMove();
        } else if (this.cursors.right.isDown) {
            this.playerGroup.x += steppingHoriz;
            this.sendMove();
        }

        if (this.game.input.keyboard.isDown(Phaser.Keyboard.Q)) {
            this.worldScale += 0.05;
        } else if (this.game.input.keyboard.isDown(Phaser.Keyboard.A)) {
            this.worldScale -= 0.05;
        }

        // set a minimum and maximum scale value
        this.worldScale = Phaser.Math.clamp(this.worldScale, 0.25, 2);

        // set our world scale as needed
        this.stageGroup.scale.set(this.worldScale);

        // XXX game.camera
        this.game.camera.setSize(this.game.width, this.game.height);
        this.game.camera.setBoundsToWorld();
        this.game.camera.follow(this.playerGroup);
        this.game.camera.update();

        this.groundLayer.resizeWorld();
    }

    render()
    {
        this.game.debug.text(this.game.time.fps || '--', 1, 14, "#00ff00");

        //game.debug.spriteInfo(this.player, 32, 32);
        //game.debug.cameraInfo(game.camera, 10, 32);

        //game.debug.soundInfo(this.music, 10, 140);
    }



    updateShadowTexture()
    {
        // This function updates the shadow texture (this.shadowTexture).
        // First, it fills the entire texture with a dark shadow color.
        // Then it draws a white circle centered on the pointer position.
        // Because the texture is drawn to the screen using the MULTIPLY
        // blend mode, the dark areas of the texture make all of the colors
        // underneath it darker, while the white area is unaffected.

        // Draw shadow
        this.shadowTexture.context.fillStyle = 'rgb(100, 100, 100)';
        this.shadowTexture.context.fillRect(0, 0, this.game.width, this.game.height);

        // Draw circle of light with a soft edge
        var gradient = this.shadowTexture.context.createRadialGradient(
            this.playerGroup.x, this.playerGroup.y, this.LIGHT_RADIUS * 0.75,
            this.playerGroup.x, this.playerGroup.y, this.LIGHT_RADIUS
        );
        gradient.addColorStop(0, 'rgba(255, 255, 255, 1.0)');
        gradient.addColorStop(1, 'rgba(255, 255, 255, 0.0)');

        this.shadowTexture.context.beginPath();
        this.shadowTexture.context.fillStyle = gradient;
        this.shadowTexture.context.arc(
            this.playerGroup.x,
            this.playerGroup.y,
            this.LIGHT_RADIUS,
            0,
            Math.PI * 2
        );
        this.shadowTexture.context.fill();

        // This just tells the engine it should update the texture cache
        this.shadowTexture.dirty = true;
    }

    initWebsockets()
    {
        var url = 'ws://localhost:3322/ws';
        this.socket = new WebSocket(url);

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
            //console.log('Websocket connected');
            this.send("new_player " + this.playerName);
        };
    }


    sendMove()
    {
        var newX = Math.floor(this.playerGroup.x / this.tileWidth);
        var newY = Math.floor(this.playerGroup.y / this.tileHeight);

        if (this.prevX == newX && this.prevY == newY) {
            // dont spam server when coords havent changed
            return;
        }

        this.socket.send("move " + newX + " " + newY + " " + this.token);
        this.prevX = newX;
        this.prevY = newY;
    }



    handleXyMessage(cmd)
    {
        this.playerGroup = this.game.add.group();
        this.playerGroup.z = 10;
        this.stageGroup.add(this.playerGroup);


        this.player = this.game.add.sprite(0, 0, 'characterAtlas');
        this.player.frameName = 'dwarf';
        this.player.anchor.set(0.5);
        this.game.camera.follow(this.playerGroup);

        this.game.physics.enable(this.player);


        //  Because both our body and our tiles are so tiny,
        //  and the body is moving pretty fast, we need to add
        //  some tile padding to the body. WHat this does
        this.player.body.tilePadding.set(32, 32);


        // multiply coords with tile size to scale properly.
        // sprite tiles are always in pixels
        this.playerGroup.x = cmd.X * this.tileWidth;
        this.playerGroup.y = cmd.Y * this.tileHeight;
        this.playerGroup.add(this.player);

        var txt = this.makeText(this.playerName);

        // floating name over head of player
        var aboveHead = this.game.add.image(0, -10, txt);
        aboveHead.anchor.set(0.5);
        this.playerGroup.add(aboveHead);
        console.log("spawned at " + cmd.X + ", " + cmd.Y);

        this.token = cmd.Token;

        this.renderLocalSpawns(cmd.LocalSpawns);
    }

    makeText(msg)
    {
        var oddballFontSet = "                " + // colors
            "                " + // cursor
            "!\"#$%&'()  ,-./0123456789:;<=>?@" +
            "ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`" +
            "abcdefghijklmnopqrstuvwxyz{|}~" +
            ""; // XXX more characters

        var o = this.game.add.retroFont('oddballFont', 8, 8, oddballFontSet, 16);
        o.autoUpperCase = false;
        o.text = msg;

        return o;
    }

    handleMoveResMessage(cmd)
    {
        // console.log("Rendering " + cmd.LocalSpawns.length + " spawns at " + cmd.X + ", " + cmd.Y);
        this.renderLocalSpawns(cmd.LocalSpawns);
    }

    renderLocalSpawns(spawns)
    {
        this.spawnLayer.removeAll();

        var atlas = "";

        for (var i = 0; i < spawns.length; i++) {
            var sp = spawns[i];
            if (sp.Class == "player" && sp.Name == this.playerName) {
                continue;
            }

            var values = sp.Sprite.split(':');
            switch (values[0]) {
            case 'c':
                atlas = 'characterAtlas';
                break;
            case 'i':
                atlas = 'itemAtlas';
                break;
            case 'g':
                atlas = 'ground2Atlas';
                break;
            default:
                console.log('ERROR unknown sprite: ' + sp.Sprite);
                console.log(sp);
                continue;
            }

            var spr = this.game.add.sprite(0, 0, atlas);
            spr.x = sp.X * this.tileWidth;
            spr.y = sp.Y * this.tileHeight;
            spr.frameName = values[1];
            spr.anchor.set(0.5);

            this.spawnLayer.add(spr);
        }
    }
}

export default GameState;
