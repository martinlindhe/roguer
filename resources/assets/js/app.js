import GameState from './GameState.js';

class Game extends Phaser.Game {

	constructor() {
		super(800, 400, Phaser.CANVAS, 'content', null, false, false);
		this.state.add('GameState', GameState, false);
		this.state.start('GameState');
	}
}

new Game();
