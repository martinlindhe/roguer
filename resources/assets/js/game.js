
import GameState from './GameState.js';

var game = new Phaser.Game(
    800,
    400,
    Phaser.CANVAS,
    'game',
    {},
    false,  // transparent
    false   // antialias
);
game.state.add('game', GameState, true);
