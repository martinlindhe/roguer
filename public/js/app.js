(function e(t,n,r){function s(o,u){if(!n[o]){if(!t[o]){var a=typeof require=="function"&&require;if(!u&&a)return a(o,!0);if(i)return i(o,!0);var f=new Error("Cannot find module '"+o+"'");throw f.code="MODULE_NOT_FOUND",f}var l=n[o]={exports:{}};t[o][0].call(l.exports,function(e){var n=t[o][1][e];return s(n?n:e)},l,l.exports,e,t,n,r)}return n[o].exports}var i=typeof require=="function"&&require;for(var o=0;o<r.length;o++)s(r[o]);return s})({1:[function(require,module,exports){
'use strict';

var gameWidth = 800;
var gameHeight = 600;

var tileWidth = 8;
var tileHeight = 4;

var game = new Phaser.Game(gameWidth, gameHeight, Phaser.CANVAS, 'game', {
    preload: preload,
    create: create,
    update: update,
    render: render
});

function preload() {
    game.time.advancedTiming = true; // required for fps counter

    game.stage.backgroundColor = '#262f71'; // deep water

    game.load.image('minimap', 'img/islands/current.png');

    game.load.tilemap('island', '/island/full', null, Phaser.Tilemap.TILED_JSON);
    game.load.image('tiles', 'img/tileset/oddball/ground1.png', 4, 8);

    game.load.audio('bgSound', ['audio/dead_feelings.mp3']);

    game.load.atlas('atlas', 'img/tileset/oddball/characters.png', 'sprite/character');

    game.load.image('oddballFont', 'img/tileset/oddball/font.png');
}

var map;
var layer;
var cursors;
var player;
var playerGroup;
var music;
var minimap;
var retroFont;

var token;
var worldScale = 1.0;

var oddballFontSet = "                " + // colors
"                " + // cursor
"!\"#$%&'()  ,-./0123456789:;<=>?@" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`" + "abcdefghijklmnopqrstuvwxyz{|}~" + ""; // XXX more characters

function create() {
    music = game.add.audio('bgSound');
    music.volume = 0.5; // 50%
    music.play();

    // A Tilemap object just holds the data needed to describe the map
    // You can add your own data or manipulate the data (swap tiles around, etc)
    // but in order to display it you need to create a TilemapLayer.
    map = game.add.tilemap('island');

    map.addTilesetImage('island_tiles', 'tiles');

    layer = map.createLayer(0);

    // Basically this sets EVERY SINGLE tile to fully collide on all faces
    // map.setCollisionByExclusion([7, 32, 35, 36, 47]);

    layer.resizeWorld();

    cursors = game.input.keyboard.createCursorKeys();

    var minimapScale = 3;
    minimap = game.add.sprite(gameWidth - game.cache.getImage('minimap').width / minimapScale, 0, 'minimap');
    minimap.fixedToCamera = true;
    minimap.scale.set(1.0 / minimapScale);
    minimap.alpha = 0.8;

    //minimap.setScaleMinMax(1, 1);

    initWebsockets();
}

function update() {
    if (!playerGroup) {
        return;
    }

    //game.physics.arcade.collide(player, layer);

    var steppingVert = 2;
    var steppingHoriz = 4;

    // flip horizontally
    if (player.body.velocity.x = cursors.left.isDown) {
        player.scale.x = -1;
    } else if (player.body.velocity.x = cursors.right.isDown) {
        player.scale.x = 1;
    }

    if (cursors.up.isDown) {
        playerGroup.y -= steppingVert;
        sendSocketMove();
    } else if (cursors.down.isDown) {
        playerGroup.y += steppingVert;
        sendSocketMove();
    }

    if (cursors.left.isDown) {
        playerGroup.x -= steppingHoriz;
        sendSocketMove();
    } else if (cursors.right.isDown) {
        playerGroup.x += steppingHoriz;
        sendSocketMove();
    }

    /*  zoom is broken
    // zoom
    if (game.input.keyboard.isDown(Phaser.Keyboard.Q)) {
        worldScale += 0.05;
    }
    else if (game.input.keyboard.isDown(Phaser.Keyboard.A)) {
        worldScale -= 0.05;
    }
     // set our world scale as needed
    game.world.scale.set(worldScale);
    */
}

function render() {
    game.debug.text(game.time.fps || '--', 1, 14, "#00ff00");

    //game.debug.spriteInfo(player, 32, 32);
    //game.debug.cameraInfo(game.camera, 32, 32);

    //game.debug.soundInfo(music, 20, 32);
}

var socket;

function initWebsockets() {
    var url = 'ws://localhost:3322/ws';
    socket = new WebSocket(url);

    socket.onmessage = onSocketMessage;
    socket.onopen = onSocketConnected;
}

/**
 * @param msg MessageEvent
 */
function onSocketMessage(msg) {
    var cmd = JSON.parse(msg.data);

    switch (cmd.Type) {
        case 'xy':
            // multiply coords with tile size to scale properly. sprite tiles are always in pixels

            playerGroup = game.add.group();
            playerGroup.z = 10;

            player = game.add.sprite(0, 0, 'atlas');
            player.frameName = 'dwarf';
            player.anchor.set(0.5);
            game.camera.follow(player);
            game.physics.enable(player);

            //  Because both our body and our tiles are so tiny,
            //  and the body is moving pretty fast, we need to add
            //  some tile padding to the body. WHat this does
            player.body.tilePadding.set(32, 32);

            playerGroup.x = cmd.X * tileWidth;
            playerGroup.y = cmd.Y * tileHeight;
            playerGroup.add(player);

            retroFont = game.add.retroFont('oddballFont', 8, 8, oddballFontSet, 16);
            retroFont.autoUpperCase = false;
            retroFont.text = cmd.Name;

            // floating name over head of player
            var playerName = game.add.image(0, -10, retroFont);

            playerName.anchor.set(0.5);

            playerGroup.add(playerName);
            console.log("spawned at " + cmd.X + ", " + cmd.Y);

            token = cmd.Token;

            // display all from .LocalSpawns
            //console.log(cmd.LocalSpawns);
            for (var i = 0; i < cmd.LocalSpawns.length; i++) {
                var sp = cmd.LocalSpawns[i];
                //console.log(sp)

                var spGroup = game.add.group();
                spGroup.x = sp.X * tileWidth;
                spGroup.y = sp.Y * tileHeight;
                spGroup.z = 10;

                var spr = game.add.sprite(0, 0, 'atlas');
                spr.frameName = sp.Sprite;
                spr.anchor.set(0.5);

                spGroup.add(spr);

                // XXX add to game world ...
                //game.add.sprite(sp.X * tileWidth, sp.Y * tileHeight, 'phaser');
            }
            break;

        case 'ok':
            console.log("server OK: " + msg.data);
            break;

        default:
            console.log("<-recv- " + msg.data);
            console.log("unknown command from server: " + cmd.Type);
    }
}

function sendSocketMsg(data) {
    socket.send(data);
    console.log("-sent->" + data);
}

function sendSocketMove() {
    sendSocketMsg("move " + Math.floor(player.x / tileWidth) + " " + Math.floor(player.y / tileHeight) + " " + token);
}

// Socket connected
function onSocketConnected() {
    sendSocketMsg("new_player mrcool");

    console.log('Connected to socket server');
}

},{}]},{},[1]);
