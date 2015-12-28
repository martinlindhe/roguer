import io from 'socket.io-client';

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
    socket = io('ws://localhost', {path: '/ws'});

    // Socket connection successful
    socket.on('connect', onSocketConnected);

    // Socket disconnection
    socket.on('disconnect', onSocketDisconnect);

    // socket.on('event', function(data){});

    // New player message received
    socket.on('new player', onNewPlayer);

    // Player move message received
    socket.on('move player', onMovePlayer);

    // Player removed message received
    socket.on('remove player', onRemovePlayer);
}

// Socket connected
function onSocketConnected()
{
    console.log('Connected to socket server');

    // Send local player data to the game server
    socket.emit('new player', { x: player.x, y: player.y })
}

// Socket disconnected
function onSocketDisconnect()
{
    console.log('Disconnected from socket server');
}

// New player
function onNewPlayer(data)
{
    console.log('New player connected:', data.id);

    // Add new player to the remote players array
    // enemies.push(new RemotePlayer(data.id, game, player, data.x, data.y));
}

// Move player
function onMovePlayer(data)
{
    var movePlayer = playerById(data.id);

    // Player not found
    if (!movePlayer) {
        console.log('Player not found: ', data.id);
        return;
    }

    // Update player position
    movePlayer.player.x = data.x;
    movePlayer.player.y = data.y;
}

// Remove player
function onRemovePlayer(data)
{
    var removePlayer = playerById(data.id);

    // Player not found
    if (!removePlayer) {
        console.log('Player not found: ', data.id);
        return;
    }

    removePlayer.player.kill();

    // Remove player from array
    enemies.splice(enemies.indexOf(removePlayer), 1);
}
