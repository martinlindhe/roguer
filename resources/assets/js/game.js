var gameWidth = 800;
var gameHeight = 600;

var tileWidth = 8;
var tileHeight = 4;

var game = new Phaser.Game(
    gameWidth,
    gameHeight,
    Phaser.CANVAS,
    'game',
    {
        preload: preload,
        create: create,
        update: update,
        render : render
    }
);

function preload()
{
    game.time.advancedTiming = true; // required for fps counter

    game.stage.backgroundColor = '#262f71';  // deep water

    game.load.image('minimap', 'img/islands/current.png');

    game.load.tilemap('island', '/island/full', null, Phaser.Tilemap.TILED_JSON);
    game.load.image('tiles', 'img/tileset/oddball/ground1.png', 4, 8);

    game.load.audio('bgSound', ['audio/dead_feelings.mp3']);

    game.load.atlas('atlas', 'img/tileset/oddball/characters.png', 'sprite/character');
}


var map;
var layer;
var cursors;
var player;
var music;
var minimap;

var token;
var worldScale = 1.0;

function create()
{
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



    player = game.add.sprite(10, 10, 'atlas');
    player.frameName = 'dwarf';


    player.visible = false;
    player.anchor.set(0.5);
    game.camera.follow(player);

    game.physics.enable(player);

    //  Because both our body and our tiles are so tiny,
    //  and the body is moving pretty fast, we need to add
    //  some tile padding to the body. WHat this does
    player.body.tilePadding.set(32, 32);

    var minimapScale = 3
    minimap = game.add.sprite(gameWidth - game.cache.getImage('minimap').width/minimapScale, 0, 'minimap');
    minimap.fixedToCamera = true;
    minimap.scale.set(1.0/minimapScale);
    minimap.alpha = 0.8;

    //minimap.setScaleMinMax(1, 1);


    initWebsockets()
}

function update()
{
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
        player.y -= steppingVert;
        sendSocketMove();
    } else if (cursors.down.isDown) {
        player.y += steppingVert;
        sendSocketMove();
    }

    if (cursors.left.isDown) {
        player.x -= steppingHoriz;
        sendSocketMove();
    } else if (cursors.right.isDown) {
        player.x += steppingHoriz;
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

function render()
{
    game.debug.text(game.time.fps || '--', 1, 14, "#00ff00");

    //game.debug.spriteInfo(player, 32, 32);
    //game.debug.cameraInfo(game.camera, 32, 32);

    //game.debug.soundInfo(music, 20, 32);
}




var socket;

function initWebsockets()
{
    var url = 'ws://localhost:3322/ws';
    socket = new WebSocket(url);

    socket.onmessage = onSocketMessage;
    socket.onopen = onSocketConnected;
}

/**
 * @param msg MessageEvent
 */
function onSocketMessage(msg)
{
    var cmd = JSON.parse(msg.data);

    switch (cmd.Type) {
    case 'xy':
        // multiply coords with tile size to scale properly. sprite tiles are always in pixels
        player.x = cmd.X * tileWidth;
        player.y = cmd.Y * tileHeight;
        player.visible = true;

        // floating name over head of player
        var t = game.add.text(0, -16, cmd.Name, { font: "8px Arial", fill: "#ffffff", align: "center" });
        t.anchor.set(0.5);
        player.addChild(t);

        token = cmd.Token;

        // display all from .LocalSpawns
        //console.log(cmd.LocalSpawns);
        for (var i = 0; i < cmd.LocalSpawns.length; i++) {
            var sp = cmd.LocalSpawns[i]
            // console.log(sp)

            // XXX add to game world ...
            game.add.sprite(sp.X * tileWidth, sp.Y * tileHeight, 'phaser');
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
    sendSocketMsg("move " + Math.floor(player.x/tileWidth) + " " + Math.floor(player.y/tileHeight) + " " + token);
}

// Socket connected
function onSocketConnected()
{
    sendSocketMsg("new_player mrcool");

    console.log('Connected to socket server');
}
