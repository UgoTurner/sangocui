package sangocui

type Panel struct {
	Title                                         string
	Name                                          string
	Highlight, Frame, Overwrite, Hidden, Editable bool
	Coordinate                                    Coordinate
	SelectionColor                                SelectionColor
}

// EnableSelection : Set the current selection colors to "active"
func (p *Panel) EnableSelection() {
	p.SelectionColor.BgColorCurrent = p.SelectionColor.BgColorActive
	p.SelectionColor.FgColorCurrent = p.SelectionColor.FgColorActive
}

// DisableSelection : Set the current selection colors to "unactive"
func (p *Panel) DisableSelection() {
	p.SelectionColor.BgColorCurrent = p.SelectionColor.BgColorUnactive
	p.SelectionColor.FgColorCurrent = p.SelectionColor.FgColorUnactive
}
