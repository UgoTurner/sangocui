package songocui

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

func (s *Songocui) CreateViews() []*gocui.View {
	var views []*gocui.View

	for _, pan := range s.panels {
		if pan.Hidden {
			continue
		}

		lv := s.CreateView(pan)
		if lv == nil {
			continue
		}
		views = append(views, lv)

	}

	return views
}

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

func getWhiteSpaces(width int) string {
	var ws string
	for i := 0; i < width; i++ {
		ws += " "
	}

	return ws
}

func (s *Songocui) UpdateListView(viewName string, data []string) {
	go s.g.Update(func(g *gocui.Gui) error {
		v, err := g.View(viewName)
		if err != nil {
			// handle error
		}
		v.Clear()
		width, _ := v.Size()
		for _, item := range data {
			fmt.Fprintln(v, item+getWhiteSpaces(width))
		}

		return nil
	})
}

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

func (s *Songocui) CursorDown(viewName string) error {
	v, err := s.g.View(viewName)
	if err != nil {
		return err
	}
	if v != nil && s.GetNextLine(v) != "" {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Songocui) CursorUp(viewName string) error {
	v, err := s.g.View(viewName)
	if err != nil {
		return err
	}
	ox, oy := v.Origin()
	cx, cy := v.Cursor()
	if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
		if err := v.SetOrigin(ox, oy-1); err != nil {
			return err
		}
	}
	return nil
}

func (s *Songocui) ResetCursor(viewName string) error {
	v, err := s.g.View(viewName)
	if err != nil {
		return err
	}
	if err := v.SetCursor(0, 0); err != nil {
		return err
	}

	return nil
}

func (s *Songocui) getPanelByViewName(viewName string) *Panel {
	for i, p := range s.panels {
		if p.Name == viewName {
			return s.panels[i]
		}
	}

	return nil
}

func (s *Songocui) EnableSelection(viewName string) error {
	pan := s.getPanelByViewName(viewName)
	if pan == nil {
		return nil
	}
	pan.EnableSelection()
	return nil
}

func (s *Songocui) DisableSelection(viewName string) error {
	pan := s.getPanelByViewName(viewName)
	if pan == nil {
		return nil
	}
	pan.DisableSelection()
	return nil
}

func (s *Songocui) Quit() error {
	return gocui.ErrQuit
}

func (s *Songocui) GetCurrentLine(v *gocui.View) string {
	_, cursorY := v.Cursor()
	l, _ := v.Line(cursorY)

	return l
}

func (s *Songocui) GetNextLine(v *gocui.View) string {
	_, cursorY := v.Cursor()
	l, _ := v.Line(cursorY + 1)

	return l
}

func (s *Songocui) Focus(viewName string) error {
	if _, err := s.g.SetCurrentView(viewName); err != nil {
		log.Panicln("Try to focus on " + viewName)
		return err
	}

	return nil
}

func (s *Songocui) Show(viewName string) error {
	pan := s.getPanelByViewName(viewName)
	if pan == nil {
		return nil
	}
	pan.Hidden = false
	s.CreateViews()

	return nil
}

func (s *Songocui) Hide(viewName string) error {
	pan := s.getPanelByViewName(viewName)
	if pan == nil {
		return nil
	}
	pan.Hidden = true
	s.g.DeleteView(viewName)

	return nil
}

func (s *Songocui) GetCurrentBuffer(viewName string) string {
	v, err := s.g.View(viewName)
	if err != nil {
		log.Panicln("View not found")
	}
	return v.Buffer()
}
