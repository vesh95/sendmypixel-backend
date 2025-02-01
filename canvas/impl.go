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
	wg := sync.WaitGroup{}
	canvas := Area{}
	for y, _ := range canvas {
		go func(area *Area, y int) {
			wg.Add(1)
			canvas[y] = &Row{}
			for x, _ := range canvas[y] {
				canvas[y][x] = Column{
					Color:  "#000000",
					UserId: 0,
				}
			}
			wg.Done()
		}(&canvas, y)
	}

	wg.Wait()

	return &SyncCanvas{
		canvas: canvas,
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
	meta.UserId = dto.UserId
	meta.Color = dto.Color

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
