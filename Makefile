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
	go run cli/tileslicer/tileslicer.go resources/assets/tilesets/oddball/font.png       resources/assets/tilesets/oddball/font 8 8
	go run cli/tileslicer/tileslicer.go resources/assets/tilesets/oddball/characters.png resources/assets/tilesets/oddball/characters 8 8
	go run cli/tileslicer/tileslicer.go resources/assets/tilesets/oddball/bosses.png     resources/assets/tilesets/oddball/bosses 16 16
	go run cli/tileslicer/tileslicer.go resources/assets/tilesets/oddball/items.png      resources/assets/tilesets/oddball/items 8 8
	go run cli/tileslicer/tileslicer.go resources/assets/tilesets/oddball/tiles.png      resources/assets/tilesets/oddball/tiles 8 12
	# move 8x12, 8x4 tiles into own folders:
	php move_tiles.php
	# cut the 8x12 tiles into 8x4:
	go run cli/tilecutter/tilecutter.go
	find resources/assets/tilesets/oddball -name '*.png' -print0 | xargs -0 -n1 ./shrink-png.sh


jointiles:
	mkdir -p public/img/tileset/oddball
	go run cli/tilejoiner/tilejoiner.go resources/assets/tilesets/oddball/tiles/8x4  public/img/tileset/oddball/tiles.png 8; imgcat public/img/tileset/oddball/tiles.png
	go run cli/tilejoiner/tilejoiner.go resources/assets/tilesets/oddball/tiles/8x12 public/img/tileset/oddball/tiles-big.png 8; imgcat public/img/tileset/oddball/tiles-big.png
	go run cli/tilejoiner/tilejoiner.go resources/assets/tilesets/oddball/items      public/img/tileset/oddball/items.png 8; imgcat public/img/tileset/oddball/items.png
	go run cli/tilejoiner/tilejoiner.go resources/assets/tilesets/oddball/font       public/img/tileset/oddball/font.png 16; imgcat public/img/tileset/oddball/font.png
	go run cli/tilejoiner/tilejoiner.go resources/assets/tilesets/oddball/characters public/img/tileset/oddball/characters.png 16; imgcat public/img/tileset/oddball/characters.png
	go run cli/tilejoiner/tilejoiner.go resources/assets/tilesets/oddball/bosses     public/img/tileset/oddball/bosses.png 8; imgcat public/img/tileset/oddball/bosses.png
	find public/img/tileset/oddball -name '*.png' -print0 | xargs -0 -n1 ./shrink-png.sh

validate-yaml:
	find . -name '*.yml' -not -path "./node_modules/*" -print0 | xargs -0 -n1 validyaml

validate-json:
	find . -name '*.json' -not -path "./node_modules/*" -print0 | xargs -0 -n1 validjson
