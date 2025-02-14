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
	for y := range canvas {
		wg.Add(1)
		go func(area *Area, y int) {
			defer wg.Done()
			canvas[y] = &Row{}
			for x := range canvas[y] {
				canvas[y][x] = Column{
					Color:  "#000000",
					UserId: 0,
				}
			}
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

	meta := &s.canvas[dto.Y][dto.X]
	meta.UserId = dto.UserId
	meta.Color = dto.Color

	return true, nil
}

func (s *SyncCanvas) GetFull() Area {
	var result Area
	wg := sync.WaitGroup{}
	for y := range s.canvas {
		y := y
		wg.Add(1)
		go func() {
			defer wg.Done()
			result[y] = &Row{}
			for x := range s.canvas[y] {
				result[y][x] = Column{
					Color:  s.canvas[y][x].Color,
					UserId: s.canvas[y][x].UserId,
				}
			}
		}()
	}

	wg.Wait()

	return result
}

func (s *SyncCanvas) GetArea(xLeft, yTop, xRight, yBottom int32) (SlicedArea, error) {

	err := validateSize(xLeft, yTop)
	if err != nil {
		return SlicedArea{}, err
	}

	err = validateSize(xRight, yBottom)
	if err != nil {
		return SlicedArea{}, nil
	}

	var result SlicedArea = make(SlicedArea, yBottom-yTop)
	for y, row := range s.canvas[yTop:yBottom] {
		result[y] = make([]Column, xRight-xLeft)
		for x, column := range row[xLeft:xRight] {
			result[y][x] = Column{
				Color:  column.Color,
				UserId: column.UserId,
			}
		}
	}

	return result, nil
}
