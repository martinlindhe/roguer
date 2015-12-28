var game = new Phaser.Game(
    800,
    600,
    Phaser.CANVAS,
    'phaser-example',
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

    // ground sprites:
    game.load.image('gr_shallow_water', 'img/ground/shallow_water.png');
    game.load.image('gr_beach', 'img/ground/beach.png');
    game.load.image('gr_grass', 'img/ground/grass.png');
    game.load.image('gr_forest', 'img/ground/forest.png');
    game.load.image('gr_hills', 'img/ground/hills.png');
    game.load.image('gr_mountains', 'img/ground/mountains.png');


    // world objects:
    game.load.image('mushroom', 'img/sprites/mushroom.png');
    game.load.image('sonic', 'img/sprites/wabbit.png');
    game.load.image('phaser', 'img/sprites/phaser-dude.png');

    // load world
    game.load.json('island', '/island/full');
}

var cursors;
var logo1;
var logo2;


function create()
{
    var island = game.cache.getJSON('island');

    //  Modify the world and camera bounds
    //game.world.setBounds(-1000, -1000, 1000, 1000);
    game.world.setBounds(0, 0, island.Width * 16, island.Height * 16);

    // height constants
    const deepWater    = 80;
    const shallowWater = 90;
    const beach        = 95;
    const grass        = 150;
    const forest       = 230;
    const hills        = 240;


    for (var y = 0; y < island.Height; y++) {
        for (var x = 0; x < island.Width; x++) {
            var b = island.HeightMap[y][x];
            var sprite = ""

            if (b <= deepWater) {
                continue;
            } else if (b <= shallowWater) {
                sprite = "gr_shallow_water";
            } else if (b <= beach) {
                sprite = "gr_beach";
            } else if (b <= grass) {
                sprite = "gr_grass";
            } else if (b <= forest) {
                sprite = "gr_forest";
            } else if (b <= hills) {
                sprite = "gr_hills";
            } else {
                sprite = "gr_mountains";
            }

            game.add.sprite(x * 16, y * 16, sprite);
        }
    }

    logo1 = game.add.sprite(0, 0, 'phaser');
    logo1.fixedToCamera = true;
    logo1.cameraOffset.setTo(100, 100);

    logo2 = game.add.sprite(0, 0, 'phaser');
    logo2.fixedToCamera = true;
    logo2.cameraOffset.setTo(500, 100);


    game.add.tween(logo2.cameraOffset).to(
        { y: 400 },
        2000,
        Phaser.Easing.Back.InOut,
        true,
        0,
        2000,
        true
    );

    cursors = game.input.keyboard.createCursorKeys();
}

function update()
{
    if (cursors.up.isDown) {
        game.camera.y -= 4;
    } else if (cursors.down.isDown) {
        game.camera.y += 4;
    }

    if (cursors.left.isDown) {
        game.camera.x -= 4;
    } else if (cursors.right.isDown) {
        game.camera.x += 4;
    }
}

function render()
{
    game.debug.cameraInfo(game.camera, 32, 32);
}
