.PHONY: views

bench:
	go test -bench=.

run:
	go run cli/server/*

shrink-png:
	find ./public/img -name '*.png' -print0 | xargs -0 -n1 ./shrink-png.sh

# generate sliced up tilesets from source images
tiles:
	go run cli/tileslicer/tileslicer.go resources/assets/tilesets/oddball/source/font.png       resources/assets/tilesets/oddball/font 8 8
	go run cli/tileslicer/tileslicer.go resources/assets/tilesets/oddball/source/characters.png resources/assets/tilesets/oddball/characters 8 8
	go run cli/tileslicer/tileslicer.go resources/assets/tilesets/oddball/source/bosses.png     resources/assets/tilesets/oddball/bosses 16 16
	go run cli/tileslicer/tileslicer.go resources/assets/tilesets/oddball/source/items.png      resources/assets/tilesets/oddball/items 8 8
	go run cli/tileslicer/tileslicer.go resources/assets/tilesets/oddball/source/ground.png     resources/assets/tilesets/oddball/ground 8 12
	# move 8x12, 8x4 tiles into own folders:
	php move_tiles.php
	# cut the first set of tiles into 8x4:
	go run cli/tilecutter/tilecutter.go resources/assets/tilesets/oddball/ground/8x4
	find resources/assets/tilesets/oddball -name '*.png' -print0 | xargs -0 -n1 ./shrink-png.sh

# produce our tilesets
jointiles:
	mkdir -p public/img/tileset/oddball
	go run cli/tilejoiner/tilejoiner.go resources/assets/tilesets/oddball/ground/8x4  resources/assets/tilesets/oddball/ground.png 8;      imgcat resources/assets/tilesets/oddball/ground.png
	go run cli/tilejoiner/tilejoiner.go resources/assets/tilesets/oddball/ground/8x12 resources/assets/tilesets/oddball/ground2.png 8;     imgcat resources/assets/tilesets/oddball/ground2.png
	go run cli/tilejoiner/tilejoiner.go resources/assets/tilesets/oddball/items       resources/assets/tilesets/oddball/items.png 8;       imgcat resources/assets/tilesets/oddball/items.png
	go run cli/tilejoiner/tilejoiner.go resources/assets/tilesets/oddball/font        resources/assets/tilesets/oddball/font.png 16;       imgcat resources/assets/tilesets/oddball/font.png
	go run cli/tilejoiner/tilejoiner.go resources/assets/tilesets/oddball/characters  resources/assets/tilesets/oddball/characters.png 16; imgcat resources/assets/tilesets/oddball/characters.png
	go run cli/tilejoiner/tilejoiner.go resources/assets/tilesets/oddball/bosses      resources/assets/tilesets/oddball/bosses.png 8;      imgcat resources/assets/tilesets/oddball/bosses.png
	cp resources/assets/tilesets/oddball/*.png public/img/tileset/oddball/
	go run cli/tilejoiner/tilejoiner.go resources/assets/tilesets/ui/buttons public/img/tileset/ui/buttons.png 1;                          imgcat public/img/tileset/ui/buttons.png
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
