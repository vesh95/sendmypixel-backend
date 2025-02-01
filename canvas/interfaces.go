package canvas

type Row [size]*Column
type Area [size]Row
type SlicedArea [][]Column

type PixelDto struct {
	Y      int32  `json:"y"`
	X      int32  `json:"x"`
	Color  string `json:"color"`
	UserId int64  `json:"user_id"`
}

// Canvas
// TODO segregate interface
type Canvas interface {
	GetPixel(x, y int32) (PixelDto, error)
	SetPixel(dto PixelDto) (bool, error)
	GetFull() SlicedArea
	GetArea(xLeft, yTop, xBottom, yBottom int32) (SlicedArea, error)
}
