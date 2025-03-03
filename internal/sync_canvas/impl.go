package sync_canvas

import (
	"backend/pkg/canvas"
	"sync"
)

type SyncCanvas struct {
	canvas canvas.Area
	mux    *sync.Mutex
}

// NewSyncCanvas is constructor of canvas object
func NewSyncCanvas() *SyncCanvas {
	wg := sync.WaitGroup{}
	canv := canvas.Area{}
	for y := range canv {
		wg.Add(1)
		go func(area *canvas.Area, y int) {
			defer wg.Done()
			canv[y] = &canvas.Row{}
			for x := range canv[y] {
				canv[y][x] = canvas.Column{
					Color:  "#000000",
					UserId: 0,
				}
			}
		}(&canv, y)
	}

	wg.Wait()

	return &SyncCanvas{
		canvas: canv,
		mux:    &sync.Mutex{},
	}
}

func (s *SyncCanvas) GetPixel(x, y int32) (canvas.PixelDto, error) {
	if err := canvas.ValidateSize(x, y); err != nil {
		return canvas.PixelDto{}, err
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	meta := s.canvas[y][x]
	return canvas.PixelDto{X: x, Y: y, Color: meta.Color, UserId: meta.UserId}, nil
}

func (s *SyncCanvas) SetPixel(dto canvas.PixelDto) (bool, error) {
	if err := canvas.ValidateDto(dto); err != nil {
		return false, err
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	meta := &s.canvas[dto.Y][dto.X]
	meta.UserId = dto.UserId
	meta.Color = dto.Color

	return true, nil
}

func (s *SyncCanvas) GetFull() canvas.Area {
	var result canvas.Area
	wg := sync.WaitGroup{}
	for y := range s.canvas {
		y := y
		wg.Add(1)
		go func() {
			defer wg.Done()
			result[y] = &canvas.Row{}
			for x := range s.canvas[y] {
				result[y][x] = canvas.Column{
					Color:  s.canvas[y][x].Color,
					UserId: s.canvas[y][x].UserId,
				}
			}
		}()
	}

	wg.Wait()

	return result
}

func (s *SyncCanvas) GetArea(xLeft, yTop, xRight, yBottom int32) (canvas.SlicedArea, error) {
	err := canvas.ValidateSize(xLeft, yTop)
	if err != nil {
		return canvas.SlicedArea{}, err
	}

	err = canvas.ValidateSize(xRight, yBottom)
	if err != nil {
		return canvas.SlicedArea{}, err
	}

	var result canvas.SlicedArea = make(canvas.SlicedArea, yBottom-yTop)
	for y, row := range s.canvas[yTop:yBottom] {
		result[y] = make([]canvas.Column, xRight-xLeft)
		for x, column := range row[xLeft:xRight] {
			result[y][x] = canvas.Column{
				Color:  column.Color,
				UserId: column.UserId,
			}
		}
	}

	return result, nil
}
