package canvas

import (
	"errors"
	"fmt"
)

var ErrorOutOfBounds = errors.New("out of bounds error")

type OutOfBounds struct {
	x, y int32
}

func (o *OutOfBounds) Error() string {
	return fmt.Sprintf("Values out of bouns X: %d Y: %d", o.x, o.y)
}
func NewOutOfBoundsError(x, y int32) error {

	return &OutOfBounds{x, y}
}
func (o *OutOfBounds) Unwrap() error {
	return ErrorOutOfBounds
}

var ErrorColorInvalid = errors.New("color invalid")

type ColorInvalid struct {
	color string
}

func NewColorInvalidError(color string) error {
	return &ColorInvalid{color}
}
func (c *ColorInvalid) Error() string {
	return fmt.Sprintf("Color value %s is invalid", c.color)
}
func (c *ColorInvalid) Unwrap() error {
	return ErrorColorInvalid
}
