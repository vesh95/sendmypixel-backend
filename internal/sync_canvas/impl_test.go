package sync_canvas

import (
	"backend/pkg/canvas"
	"errors"
	"log"
	"strings"
	"testing"
)

func TestSyncCanvas_SetPixel(t *testing.T) {
	fixture := NewSyncCanvas()

	cases := []struct {
		name  string
		dto   canvas.PixelDto
		err   error
		exptd bool
	}{
		{
			name: "Successfully set",
			dto: canvas.PixelDto{
				Y:      0,
				X:      0,
				Color:  "#FAFAFA",
				UserId: 0,
			},
			err:   nil,
			exptd: true,
		},
		{
			name: "Set out of bound pixel by Y",
			dto: canvas.PixelDto{
				Y:      canvas.Size,
				X:      0,
				Color:  "#F0F0F0",
				UserId: 0,
			},
			err:   canvas.ErrorOutOfBounds,
			exptd: false,
		},
		{
			name: "Set out of bound pixel by X",
			dto: canvas.PixelDto{
				Y:      0,
				X:      canvas.Size,
				Color:  "#F0F0F0",
				UserId: 0,
			},
			err:   canvas.ErrorOutOfBounds,
			exptd: false,
		},
		{
			name: "Set out of bound pixel by Y with under 0 value",
			dto: canvas.PixelDto{
				Y:      -1,
				X:      0,
				Color:  "#F0F0F0",
				UserId: 0,
			},
			err:   canvas.ErrorOutOfBounds,
			exptd: false,
		},
		{
			name: "Set out of bound pixel by X with under 0 value",
			dto: canvas.PixelDto{
				Y:      0,
				X:      -1,
				Color:  "#F0F0F0",
				UserId: 0,
			},
			err:   canvas.ErrorOutOfBounds,
			exptd: false,
		},
		{
			name: "Set invalid color",
			dto: canvas.PixelDto{
				Y:      0,
				X:      0,
				Color:  "yellow",
				UserId: 0,
			},
			err:   canvas.ErrorColorInvalid,
			exptd: false,
		},
		{
			name: "Successfully reset",
			dto: canvas.PixelDto{
				Y:      0,
				X:      0,
				Color:  "#FAFAF3",
				UserId: 1,
			},
			err:   nil,
			exptd: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ok, err := fixture.SetPixel(c.dto)
			if c.err != nil && !errors.Is(c.err, err) {
				t.Errorf("Error check: %v was excepted, %v was received", c.err.Error(), err.Error())
			} else if c.err == nil && err != nil {
				t.Fatalf("Unexcepted error: %T", err)
			}

			if c.exptd != ok {
				t.Errorf("%v was expected, %v was received.", c.exptd, ok)
			}
		})
	}
}

func TestSyncCanvas_GetPixel(t *testing.T) {
	fixture := NewSyncCanvas()
	fixture.canvas[0][0] = canvas.Column{
		Color:  "#FAFAFA",
		UserId: 1,
	}

	cases := []struct {
		name string
		args struct {
			X, Y int32
		}
		err    error
		expctd canvas.PixelDto
	}{
		{
			name: "Successfully getting pixel",
			args: struct{ X, Y int32 }{X: 0, Y: 0},
			err:  nil,
			expctd: canvas.PixelDto{
				Y:      0,
				X:      0,
				Color:  "#FAFAFA",
				UserId: 1,
			},
		},
		{
			name: "Successfully getting uninitialized pixel",
			args: struct{ X, Y int32 }{X: 1, Y: 0},
			err:  nil,
			expctd: canvas.PixelDto{
				Y:      0,
				X:      1,
				Color:  "#000000",
				UserId: 0,
			},
		},
		{
			name: "Out of bounds X",
			args: struct{ X, Y int32 }{X: canvas.Size, Y: 0},
			err:  canvas.ErrorOutOfBounds,
			expctd: canvas.PixelDto{
				Y:      0,
				X:      0,
				Color:  "",
				UserId: 0,
			},
		},
		{
			name: "Out of bounds Y",
			args: struct{ X, Y int32 }{X: 0, Y: canvas.Size},
			err:  canvas.ErrorOutOfBounds,
			expctd: canvas.PixelDto{
				Y:      0,
				X:      0,
				Color:  "",
				UserId: 0,
			},
		},
		{
			name: "Out of bounds X with under 0 value",
			args: struct{ X, Y int32 }{X: -1, Y: 0},
			err:  canvas.ErrorOutOfBounds,
			expctd: canvas.PixelDto{
				Y:      0,
				X:      0,
				Color:  "",
				UserId: 0,
			},
		},
		{
			name: "Out of bounds Y with under 0 value",
			args: struct{ X, Y int32 }{X: 0, Y: -1},
			err:  canvas.ErrorOutOfBounds,
			expctd: canvas.PixelDto{
				Y:      0,
				X:      0,
				Color:  "",
				UserId: 0,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			pixel, err := fixture.GetPixel(c.args.X, c.args.Y)

			if c.err != nil && !errors.Is(err, c.err) {
				t.Errorf("Error check: %v was excepted, %v was received", c.err.Error(), err.Error())
			} else if c.err == nil && err != nil {
				t.Fatalf("Unexcepted error: %T", err)
			}

			if !comparePixelDto(pixel, c.expctd) {
				t.Errorf("%v was excepted, %v was received", pixel, c.expctd)
			}
		})
	}
}

