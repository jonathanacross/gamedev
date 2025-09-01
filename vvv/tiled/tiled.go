package tiled

// --- Intermediate Tiled JSON structures to help with unmarshalling ---

type tiledMap struct {
	Width            int            `json:"width"`
	Height           int            `json:"height"`
	TileWidth        int            `json:"tilewidth"`
	TileHeight       int            `json:"tileheight"`
	Layers           []tiledLayer   `json:"layers"`
	Tilesets         []tiledTileset `json:"tilesets"`
	CompressionLevel int            `json:"compressionlevel"`
}

type tiledLayer struct {
	Name    string        `json:"name"`
	Type    string        `json:"type"`
	Width   int           `json:"width"`
	Height  int           `json:"height"`
	Data    []int         `json:"data"`
	Objects []tiledObject `json:"objects"`
}

type tiledTileset struct {
	FirstGID   int         `json:"firstgid"`
	Source     string      `json:"source"`
	Image      string      `json:"image"`
	Tiles      []tiledTile `json:"tiles"`
	Name       string      `json:"name"`
	TileWidth  int         `json:"tilewidth"`
	TileHeight int         `json:"tileheight"`
	TileCount  int         `json:"tilecount"`
	Columns    int         `json:"columns"`
}

type tiledTile struct {
	ID          int              `json:"id"`
	Image       string           `json:"image"`
	ImageWidth  int              `json:"imagewidth"`
	ImageHeight int              `json:"imageheight"`
	Properties  []tiledProperty  `json:"properties"`
	Type        string           `json:"type"`
	ObjectGroup tiledObjectGroup `json:"objectgroup"`
}

type tiledProperty struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type tiledObject struct {
	ID         int             `json:"id"`
	Name       string          `json:"name"`
	Type       string          `json:"type"`
	X          float64         `json:"x"`
	Y          float64         `json:"y"`
	Width      float64         `json:"width"`
	Height     float64         `json:"height"`
	Rotation   float64         `json:"rotation"`
	GID        int             `json:"gid"`
	Properties []tiledProperty `json:"properties"`
}

type tiledObjectGroup struct {
	Name    string        `json:"name"`
	Objects []tiledObject `json:"objects"`
}
