package canvas

import "fmt"

type OutOfBounds struct {
	x, y int32
}
type ColorInvalid struct {
	color string
}

func (o OutOfBounds) Error() string {
	return fmt.Sprintf("Values out of bouns x: %d y: %d", o.x, o.y)
}
func (c ColorInvalid) Error() string {
	return fmt.Sprintf("Color value %s is invalid", c.color)
}
