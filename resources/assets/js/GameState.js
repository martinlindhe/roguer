import {Client} from './Client.js';
import {MessageLog} from './MessageLog.js';
import {GameTime} from './GameTime.js';

export class GameState extends Phaser.State
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

        // NOTE: topaz-8.woff is force loaded with css hack, see fontLoader

        this.game.scale.scaleMode = Phaser.ScaleManager.RESIZE;
        //this.scale.pageAlignVertically = true;

        //window.addEventListener('resize', () => this.game.scale.refresh());
        //this.game.scale.refresh();
    }

    create()
    {
        this.playerName = "Jimpson";

        this.worldScale = 1.0;
        this.tileWidth = 8;
        this.tileHeight = 4;

        this.maxMessages = 15;
        this.logTextHeight = 15;
        this.minimapScale = 3;

        this.serverTime = new GameTime(0);

        // scale to whole window
        //this.game.scale.setGameSize(window.innerWidth, window.innerHeight);

        this.game.physics.startSystem(Phaser.Physics.ARCADE);


        this.music = this.game.add.audio('bgSound');
        this.music.volume = 0.20; // 20%



        this.groundMap = this.game.add.tilemap('islandMap');
        this.groundMap.addTilesetImage('island_tiles', 'ground');
        this.groundMap.setCollisionBetween(0, 112); // 112 = beach line

        this.groundLayer = this.groundMap.createLayer(0);
        //this.groundLayer.resizeWorld();

        // Un-comment this on to see the collision tiles
        //this.groundLayer.debug = true;


        this.playerGroup = this.game.add.group();
        //this.playerGroup.z = 10;


        this.spawnLayer = this.game.add.group();
        //this.spawnLayer.z = 0;



        this.initUi();

        //  Capture all key presses
        this.game.input.keyboard.addCallbacks(this, null, null, function(char) {
            // console.log("pressed "+char);
            switch (char) {
            case 't':
                this.ui.visible = !this.ui.visible;
                return;
            case 'q':
                this.worldScale += 0.2;
                break;
            case 'a':
                this.worldScale -= 0.2;
                break;
            default:
                console.log("unhandled char " + char);
                return;
            }

            // set a minimum and maximum scale value
            this.worldScale = Phaser.Math.clamp(this.worldScale, 0.5, 4);

            // set our world scale as needed
            this.game.scale.set(this.worldScale); // XXX
        });



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
    }

    update()
    {
        if (!this.playerSprite) {
            return;
        }

        // Update the shadow texture each frame
        //this.updateShadowTexture();

        this.game.physics.arcade.collide(this.playerSprite, this.groundLayer);

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
            this.playerSprite.y -= steppingVert;
            this.client.sendMove();
        } else if (this.cursors.down.isDown) {
            this.playerSprite.y += steppingVert;
            this.client.sendMove();
        }

        if (this.cursors.left.isDown) {
            this.playerSprite.x -= steppingHoriz;
            this.client.sendMove();
        } else if (this.cursors.right.isDown) {
            this.playerSprite.x += steppingHoriz;
            this.client.sendMove();
        }

