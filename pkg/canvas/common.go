package canvas

import (
	"regexp"
)

const (
	Size = 1000
)

type Column struct {
	Color  string
	UserId int64
}

func ColorStringValidate(clr string) error {
	ok, _ := regexp.MatchString("#[0-9A-Fa-f]{6}", clr)

	if !ok {
		return ErrorColorInvalid
	}

	return nil
}

func ValidateSize(x, y int32) error {
	if x >= Size || y >= Size || x < 0 || y < 0 {
		return ErrorOutOfBounds
	}

	return nil
}

func ValidateDto(dto PixelDto) error {
	err := ValidateSize(dto.X, dto.Y)

	if err != nil {
		return err
	}

	err = ColorStringValidate(dto.Color)
	if err != nil {
		return err
	}

	return nil
}