type expectedColumnsList []struct {
	x, y      int32
	expctdCol *canvas.Column
}

func TestSyncCanvas_GetFull(t *testing.T) {
	fixture := NewSyncCanvas()
	fixture.canvas[0][0] = canvas.Column{
		Color:  "#FAFAFA",
		UserId: 1,
	}

	cases := []struct {
		name   string
		expctd expectedColumnsList
	}{
		{
			name: "Get successfully",
			expctd: expectedColumnsList{
				{
					x: 0,
					y: 0,
					expctdCol: &canvas.Column{
						Color:  "#FAFAFA",
						UserId: 1,
					},
				},
				{
					x: 1,
					y: 1,
					expctdCol: &canvas.Column{
						Color:  "#000000",
						UserId: 0,
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			area := fixture.GetFull()
			for _, col := range c.expctd {
				colorConsist := col.expctdCol.Color == area[col.y][col.x].Color
				userIdConsist := col.expctdCol.UserId == area[col.y][col.x].UserId
				if !colorConsist {
					t.Errorf("Colors not consists %s != %s", col.expctdCol.Color, area[col.y][col.x].Color)
				}

				if !userIdConsist {
					t.Errorf("UserId not consists %d != %d", col.expctdCol.UserId, area[col.y][col.x].UserId)
				}
			}
		})
	}
}

type areaTestCase struct {
	name           string
	x1, x2, y1, y2 int32
	err            error
	trsPtsMatch    bool
}

var tracePointColor string = "#FAFAFA"

func TestSyncCanvas_GetArea(t *testing.T) {

	fixture := NewSyncCanvas()
	fixture.canvas[0][0] = canvas.Column{
		Color:  tracePointColor,
		UserId: 1,
	}
	fixture.canvas[4][4] = canvas.Column{
		Color:  tracePointColor,
		UserId: 1,
	}
	fixture.canvas[0][4] = canvas.Column{
		Color:  tracePointColor,
		UserId: 1,
	}
	fixture.canvas[4][0] = canvas.Column{
		Color:  tracePointColor,
		UserId: 1,
	}

	cases := []areaTestCase{
		{
			name:        "Get successfully",
			x1:          0,
			x2:          4,
			y1:          0,
			y2:          4,
			err:         nil,
			trsPtsMatch: true,
		},
		{
			name:        "Get out of range x1",
			x1:          canvas.Size + 1,
			x2:          4,
			y1:          0,
			y2:          4,
			err:         canvas.ErrorOutOfBounds,
			trsPtsMatch: false,
		},
		{
			name:        "Get out of range x2",
			x1:          0,
			x2:          canvas.Size + 1,
			y1:          0,
			y2:          4,
			err:         canvas.ErrorOutOfBounds,
			trsPtsMatch: false,
		},
		{
			name:        "Get out of range y1",
			x1:          0,
			x2:          4,
			y1:          canvas.Size + 1,
			y2:          4,
			err:         canvas.ErrorOutOfBounds,
			trsPtsMatch: false,
		},
		{
			name:        "Get out of range x2",
			x1:          0,
			x2:          4,
			y1:          0,
			y2:          canvas.Size + 1,
			err:         canvas.ErrorOutOfBounds,
			trsPtsMatch: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			area, err := fixture.GetArea(c.x1, c.y1, c.x2, c.y2)

			if err != nil {
				log.Printf("%T\n", err)
				if !errors.Is(c.err, err) {
					t.Errorf("Exceped: %s, got: %s", err, c.err)
				}
			} else {
				if len(area) != int(c.y2-c.y1) {
					t.Errorf("Unexcepted area rows count %d. Excepted: %d", len(area), int(c.y2-c.y1))
				}

				if c.trsPtsMatch == checkTracePoints(area) {
					t.Error("TracePoints not consist")
				}
			}
		})
	}
}

func checkTracePoints(area canvas.SlicedArea) bool {
	rows := len(area) - 1
	cols := len(area[0]) - 1
	return area[0][0].Color == tracePointColor &&
		area[0][cols].Color == tracePointColor &&
		area[rows][0].Color == tracePointColor &&
		area[rows][cols].Color == tracePointColor
}

func comparePixelDto(value canvas.PixelDto, exceptedValue canvas.PixelDto) bool {
	return value.Y == exceptedValue.Y &&
		value.X == exceptedValue.X &&
		strings.Compare(value.Color, exceptedValue.Color) == 0 &&
		value.UserId == exceptedValue.UserId
}
