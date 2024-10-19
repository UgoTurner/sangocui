package songocui

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

// CreateViews generates and returns all visible views from the panels.
func (s *Songocui) CreateViews() []*gocui.View {
	var views []*gocui.View
	for _, pan := range s.panels {
		if !pan.Hidden {
			if v := s.CreateView(pan); v != nil {
				views = append(views, v)
			}
		}
	}
	return views
}

// CreateView creates a gocui.View based on panel specifications.
func (s *Songocui) CreateView(p *Panel) *gocui.View {
	v, err := s.g.SetView(
		p.Name,
		p.Coordinate.TopLeftXabs,
		p.Coordinate.TopLeftYabs,
		p.Coordinate.BottomRightXabs,
		p.Coordinate.BottomRightYabs,
	)
	if err != nil && err != gocui.ErrUnknownView {
		return nil
	}

	v.SelBgColor = ColorStrToCode(p.SelectionColor.BgColorCurrent)
	v.SelFgColor = ColorStrToCode(p.SelectionColor.FgColorCurrent)
	v.Highlight = p.Highlight
	v.Frame = p.Frame
	v.Title = p.Title
	v.Editable = p.Editable
	v.Wrap = p.Wrap
	v.Overwrite = p.Overwrite

	return v
}

// UpdateListView updates a view with a list of strings and maintains the correct width.
func (s *Songocui) UpdateListView(viewName string, data []string) {
	s.g.Update(func(g *gocui.Gui) error {
		v, err := g.View(viewName)
		if err != nil {
			log.Printf("Error finding view %s: %v", viewName, err)
			return err
		}
		v.Clear()
		width, _ := v.Size()
		for _, item := range data {
			fmt.Fprintln(v, fmt.Sprintf("%s%s", item, getWhiteSpaces(width-len(item))))
		}
		return nil
	})
}

// UpdateTextView updates a text view with a single string.
func (s *Songocui) UpdateTextView(viewName string, data string) error {
	go s.g.Update(func(g *gocui.Gui) error {
		v, err := g.View(viewName)
		if err != nil {
			return err
		}
		v.Clear()
		fmt.Fprintln(v, data)

		return nil
	})
	return nil
}

// CursorDown moves the cursor one line down.
func (s *Songocui) CursorDown(viewName string) error {
	v, err := s.g.View(viewName)
	if err != nil {
		return err
	}
	if s.GetNextLine(v) != "" {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			return v.SetOrigin(ox, oy+1)
		}
	}
	return nil
}

// CursorUp moves the cursor one line up.
func (s *Songocui) CursorUp(viewName string) error {
	v, err := s.g.View(viewName)
	if err != nil {
		return err
	}
	cx, cy := v.Cursor()
	if err := v.SetCursor(cx, cy-1); err != nil {
		ox, oy := v.Origin()
		if oy > 0 {
			return v.SetOrigin(ox, oy-1)
		}
	}
	return nil
}

// ResetCursor resets the cursor position to the top-left of the view.
func (s *Songocui) ResetCursor(viewName string) error {
	v, err := s.g.View(viewName)
	if err != nil {
		return err
	}
	return v.SetCursor(0, 0)
}

// EnableSelection enables selection for a given view.
func (s *Songocui) EnableSelection(viewName string) error {
	if pan := s.getPanelByViewName(viewName); pan != nil {
		pan.EnableSelection()
	}
	return nil
}

// DisableSelection disables selection for a given view.
func (s *Songocui) DisableSelection(viewName string) error {
	if pan := s.getPanelByViewName(viewName); pan != nil {
		pan.DisableSelection()
	}
	return nil
}

// Quit exits the application.
func (s *Songocui) Quit() error {
	return gocui.ErrQuit
}

// GetCurrentLine returns the line where the cursor is located in a view.
func (s *Songocui) GetCurrentLine(v *gocui.View) string {
	_, cy := v.Cursor()
	l, _ := v.Line(cy)
	return l
}

// GetNextLine returns the line below the cursor.
func (s *Songocui) GetNextLine(v *gocui.View) string {
	_, cy := v.Cursor()
	l, _ := v.Line(cy + 1)
	return l
}

// Focus sets the focus to a specific view.
func (s *Songocui) Focus(viewName string) error {
	_, err := s.g.SetCurrentView(viewName)
	if err != nil {
		log.Printf("Error focusing view %s: %v", viewName, err)
	}
	return err
}

// Show unhides a view and recreates the visible views.
func (s *Songocui) Show(viewName string) error {
	if pan := s.getPanelByViewName(viewName); pan != nil {
		pan.Hidden = false
		s.CreateViews()
	}
	return nil
}

// Hide hides a view and removes it from the UI.
func (s *Songocui) Hide(viewName string) error {
	if pan := s.getPanelByViewName(viewName); pan != nil {
		pan.Hidden = true
		s.g.DeleteView(viewName)
	}
	return nil
}

// GetCurrentBuffer returns the current buffer content of a view.
func (s *Songocui) GetCurrentBuffer(viewName string) string {
	v, err := s.g.View(viewName)
	if err != nil {
		log.Printf("Error fetching view %s: %v", viewName, err)
	}
	return v.Buffer()
}

// Helper function to generate a string of whitespaces.
func getWhiteSpaces(width int) string {
	return fmt.Sprintf("%*s", width, "")
}

// getPanelByViewName fetches a panel by its view name.
func (s *Songocui) getPanelByViewName(viewName string) *Panel {
	for _, p := range s.panels {
		if p.Name == viewName {
			return p
		}
	}
	return nil
}
