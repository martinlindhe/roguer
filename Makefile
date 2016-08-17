.PHONY: views

bench:
	go test -bench=.

run:
	go run cli/server/*

shrink-png:
	find ./public/img -name '*.png' -print0 | xargs -0 -n1 ./shrink-png.sh

# generate sliced up tilesets from source images
tiles:
	tileslicer resources/assets/tilesets/oddball/source/font.png       --out resources/assets/tilesets/oddball/font       --width 8  --height 8
	tileslicer resources/assets/tilesets/oddball/source/characters.png --out resources/assets/tilesets/oddball/characters --width 8  --height 8
	tileslicer resources/assets/tilesets/oddball/source/bosses.png     --out resources/assets/tilesets/oddball/bosses     --width 16 --height 16
	tileslicer resources/assets/tilesets/oddball/source/items.png      --out resources/assets/tilesets/oddball/items      --width 8  --height 8
	tileslicer resources/assets/tilesets/oddball/source/ground.png     --out resources/assets/tilesets/oddball/ground     --width 8  --height 12
	# move 8x12, 8x4 tiles into own folders:
	php move_tiles.php
	# cut the first set of tiles into 8x4:
	tilecutter resources/assets/tilesets/oddball/ground/8x4
	find resources/assets/tilesets/oddball -name '*.png' -print0 | xargs -0 -n1 ./shrink-png.sh

# produce our tilesets
jointiles:
	mkdir -p public/img/tileset/oddball
	tilejoiner resources/assets/tilesets/oddball/ground/8x4  --out resources/assets/tilesets/oddball/ground.png     --tiles-per-row 8;  imgcat resources/assets/tilesets/oddball/ground.png
	tilejoiner resources/assets/tilesets/oddball/ground/8x12 --out resources/assets/tilesets/oddball/ground2.png    --tiles-per-row 8;  imgcat resources/assets/tilesets/oddball/ground2.png
	tilejoiner resources/assets/tilesets/oddball/items       --out resources/assets/tilesets/oddball/items.png      --tiles-per-row 8;  imgcat resources/assets/tilesets/oddball/items.png
	tilejoiner resources/assets/tilesets/oddball/font        --out resources/assets/tilesets/oddball/font.png       --tiles-per-row 16; imgcat resources/assets/tilesets/oddball/font.png
	tilejoiner resources/assets/tilesets/oddball/characters  --out resources/assets/tilesets/oddball/characters.png --tiles-per-row 16; imgcat resources/assets/tilesets/oddball/characters.png
	tilejoiner resources/assets/tilesets/oddball/bosses      --out resources/assets/tilesets/oddball/bosses.png     --tiles-per-row 8;  imgcat resources/assets/tilesets/oddball/bosses.png
	cp resources/assets/tilesets/oddball/*.png public/img/tileset/oddball/
	tilejoiner resources/assets/tilesets/ui/buttons          --out public/img/tileset/ui/buttons.png                --tiles-per-row 1;  imgcat public/img/tileset/ui/buttons.png
	find public/img/tileset -name '*.png' -print0 | xargs -0 -n1 ./shrink-png.sh

validate-yaml:
	find . -name '*.yml' -not -path "./node_modules/*" -print0 | xargs -0 -n1 validyaml

validate-json:
	find . -name '*.json' -not -path "./node_modules/*" -print0 | xargs -0 -n1 validjson

validate-js:
	node_modules/jshint/bin/jshint resources/assets/js

watch:
	# BUG: watch dont run sass.. elixir 4.2
	node_modules/gulp/bin/gulp.js
	node_modules/gulp/bin/gulp.js watch

deps:
	npm install

build-custom-phaser:
	cd node_modules/phaser && npm install && grunt --split true
