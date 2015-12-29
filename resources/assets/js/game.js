var game = new Phaser.Game(
    800,
    600,
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
    game.stage.backgroundColor = '#262f71';  // deep water

    // load world
    game.load.tilemap('island', '/island/full', null, Phaser.Tilemap.TILED_JSON);
    game.load.image('tiles', 'img/sprites/island_tiles.png', 32, 32);

    game.load.image('chunk', 'img/sprites/chunk.png');
    game.load.image('phaser', 'img/sprites/phaser-dude.png');


    game.load.audio('carter', ['audio/dead_feelings.mp3']);
}


var socket;
var map;
var layer;
var cursors;
var player;
var emitter;
var music;

var token;

function create()
{
    music = game.add.audio('carter');
    music.volume = 0.5; // 50%
    music.play();



    // A Tilemap object just holds the data needed to describe the map
    // You can add your own data or manipulate the data (swap tiles around, etc)
    // but in order to display it you need to create a TilemapLayer.
    map = game.add.tilemap('island');

    map.addTilesetImage('island_tiles', 'tiles');

    layer = map.createLayer(0);

    // Basically this sets EVERY SINGLE tile to fully collide on all faces
    map.setCollisionByExclusion([7, 32, 35, 36, 47]);

    layer.resizeWorld();

    cursors = game.input.keyboard.createCursorKeys();

    emitter = game.add.emitter(0, 0, 200);

    emitter.makeParticles('chunk');
    emitter.minRotation = 0;
    emitter.maxRotation = 0;
    emitter.gravity = 150;
    emitter.bounce.setTo(0.5, 0.5);

    player = game.add.sprite(0, 0, 'phaser');
    player.visible = false;
    player.anchor.set(0.5);

    game.physics.enable(player);

    //  Because both our body and our tiles are so tiny,
    //  and the body is moving pretty fast, we need to add
    //  some tile padding to the body. WHat this does
    player.body.tilePadding.set(32, 32);

    game.camera.follow(player);


    initWebsockets()
}

function particleBurst()
{
    emitter.x = player.x;
    emitter.y = player.y;
    emitter.start(true, 2000, null, 1);
}

function update()
{
    game.physics.arcade.collide(player, layer);
    game.physics.arcade.collide(emitter, layer);

    player.body.velocity.x = 0;
    player.body.velocity.y = 0;

    if (cursors.up.isDown) {
        player.body.velocity.y = -200;
        particleBurst();
    } else if (cursors.down.isDown) {
        player.body.velocity.y = 200;
        particleBurst();
    }

    if (cursors.left.isDown) {
        player.body.velocity.x = -200;
        player.scale.x = -1;
        particleBurst();
    } else if (cursors.right.isDown) {
        player.body.velocity.x = 200;
        player.scale.x = 1;
        particleBurst();
    }
}

function render()
{
    game.debug.spriteInfo(player, 32, 32);
    //game.debug.cameraInfo(game.camera, 32, 32);

    //game.debug.soundInfo(music, 20, 32);
}


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
    console.log("<-recv- " + msg.data);

    var cmd = JSON.parse(msg.data);

    switch (cmd.Type) {
    case 'xy':
        // multiply coords with tile size to scale properly. sprite tiles are always in pixels
        player.x = cmd.X * 16;
        player.y = cmd.Y * 16;
        player.visible = true;
        // XXX todo set floating name over head of player

        var t = game.add.text(0, -24, cmd.Name, { font: "14px Arial", fill: "#ffffff", align: "center" });
        t.anchor.set(0.5);
        player.addChild(t);

        token = cmd.Token;
        break;

    case 'pong':
        console.log("pong");
        break;

    default:
        console.log("unknown command from server: " + cmd.Type);
    }
}

function sendSocketMsg(data) {
    socket.send(data);
    console.log("-sent->" + data);
}

// Socket connected
function onSocketConnected()
{
    sendSocketMsg("ping");
    sendSocketMsg("new_player mrcool");

    setInterval(
        function() {
            sendSocketMsg("ping");
        },
        5000
    );

    console.log('Connected to socket server');
}
