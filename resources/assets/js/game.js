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
    // XXX post to new player backend, get initial coordinates




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
var sprite;
var emitter;
var music;

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

    sprite = game.add.sprite(300, 90, 'phaser');
    sprite.anchor.set(0.5);

    game.physics.enable(sprite);

    //  Because both our body and our tiles are so tiny,
    //  and the body is moving pretty fast, we need to add
    //  some tile padding to the body. WHat this does
    sprite.body.tilePadding.set(32, 32);

    game.camera.follow(sprite);


    initWebsockets()
}

function particleBurst()
{
    emitter.x = sprite.x;
    emitter.y = sprite.y;
    emitter.start(true, 2000, null, 1);
}

function update()
{
    game.physics.arcade.collide(sprite, layer);
    game.physics.arcade.collide(emitter, layer);

    sprite.body.velocity.x = 0;
    sprite.body.velocity.y = 0;

    if (cursors.up.isDown) {
        sprite.body.velocity.y = -200;
        particleBurst();
    } else if (cursors.down.isDown) {
        sprite.body.velocity.y = 200;
        particleBurst();
    }

    if (cursors.left.isDown) {
        sprite.body.velocity.x = -200;
        sprite.scale.x = -1;
        particleBurst();
    } else if (cursors.right.isDown) {
        sprite.body.velocity.x = 200;
        sprite.scale.x = 1;
        particleBurst();
    }
}

function render()
{
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

function onSocketMessage(msg)
{
    console.log("<-- " + msg);
}

// Socket connected
function onSocketConnected()
{
    var send = function(data) {
        console.log("-->" + data);
        socket.send(data);
    }

    setInterval(
        function(){ send("ping") },
        1000
    );


    console.log('Connected to socket server');
}
