package rogue

// PhaserTileMap represents the tiled json format, recognized by phaser.io
type PhaserTileMap struct {
	Version     int               `json:"version"`
	Width       int               `json:"width"`
	Height      int               `json:"height"`
	TileWidth   int               `json:"tilewidth"`
	TileHeight  int               `json:"tileheight"`
	Orientation string            `json:"orientation"`
	Layers      []PhaserTileLayer `json:"layers"`
	TileSets    []PhaserTileSet   `json:"tilesets"`
	// Properties ....  we skipped this
}

type PhaserTileLayer struct {
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

type PhaserTileSet struct {
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
