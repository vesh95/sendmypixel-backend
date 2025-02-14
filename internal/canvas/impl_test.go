package canvas

import (
	"errors"
	"strings"
	"testing"
)

func TestSyncCanvas_SetPixel(t *testing.T) {
	fixture := NewSyncCanvas()

	cases := []struct {
		name  string
		dto   PixelDto
		err   error
		exptd bool
	}{
		{
			name: "Successfully set",
			dto: PixelDto{
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
			dto: PixelDto{
				Y:      size,
				X:      0,
				Color:  "#F0F0F0",
				UserId: 0,
			},
			err:   OutOfBounds{x: 0, y: size},
			exptd: false,
		},
		{
			name: "Set out of bound pixel by X",
			dto: PixelDto{
				Y:      0,
				X:      size,
				Color:  "#F0F0F0",
				UserId: 0,
			},
			err:   OutOfBounds{x: size, y: 0},
			exptd: false,
		},
		{
			name: "Set out of bound pixel by Y with under 0 value",
			dto: PixelDto{
				Y:      -1,
				X:      0,
				Color:  "#F0F0F0",
				UserId: 0,
			},
			err:   OutOfBounds{x: 0, y: -1},
			exptd: false,
		},
		{
			name: "Set out of bound pixel by X with under 0 value",
			dto: PixelDto{
				Y:      0,
				X:      -1,
				Color:  "#F0F0F0",
				UserId: 0,
			},
			err:   OutOfBounds{x: -1, y: 0},
			exptd: false,
		},
		{
			name: "Set invalid color",
			dto: PixelDto{
				Y:      0,
				X:      0,
				Color:  "yellow",
				UserId: 0,
			},
			err:   ColorInvalid{color: "yellow"},
			exptd: false,
		},
		{
			name: "Successfully reset",
			dto: PixelDto{
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

			if c.err != nil && errors.Is(err, c.err) {
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
	fixture.canvas[0][0] = Column{
		Color:  "#FAFAFA",
		UserId: 1,
	}

	cases := []struct {
		name string
		args struct {
			X, Y int32
		}
		err    error
		expctd PixelDto
	}{
		{
			name: "Successfully getting pixel",
			args: struct{ X, Y int32 }{X: 0, Y: 0},
			err:  nil,
			expctd: PixelDto{
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
			expctd: PixelDto{
				Y:      0,
				X:      1,
				Color:  "#000000",
				UserId: 0,
			},
		},
		{
			name: "Out of bounds X",
			args: struct{ X, Y int32 }{X: size, Y: 0},
			err:  OutOfBounds{size, 0},
			expctd: PixelDto{
				Y:      0,
				X:      0,
				Color:  "",
				UserId: 0,
			},
		},
		{
			name: "Out of bounds Y",
			args: struct{ X, Y int32 }{X: 0, Y: size},
			err:  OutOfBounds{0, size},
			expctd: PixelDto{
				Y:      0,
				X:      0,
				Color:  "",
				UserId: 0,
			},
		},
		{
			name: "Out of bounds X with under 0 value",
			args: struct{ X, Y int32 }{X: -1, Y: 0},
			err:  OutOfBounds{-1, 0},
			expctd: PixelDto{
				Y:      0,
				X:      0,
				Color:  "",
				UserId: 0,
			},
		},
		{
			name: "Out of bounds Y with under 0 value",
			args: struct{ X, Y int32 }{X: 0, Y: -1},
			err:  OutOfBounds{0, -1},
			expctd: PixelDto{
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

			if c.err != nil && errors.Is(err, c.err) {
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
	expctdCol *Column
}

func TestSyncCanvas_GetFull(t *testing.T) {
	fixture := NewSyncCanvas()
	fixture.canvas[0][0] = Column{
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
					expctdCol: &Column{
						Color:  "#FAFAFA",
						UserId: 1,
					},
				},
				{
					x: 1,
					y: 1,
					expctdCol: &Column{
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
	fixture.canvas[0][0] = Column{
		Color:  tracePointColor,
		UserId: 1,
	}
	fixture.canvas[4][4] = Column{
		Color:  tracePointColor,
		UserId: 1,
	}
	fixture.canvas[0][4] = Column{
		Color:  tracePointColor,
		UserId: 1,
	}
	fixture.canvas[4][0] = Column{
		Color:  tracePointColor,
		UserId: 1,
	}

	cases := []areaTestCase{
		{
			name:        "Get succefully",
			x1:          0,
			x2:          4,
			y1:          0,
			y2:          4,
			err:         nil,
			trsPtsMatch: true,
		},

		{
			name:        "Get out of range",
			x1:          0,
			x2:          size + 1,
			y1:          0,
			y2:          4,
			err:         OutOfBounds{x: size, y: 4},
			trsPtsMatch: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			area, err := fixture.GetArea(c.x1, c.y1, c.x2, c.y2)
			t.Log(err)
			if !errors.Is(err, c.err) {
				t.Errorf("Exceped %s, recieved %s", c.err, err)
				return
			}

			if len(area) != int(c.y2-c.y1) {
				t.Errorf("Unexcepted area rows count %d. Excepted: %d", len(area), int(c.y2-c.y1))
			}

			if c.trsPtsMatch == checkTracePoints(area) {
				t.Error("TracePoints not consist")
			}
		})
	}
}

func checkTracePoints(area SlicedArea) bool {
	rows := len(area) - 1
	cols := len(area[0]) - 1
	return area[0][0].Color == tracePointColor &&
		area[0][cols].Color == tracePointColor &&
		area[rows][0].Color == tracePointColor &&
		area[rows][cols].Color == tracePointColor
}

func comparePixelDto(value PixelDto, exceptedValue PixelDto) bool {
	return value.Y == exceptedValue.Y &&
		value.X == exceptedValue.X &&
		strings.Compare(value.Color, exceptedValue.Color) == 0 &&
		value.UserId == exceptedValue.UserId
}
