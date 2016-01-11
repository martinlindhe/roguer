
import GameState from './GameState.js';
/*
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
*/




class Game extends Phaser.Game {

	constructor() {
		super(800, 400, Phaser.CANVAS, 'content', null, false, false);
		this.state.add('GameState', GameState, false);
		this.state.start('GameState');
	}
}

new Game();
