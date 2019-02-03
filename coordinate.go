package sangocui

// Coordinate : x, y of the top left and x, y of the bottom right of a panel
type Coordinate struct {
	TopLeftXrel, TopLeftYrel, BottomRightXrel, BottomRightYrel int
	TopLeftXabs, TopLeftYabs, BottomRightXabs, BottomRightYabs int
}

// Scale : Converts relative coordinates to absolute
// according max width and max height of the current terminal
func (c *Coordinate) Scale(maxX, maxY int) {
	if c.TopLeftXrel != 0 {
		c.TopLeftXabs = maxX + c.TopLeftXrel
	}
	if c.TopLeftYrel != 0 {
		c.TopLeftYabs = maxY + c.TopLeftYrel
	}
	if c.BottomRightXrel != 0 {
		c.BottomRightXabs = maxX + c.BottomRightXrel
	}
	if c.BottomRightYrel != 0 {
		c.BottomRightYabs = maxY + c.BottomRightYrel
	}
}
