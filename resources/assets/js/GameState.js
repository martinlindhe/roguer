import Client from './Client.js';

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
        this.playerName = "Jimpson";

        this.worldScale = 1.0;
        this.tileWidth = 8;
        this.tileHeight = 4;

        this.logMessages = [
            {text: "hello there", time: 1222},
            {text: "later on", time: 1234},
        ];

        this.game.physics.startSystem(Phaser.Physics.ARCADE);

        // scale to whole window
        this.game.scale.setGameSize(window.innerWidth, window.innerHeight);



        this.music = this.game.add.audio('bgSound');
        this.music.volume = 0.20; // 20%



        // world (except UI) is in this group, so it can be scaled
        this.stageGroup = this.game.add.group();

        this.groundMap = this.game.add.tilemap('islandMap');
        this.groundMap.addTilesetImage('island_tiles', 'ground');

        this.groundMap.setCollisionBetween(0, 112); // 112 = beach line

        this.groundLayer = this.groundMap.createLayer(0);
        this.groundLayer.resizeWorld();

        // Un-comment this on to see the collision tiles
        //this.groundLayer.debug = true;

        this.stageGroup.add(this.groundLayer);


        this.playerGroup = this.game.add.group();
        this.game.camera.follow(this.playerGroup);

        //this.playerGroup.z = 10;
        this.stageGroup.add(this.playerGroup);


        this.spawnLayer = this.game.add.group();
        //this.spawnLayer.z = 0;
        this.stageGroup.add(this.spawnLayer);


        this.initUi();




/*
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
*/


        this.client = new Client(this);
        this.redrawLogMessages();
    }

    update()
    {
        if (!this.playerSprite) {
            return;
        }

        // Update the shadow texture each frame
        //this.updateShadowTexture();

        this.game.physics.arcade.collide(this.playerGroup, this.groundLayer);

        var steppingHoriz = 1;
        var steppingVert = steppingHoriz / 2;

        /*
        // flip horizontally, XXX this should be done on all npc movements to show facing direction
        if (this.playerSprite.body.velocity.x == this.cursors.left.isDown) {
            this.playerSprite.scale.x = -1;
        } else if (this.playerSprite.body.velocity.x == this.cursors.right.isDown) {
            this.playerSprite.scale.x = 1;
        }
        */

        if (this.cursors.up.isDown) {
            this.playerGroup.y -= steppingVert;
            this.client.sendMove();
        } else if (this.cursors.down.isDown) {
            this.playerGroup.y += steppingVert;
            this.client.sendMove();
        }

        if (this.cursors.left.isDown) {
            this.playerGroup.x -= steppingHoriz;
            this.client.sendMove();
        } else if (this.cursors.right.isDown) {
            this.playerGroup.x += steppingHoriz;
            this.client.sendMove();
        }

        if (this.game.input.keyboard.isDown(Phaser.Keyboard.Q)) {
            this.worldScale += 0.05;
        } else if (this.game.input.keyboard.isDown(Phaser.Keyboard.A)) {
            this.worldScale -= 0.05;
        }

        // set a minimum and maximum scale value
        this.worldScale = Phaser.Math.clamp(this.worldScale, 0.5, 4);

        // set our world scale as needed
        this.stageGroup.scale.set(this.worldScale);

/*
        // XXX game.camera
        this.game.camera.setSize(this.game.width, this.game.height);
        this.game.camera.setBoundsToWorld();
        this.game.camera.follow(this.playerGroup);
        this.game.camera.update();

        this.groundLayer.resizeWorld();
*/
    }

    render()
    {
        this.game.debug.text(this.game.time.fps || '--', 1, 14, "#00ff00");

        //game.debug.spriteInfo(this.player, 32, 32);
        //game.debug.cameraInfo(game.camera, 10, 32);

        //game.debug.soundInfo(this.music, 10, 140);
    }

    initUi()
    {
        this.uiGroup = this.game.add.group();


        this.cursors = this.game.input.keyboard.createCursorKeys();

        var minimapScale = 3;
        var minimapX = this.game.width - this.game.cache.getImage('minimap').width / minimapScale;
        var minimap = this.game.add.sprite(minimapX, 0, 'minimap');
        minimap.fixedToCamera = true;
        minimap.scale.set(1.0 / minimapScale);
        minimap.alpha = 0.8;
        minimap.setScaleMinMax(1.0 / minimapScale, 1.0 / minimapScale);
        this.uiGroup.add(minimap);

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
        this.uiGroup.add(button);



        this.logMessageWindow = this.game.add.group();
        this.logMessageWindow.x = 10;
        this.logMessageWindow.y = 10;
        this.logMessageWindow.z = 20;
        this.logMessageWindow.fixedToCamera = true;

        this.uiGroup.add(this.logMessageWindow);
    }

    spawnPlayer(cmd)
    {
        this.playerSprite = this.game.add.sprite(0, 0, 'characterAtlas');
        this.playerSprite.frameName = 'dwarf';
        this.playerSprite.anchor.set(0.5);



        //  Because both our body and our tiles are so tiny,
        //  and the body is moving pretty fast, we need to add
        //  some tile padding to the body. WHat this does
        //this.playerSprite.body.tilePadding.set(32, 32);


        // multiply coords with tile size to scale properly.
        // sprite tiles are always in pixels
        this.playerGroup.x = cmd.X * this.tileWidth;
        this.playerGroup.y = cmd.Y * this.tileHeight;
        this.playerGroup.add(this.playerSprite);


        // floating name over head of player
        var txt = this.makeText(this.playerName);
        var aboveHead = this.game.add.image(0, -10, txt);
        aboveHead.anchor.set(0.5);
        this.playerGroup.add(aboveHead);

        console.log("spawned at " + cmd.X + ", " + cmd.Y);

        this.game.physics.enable(this.playerGroup);

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


    redrawLogMessages()
    {
        this.logMessageWindow.removeAll();

        // XXX log window with scroll
        for (let msg of this.logMessages) {
            console.log(msg);

            var txt = this.makeText(msg.text);

            var img = this.game.add.image(0, -10, txt);
            this.logMessageWindow.add(img);
        }
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
}

export default GameState;
