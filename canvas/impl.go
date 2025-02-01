package canvas

import (
	"sync"
)

type SyncCanvas struct {
	canvas Area
	mux    *sync.Mutex
}

// NewSyncCanvas is constructor of canvas object
func NewSyncCanvas() *SyncCanvas {
	return &SyncCanvas{
		canvas: Area{},
		mux:    &sync.Mutex{},
	}
}

func (s *SyncCanvas) GetPixel(x, y int32) (PixelDto, error) {
	if err := validateSize(x, y); err != nil {
		return PixelDto{}, err
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	meta := s.canvas[y][x]
	return PixelDto{X: x, Y: y, Color: meta.Color, UserId: meta.UserId}, nil
}

func (s *SyncCanvas) SetPixel(dto PixelDto) (bool, error) {
	if err := validateDto(dto); err != nil {
		return false, err
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	meta := s.canvas[dto.Y][dto.X]
	if meta == nil {
		s.canvas[dto.Y][dto.X] = &Column{
			Color:  dto.Color,
			UserId: dto.UserId,
		}
	} else {
		meta.UserId = dto.UserId
		meta.Color = dto.Color
	}

	return true, nil
}

func (s *SyncCanvas) GetFull() SlicedArea {
	var result SlicedArea
	for _, row := range s.canvas {
		var columns []Column
		for _, pixelMeta := range row {
			columns = append(columns, Column{
				Color:  pixelMeta.Color,
				UserId: pixelMeta.UserId,
			})
		}
		result = append(result, columns)
	}

	return result
}

func (s *SyncCanvas) GetArea(xLeft, yTop, xRight, yBottom int32) (SlicedArea, error) {
	panic("implements me")
	//err := validateSize(xLeft, yTop)
	//if err != nil {
	//	return SlicedArea{}, err
	//}
	//err = validateSize(xRight, yBottom)
	//if err != nil {
	//	return SlicedArea{}, nil
	//}
	//
	//var result SlicedArea
	//for _, row := range s.canvas[yTop:yBottom] {
	//	var columns := make([]Column)
	//	for _, column := range row {
	//		copy(columns, row)
	//	}
	//	result = append(result, make([]Column, len(row)))
	//}
	return SlicedArea{}, nil
}
