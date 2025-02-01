package canvas

import (
	"errors"
	"strings"
	"testing"
)

func TestSyncCanvas_SetPixel(t *testing.T) {
	fixture := NewSyncCanvas()

	cases := []struct {
		name          string
		dto           PixelDto
		err           error
		exceptedValue bool
	}{
		{
			name: "Successfully set",
			dto: PixelDto{
				Y:      0,
				X:      0,
				Color:  "#FAFAFA",
				UserId: 0,
			},
			err:           nil,
			exceptedValue: true,
		},
		{
			name: "Set out of bound pixel by Y",
			dto: PixelDto{
				Y:      size,
				X:      0,
				Color:  "#F0F0F0",
				UserId: 0,
			},
			err:           OutOfBounds{x: 0, y: size},
			exceptedValue: false,
		},
		{
			name: "Set out of bound pixel by X",
			dto: PixelDto{
				Y:      0,
				X:      size,
				Color:  "#F0F0F0",
				UserId: 0,
			},
			err:           OutOfBounds{x: size, y: 0},
			exceptedValue: false,
		},
		{
			name: "Set out of bound pixel by Y with under 0 value",
			dto: PixelDto{
				Y:      -1,
				X:      0,
				Color:  "#F0F0F0",
				UserId: 0,
			},
			err:           OutOfBounds{x: 0, y: -1},
			exceptedValue: false,
		},
		{
			name: "Set out of bound pixel by X with under 0 value",
			dto: PixelDto{
				Y:      0,
				X:      -1,
				Color:  "#F0F0F0",
				UserId: 0,
			},
			err:           OutOfBounds{x: -1, y: 0},
			exceptedValue: false,
		},
		{
			name: "Set invalid color",
			dto: PixelDto{
				Y:      0,
				X:      0,
				Color:  "yellow",
				UserId: 0,
			},
			err:           ColorInvalid{color: "yellow"},
			exceptedValue: false,
		},
		{
			name: "Successfully reset",
			dto: PixelDto{
				Y:      0,
				X:      0,
				Color:  "#FAFAF3",
				UserId: 1,
			},
			err:           nil,
			exceptedValue: true,
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

			if c.exceptedValue != ok {
				t.Errorf("%v was expected, %v was received.", c.exceptedValue, ok)
			}

			return
		})
	}
}

func comparePixelDto(value PixelDto, exceptedValue PixelDto) bool {
	return value.Y == exceptedValue.Y &&
		value.X == exceptedValue.X &&
		strings.Compare(value.Color, exceptedValue.Color) == 0 &&
		value.UserId == exceptedValue.UserId
}
