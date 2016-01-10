(function e(t,n,r){function s(o,u){if(!n[o]){if(!t[o]){var a=typeof require=="function"&&require;if(!u&&a)return a(o,!0);if(i)return i(o,!0);var f=new Error("Cannot find module '"+o+"'");throw f.code="MODULE_NOT_FOUND",f}var l=n[o]={exports:{}};t[o][0].call(l.exports,function(e){var n=t[o][1][e];return s(n?n:e)},l,l.exports,e,t,n,r)}return n[o].exports}var i=typeof require=="function"&&require;for(var o=0;o<r.length;o++)s(r[o]);return s})({1:[function(require,module,exports){
'use strict';

var gameWidth = 800;
var gameHeight = 400;

var tileWidth = 8;
var tileHeight = 4;

var worldScale = 1.0;

var GameState = function GameState(game) {};

GameState.prototype.preload = function () {
    // SCALE TO FIT SCREEN:
    //game.scale.scaleMode = Phaser.ScaleManager.SHOW_ALL;
    //game.scale.setGameSize(gameWidth, gameHeight);

    game.time.advancedTiming = true; // required for fps counter

    game.stage.backgroundColor = '#262f71'; // deep water

    game.load.tilemap('islandMap', '/island/full', null, Phaser.Tilemap.TILED_JSON);
    game.load.image('ground', 'img/tileset/oddball/ground.png', 4, 8);

    game.load.image('minimap', 'img/islands/current.png');

    game.load.atlas('characterAtlas', 'img/tileset/oddball/characters.png', 'sprite/character');
    game.load.atlas('itemAtlas', 'img/tileset/oddball/items.png', 'sprite/item');
    game.load.atlas('ground2Atlas', 'img/tileset/oddball/ground2.png', 'sprite/ground2');

    game.load.image('oddballFont', 'img/tileset/oddball/font.png');

    game.load.audio('bgSound', ['audio/dead_feelings.mp3']);
};

var player;
var playerName;
var playerGroup;

var token;

var spawnLayer; // this group holds all nearby spawns

var oddballFontSet = "                " + // colors
"                " + // cursor
"!\"#$%&'()  ,-./0123456789:;<=>?@" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`" + "abcdefghijklmnopqrstuvwxyz{|}~" + ""; // XXX more characters

GameState.prototype.create = function () {
    // world (except UI) is in this group, so it can be scaled
    this.stageGroup = game.add.group();

    this.groundMap = game.add.tilemap('islandMap');
    this.groundMap.addTilesetImage('island_tiles', 'ground');
    this.groundLayer = this.groundMap.createLayer(0);
    this.groundLayer.resizeWorld();

    this.stageGroup.add(this.groundLayer);

    this.music = game.add.audio('bgSound');
    this.music.volume = 0.5; // 50%
    this.music.play();

    spawnLayer = game.add.group();
    spawnLayer.z = 5;
    this.stageGroup.add(spawnLayer);

    this.cursors = game.input.keyboard.createCursorKeys();

    var minimapScale = 3;
    this.minimap = game.add.sprite(gameWidth - game.cache.getImage('minimap').width / minimapScale, 0, 'minimap');
    this.minimap.fixedToCamera = true;
    this.minimap.scale.set(1.0 / minimapScale);
    this.minimap.alpha = 0.8;
    this.minimap.setScaleMinMax(1.0 / minimapScale, 1.0 / minimapScale);

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

    this.initWebsockets();
};

GameState.prototype.update = function () {
    if (!playerGroup) {
        return;
    }

    // Update the shadow texture each frame
    //this.updateShadowTexture();

    game.physics.arcade.collide(player, this.groundLayer);

    var steppingVert = 2;
    var steppingHoriz = 4;

    // flip horizontally
    if (player.body.velocity.x == this.cursors.left.isDown) {
        player.scale.x = -1;
    } else if (player.body.velocity.x == this.cursors.right.isDown) {
        player.scale.x = 1;
    }

    if (this.cursors.up.isDown) {
        playerGroup.y -= steppingVert;
        this.sendMove();
    } else if (this.cursors.down.isDown) {
        playerGroup.y += steppingVert;
        this.sendMove();
    }

    if (this.cursors.left.isDown) {
        playerGroup.x -= steppingHoriz;
        this.sendMove();
    } else if (this.cursors.right.isDown) {
        playerGroup.x += steppingHoriz;
        this.sendMove();
    }

    // XXX  zoom is broken
    // zoom
    if (game.input.keyboard.isDown(Phaser.Keyboard.Q)) {
        worldScale += 0.05;
    } else if (game.input.keyboard.isDown(Phaser.Keyboard.A)) {
        worldScale -= 0.05;
    }

    // set our world scale as needed
    //game.world.scale.set(worldScale);

    this.stageGroup.scale.x = worldScale;
    this.stageGroup.scale.y = worldScale;
};

GameState.prototype.render = function () {
    game.debug.text(game.time.fps || '--', 1, 14, "#00ff00");

    //game.debug.spriteInfo(player, 32, 32);
    //game.debug.cameraInfo(game.camera, 32, 32);

    //game.debug.soundInfo(music, 20, 32);
};

GameState.prototype.updateShadowTexture = function () {
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
    var gradient = this.shadowTexture.context.createRadialGradient(playerGroup.x, playerGroup.y, this.LIGHT_RADIUS * 0.75, playerGroup.x, playerGroup.y, this.LIGHT_RADIUS);
    gradient.addColorStop(0, 'rgba(255, 255, 255, 1.0)');
    gradient.addColorStop(1, 'rgba(255, 255, 255, 0.0)');

    this.shadowTexture.context.beginPath();
    this.shadowTexture.context.fillStyle = gradient;
    this.shadowTexture.context.arc(playerGroup.x, playerGroup.y, this.LIGHT_RADIUS, 0, Math.PI * 2);
    this.shadowTexture.context.fill();

    // This just tells the engine it should update the texture cache
    this.shadowTexture.dirty = true;
};

var socket;

GameState.prototype.initWebsockets = function () {
    var url = 'ws://localhost:3322/ws';
    socket = new WebSocket(url);

    parent = this;

    /**
     * @param msg MessageEvent
     */
    socket.onmessage = function (msg) {
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

    socket.onopen = function () {
        //console.log('Websocket connected');
        socket.send("new_player mrcool");
    };
};

GameState.prototype.sendMove = function () {
    var newX = Math.floor(playerGroup.x / tileWidth);
    var newY = Math.floor(playerGroup.y / tileHeight);

    if (this.prevX == newX && this.prevY == newY) {
        // dont spam server when coords havent changed
        return;
    }

    socket.send("move " + newX + " " + newY + " " + token);

    this.prevX = newX;
    this.prevY = newY;
};

GameState.prototype.handleXyMessage = function (cmd) {
    playerGroup = game.add.group();
    playerGroup.z = 10;
    this.stageGroup.add(playerGroup);

    player = game.add.sprite(0, 0, 'characterAtlas');
    player.frameName = 'dwarf';
    player.anchor.set(0.5);
    game.camera.follow(playerGroup);
    game.physics.enable(player);

    //  Because both our body and our tiles are so tiny,
    //  and the body is moving pretty fast, we need to add
    //  some tile padding to the body. WHat this does
    player.body.tilePadding.set(32, 32);

    // multiply coords with tile size to scale properly. sprite tiles are always in pixels
    playerName = cmd.Name;

    playerGroup.x = cmd.X * tileWidth;
    playerGroup.y = cmd.Y * tileHeight;
    playerGroup.add(player);

    var playerNameFont = game.add.retroFont('oddballFont', 8, 8, oddballFontSet, 16);
    playerNameFont.autoUpperCase = false;
    playerNameFont.text = playerName;

    // floating name over head of player
    var aboveHead = game.add.image(0, -10, playerNameFont);
    aboveHead.anchor.set(0.5);
    playerGroup.add(aboveHead);
    console.log("spawned at " + cmd.X + ", " + cmd.Y);

    token = cmd.Token;

    this.renderLocalSpawns(cmd.LocalSpawns);
};

GameState.prototype.handleMoveResMessage = function (cmd) {
    // console.log("Rendering " + cmd.LocalSpawns.length + " spawns at " + cmd.X + ", " + cmd.Y);
    this.renderLocalSpawns(cmd.LocalSpawns);
};

GameState.prototype.renderLocalSpawns = function (spawns) {
    spawnLayer.removeAll();

    var atlas = "";

    for (var i = 0; i < spawns.length; i++) {
        var sp = spawns[i];
        if (sp.Class == "player" && sp.Name == playerName) {
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

        var spr = game.add.sprite(0, 0, atlas);
        spr.x = sp.X * tileWidth;
        spr.y = sp.Y * tileHeight;
        spr.frameName = values[1];
        spr.anchor.set(0.5);

        spawnLayer.add(spr);
    }
};

var game = new Phaser.Game(gameWidth, gameHeight, Phaser.CANVAS, 'game', {}, false, // transparent
false // antialias
);
game.state.add('game', GameState, true);

},{}]},{},[1]);
