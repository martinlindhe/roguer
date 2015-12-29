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
	# move 8x12, 8x4 tiles into own folders:
	php move_tiles.php
	# cut the 8x12 tiles into 8x4:
	go run cli/tilecutter/tilecutter.go
	# pngcrush:
	find resources/assets/tilesets/oddball -name '*.png' -print0 | xargs -0 -n1 ./shrink-png.sh
