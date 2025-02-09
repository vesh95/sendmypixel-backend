package canvas

import "regexp"

const (
	size = 10000
)

type Column struct {
	Color  string
	UserId int64
}

func colorStringValidate(clr string) error {
	ok, _ := regexp.MatchString("#[0-9A-Fa-f]{6}", clr)

	if !ok {
		return &ColorInvalid{color: clr}
	}

	return nil
}

func validateSize(x, y int32) error {
	if x >= size || y >= size || x < 0 || y < 0 {
		return &OutOfBounds{
			x: x,
			y: y,
		}
	}

	return nil
}

func validateDto(dto PixelDto) error {
	err := validateSize(dto.X, dto.Y)

	if err != nil {
		return err
	}

	err = colorStringValidate(dto.Color)
	if err != nil {
		return err
	}

	return nil
}
