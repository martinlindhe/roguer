.PHONY: views

bench:
	go test -bench=.

run:
	go run cli/app/*

views:
	gorazor views views

shrink-png:
	find ./public/img -name '*.png' -print0 | xargs -0 -n1 ./shrink-png.sh

# generates sliced up tilesets (result is already in git)
tiles:
	go run cli/tileslicer/tileslicer.go resources/assets/tilesets/oddball/font.png resources/assets/tilesets/oddball/font 8 8
	go run cli/tileslicer/tileslicer.go resources/assets/tilesets/oddball/characters.png resources/assets/tilesets/oddball/characters 8 8
	go run cli/tileslicer/tileslicer.go resources/assets/tilesets/oddball/bosses.png resources/assets/tilesets/oddball/bosses 16 16
	go run cli/tileslicer/tileslicer.go resources/assets/tilesets/oddball/items.png resources/assets/tilesets/oddball/items 8 8
	go run cli/tileslicer/tileslicer.go resources/assets/tilesets/oddball/tiles.png resources/assets/tilesets/oddball/tiles 8 12
	# tiles 0-73 is half height (4x12) how on earth do i move those with a simple command
	php move_tiles.php

	#find resources/assets/tilesets/oddball/ -name '*.png' -print0 | xargs -0 -n1 ./shrink-png.sh
