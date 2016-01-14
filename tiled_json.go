package rogue

// SEE https://github.com/bjorn/tiled/wiki/JSON-Map-Format

// TiledMapJSON represents the tiled json format, recognized by phaser.io
type TiledMapJSON struct {
	Version     int             `json:"version"`
	Width       int             `json:"width"`
	Height      int             `json:"height"`
	TileWidth   int             `json:"tilewidth"`
	TileHeight  int             `json:"tileheight"`
	Orientation string          `json:"orientation"`
	Layers      []TiledMapLayer `json:"layers"`
	TileSets    []TiledTileSet  `json:"tilesets"`
	// Properties ....  we skipped this
}

// TiledMapLayer represents a map layer
type TiledMapLayer struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Data    []int  `json:"data"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Opacity int    `json:"opacity"`
	Visible bool   `json:"visible"`
	X       int    `json:"x"`
	Y       int    `json:"y"`
}

// TiledTileSet represents a tile set in use on the map
type TiledTileSet struct {
	FirstGid    int    `json:"firstgid"`
	Image       string `json:"image"`
	Name        string `json:"name"`
	ImageHeight int    `json:"imageheight"`
	ImageWidth  int    `json:"imagewidth"`
	Margin      int    `json:"margin"`
	Spacing     int    `json:"spacing"`
	TileHeight  int    `json:"tileheight"`
	TileWidth   int    `json:"tilewidth"`
	// Properties ....  we skipped this
}