/*
        // XXX game.camera
        this.game.camera.setSize(this.game.width, this.game.height);
        this.game.camera.setBoundsToWorld();
        this.game.camera.follow(this.playerGroup);
        this.game.camera.update();
*/


        var s = "Game size: " + this.game.width + " x " + this.game.height + "\n";
        s += "Actual size: " + this.game.scale.width + " x " + this.game.scale.height + "\n";
        s += "minWidth: " + this.game.scale.minWidth + " - minHeight: " + this.game.scale.minHeight + "\n";
        s += "maxWidth: " + this.game.scale.maxWidth + " - maxHeight: " + this.game.scale.maxHeight + "\n";
        s += "aspect ratio: " + this.game.scale.aspectRatio + "\n";
        s += "parent is window: " + this.game.scale.parentIsWindow + "\n";
        s += "bounds x: " + this.game.scale.bounds.x + " y: " + this.game.scale.bounds.y + " width: " + this.game.scale.bounds.width + " height: " + this.game.scale.bounds.height + "\n";
        this.info.text = s;
    }

    render()
    {
        this.game.debug.text(this.game.time.fps || '--', 1, 14, "#00ff00");

        if (this.playerSprite) {
            this.game.debug.spriteInfo(this.playerSprite, 400, 20);
        }
        //game.debug.cameraInfo(game.camera, 10, 32);

        //game.debug.soundInfo(this.music, 10, 140);
    }

    messageToLog(o)
    {
        this.logMessageList.text = this.messageLog.add(o).render();
    }

    initUi()
    {
        this.ui = this.game.add.group();
        this.ui.fixedToCamera = true;

        this.cursors = this.game.input.keyboard.createCursorKeys();



        this.info = this.game.add.text(16, 16, ' ');
        this.info.font = "Courier";
        this.info.fontSize = 14;
        this.info.fill = "#fff";
        this.info.lineSpacing = -10;
        this.info.setShadow(2, 2);
        this.ui.add(this.info);



        this.minimap = this.game.add.sprite(0, 0, 'minimap');
        this.minimap.scale.set(1.0 / this.minimapScale);
        this.minimap.alpha = 0.8;
        this.ui.add(this.minimap);

        this.muteButton = this.game.add.button(
            0,
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
        this.ui.add(this.muteButton);


        // text-shadow hack for making text-stroke
        var style = { font: "10px topaz", fill: "#fff", backgroundColor: 'rgba(0,0,0,0.25)', wordWrap: true, wordWrapWidth: 400 };

        this.logMessageList = this.game.add.text(0, 0, '', style);
        this.logMessageList.stroke = '#000';
        this.logMessageList.strokeThickness = 2;
        this.logMessageList.lineSpacing = -8;

        this.ui.add(this.logMessageList);

        this.messageLog = new MessageLog();
        this.logMessageList.text = this.messageLog.render();

        this.timeOfDayIcon = this.game.add.text(0, 0, '', { fill : '#fff', font : '18px weathericons' });
        this.timeOfDayIcon.stroke = '#000';
        this.timeOfDayIcon.strokeThickness = 2;

        this.ui.add(this.timeOfDayIcon);


        // shows server time :
        this.serverTimeText = this.game.add.text(0, 0, "", style);
        this.ui.add(this.serverTimeText);

        this.resize(this.game.width, this.game.height);
    }


    resize(width, height)
    {
        console.log("resized to w " + width + ", h = " + height);

        this.minimap.x = width - this.game.cache.getImage('minimap').width / this.minimapScale;

        this.muteButton.x = width - 102;

        this.logMessageList.x = width - 400;
        this.logMessageList.y = height - (this.maxMessages * this.logTextHeight);

        this.timeOfDayIcon.x = width - 150;

        this.serverTimeText.x = width - 230;

        this.game.scale.refresh();
    }

    setServerTime(i)
    {
        this.serverTime.set(i);

        if (this.serverTime.time) {
            this.serverTimeText.text = this.serverTime.render();

            this.updateTimeOfDayIcon();
        }

        // auto save message log on every ping
        this.messageLog.save();
    }

    // uses https://erikflowers.github.io/weather-icons/
    updateTimeOfDayIcon()
    {
        var hour = this.serverTime.hour();
        if (hour >= 12) {
            hour -= 12;
        }

        // wi-time-*
        var timesOfDay = [
            "\uf089", "\uf08a", "\uf08b", "\uf08c", "\uf08d", "\uf08e",
            "\uf08f", "\uf090", "\uf091", "\uf092", "\uf093", "\uf094"
        ];

        var val = timesOfDay[hour];
        if (this.serverTime.isDaytime()) {
            val += " \uf00d"; // wi-day-sunny
        } else {
            val += " \uf02e"; // wi-night-clear
        }

        this.timeOfDayIcon.text = val + " ";
    }

    spawnPlayer(cmd)
    {
        this.playerSprite = this.game.add.sprite(0, 0, 'characterAtlas');
        this.game.physics.enable(this.playerSprite, Phaser.Physics.ARCADE);
        this.game.camera.follow(this.playerSprite);

        this.playerSprite.frameName = 'dwarf';
        this.playerSprite.anchor.set(0.5);
        //this.playerSprite.body.collideWorldBounds = true;



        //  Because both our body and our tiles are so tiny,
        //  and the body is moving pretty fast, we need to add
        //  some tile padding to the body. WHat this does
        //this.playerSprite.body.tilePadding.set(32, 32);


        // multiply coords with tile size to scale properly.
        // sprite tiles are always in pixels
        this.playerSprite.x = cmd.X * this.tileWidth;
        this.playerSprite.y = cmd.Y * this.tileHeight;

/*
        // floating name over head of player
        var txt = this.makeText(this.playerName);
        var aboveHead = this.game.add.image(0, -10, txt);
        aboveHead.anchor.set(0.5);
        this.playerGroup.add(aboveHead);
*/
        console.log("spawned at " + cmd.X + ", " + cmd.Y);

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
