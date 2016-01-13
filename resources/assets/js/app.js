import {GameState} from './GameState.js';

class Game extends Phaser.Game {

	constructor() {
		var innerWidth = window.innerWidth;
		var innerHeight = window.innerHeight;

		super(innerWidth, innerHeight, Phaser.CANVAS, '', null, false, false);
		this.state.add('GameState', GameState, false);
		this.state.start('GameState');
	}
}

new Game();
