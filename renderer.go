package sangocui

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

func (s *Sangocui) CreateViews() []*gocui.View {
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

func (s *Sangocui) CreateView(p *Panel) *gocui.View {

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
	v.Wrap = true
	v.Overwrite = p.Overwrite

	return v
}

func (s *Sangocui) UpdateListView(viewName string, data []string) {
	go s.g.Update(func(g *gocui.Gui) error {
		v, err := g.View(viewName)
		if err != nil {
			// handle error
		}
		v.Clear()
		for _, item := range data {
			fmt.Fprintln(v, item)
		}

		return nil
	})
}

func (s *Sangocui) UpdateTextView(viewName string, data string) error {
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

func (s *Sangocui) CursorDown(viewName string) error {
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

func (s *Sangocui) CursorUp(viewName string) error {
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

func (s *Sangocui) ResetCursor(viewName string) error {
	v, err := s.g.View(viewName)
	if err != nil {
		return err
	}
	if err := v.SetCursor(0, 0); err != nil {
		return err
	}

	return nil
}

func (s *Sangocui) getPanelByViewName(viewName string) *Panel {
	for i, p := range s.panels {
		if p.Name == viewName {
			return s.panels[i]
		}
	}

	return nil
}

func (s *Sangocui) EnableSelection(viewName string) error {
	pan := s.getPanelByViewName(viewName)
	if pan == nil {
		return nil
	}
	pan.EnableSelection()
	return nil
}

func (s *Sangocui) DisableSelection(viewName string) error {
	pan := s.getPanelByViewName(viewName)
	if pan == nil {
		return nil
	}
	pan.DisableSelection()
	return nil
}

func (s *Sangocui) Quit() error {
	return gocui.ErrQuit
}

func (s *Sangocui) GetCurrentLine(v *gocui.View) string {
	_, cursorY := v.Cursor()
	l, _ := v.Line(cursorY)

	return l
}

func (s *Sangocui) GetNextLine(v *gocui.View) string {
	_, cursorY := v.Cursor()
	l, _ := v.Line(cursorY + 1)

	return l
}

func (s *Sangocui) Focus(viewName string) error {
	if _, err := s.g.SetCurrentView(viewName); err != nil {
		log.Panicln("Try to focus on " + viewName)
		return err
	}

	return nil
}

func (s *Sangocui) Show(viewName string) error {
	pan := s.getPanelByViewName(viewName)
	if pan == nil {
		return nil
	}
	pan.Hidden = false
	s.CreateViews()

	return nil
}

func (s *Sangocui) Hide(viewName string) error {
	pan := s.getPanelByViewName(viewName)
	if pan == nil {
		return nil
	}
	pan.Hidden = true
	s.g.DeleteView(viewName)

	return nil
}

func (s *Sangocui) GetCurrentBuffer(viewName string) string {
	v, err := s.g.View(viewName)
	if err != nil {
		log.Panicln("View not found")
	}
	return v.Buffer()
}
