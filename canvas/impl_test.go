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
				Color:  "#FAFAFC",
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

			return
		})
	}
}

func TestSyncCanvas_GetPixel(t *testing.T) {
	fixture := NewSyncCanvas()
	fixture.canvas[0][0] = Column{
		Color:  "#FAFAFC",
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

func comparePixelDto(value PixelDto, exceptedValue PixelDto) bool {
	return value.Y == exceptedValue.Y &&
		value.X == exceptedValue.X &&
		strings.Compare(value.Color, exceptedValue.Color) == 0 &&
		value.UserId == exceptedValue.UserId
}
