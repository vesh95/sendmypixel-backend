package canvas

import "sync"

type PixelDto struct {
	Row   int    `json:"row"`
	Col   int    `json:"col"`
	Color string `json:"color"`
}

type Canvas [100][100]string

type SyncCanvas struct {
	canvas Canvas
	mux    *sync.Mutex
}

func CreateNewState() SyncCanvas {
	mux := &sync.Mutex{}
	cnv := Canvas{}
	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			cnv[i][j] = "white"
		}
	}

	return SyncCanvas{mux: mux, canvas: cnv}
}

func (s *SyncCanvas) SetPixel(dto PixelDto) (PixelDto, error) {
	s.mux.Lock()
	s.canvas[dto.Row][dto.Col] = dto.Color
	s.mux.Unlock()

	return dto, nil
}

func (s *SyncCanvas) GetState() Canvas {
	result := Canvas{}
	s.mux.Lock()
	for i, row := range s.canvas {
		for j, col := range row {
			result[i][j] = col
		}
	}
	s.mux.Unlock()

	return result
}
