package songocui

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/sirupsen/logrus"

	"github.com/jroimartin/gocui"
)

// Subscriber : Structs binding Songocui event must implements this interface
type Subscriber interface {
	On(string) error
}

// Songocui : Wrapper for Gocui
type Songocui struct {
	g           *gocui.Gui
	panels      []*Panel
	subscribers []Subscriber
	logger      *logrus.Logger
}

// NewWithLogger : Instanciate a new Songocui object with a logger
func NewWithLogger(logger *logrus.Logger) *Songocui {
	return &Songocui{logger: logger}
}

/* Configure : Create keybinds and panels according conf files
   and view name to focus on at start
*/
func (s *Songocui) Configure(pathConfViews, pathconfKeybinds, defaultFocus string) {
	var err error
	s.g, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln("Error when creating GUI")
	}
	s.g.SetManagerFunc(
		func(gg *gocui.Gui) error {
			s.CreateViews()
			// Focus on side view if no current view
			if gg.CurrentView() == nil {
				if _, err := gg.SetCurrentView(defaultFocus); err != nil {
					log.Panicln(err)
				}
			}
			return nil
		},
	)
	maxX, maxY := s.g.Size()
	s.panels = s.loadPanels(pathConfViews, maxX, maxY)
	s.CreateKeybinds(s.loadKeybinds(pathconfKeybinds))
}

func (s *Songocui) Boot() {
	s.dispatch("Launch")
	defer s.g.Close()
	if err := s.g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panic(err)
	}
}

func (s *Songocui) loadPanels(path string, maxX, maxY int) []*Panel {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicln("Error when opening json file")
	}
	var panels []*Panel
	json.Unmarshal(data, &panels)
	for i := range panels {
		panels[i].Coordinate.Scale(maxX, maxY)
	}

	return panels
}

func (s *Songocui) RegisterSubscribers(subscribers []Subscriber) {
	s.subscribers = subscribers
}

func (s *Songocui) dispatch(eventName string) error {
	for i := range s.subscribers {
		err := s.subscribers[i].On(eventName)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Songocui) loadKeybinds(path string) []confViewsKeybind {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicln(err)
	}
	var viewsKeybinds = []confViewsKeybind{}
	json.Unmarshal(data, &viewsKeybinds)

	return viewsKeybinds
}

type confViewsKeybind struct {
	ViewName string
	Keybinds []confKeybind
}

type confKeybind struct {
	Key    string
	Action string
}

func (s *Songocui) CreateKeybinds(viewsKeybinds []confViewsKeybind) {
	for _, ckb := range viewsKeybinds {
		for _, kb := range ckb.Keybinds {
			viewName := ckb.ViewName
			key := kb.Key
			action := kb.Action
			s.logger.Info("Registering for " + ckb.ViewName + ": " + kb.Key + " --- " + kb.Action)
			s.g.SetKeybinding(
				viewName,
				KeyStrToCode(key),
				gocui.ModNone,
				func(g *gocui.Gui, v *gocui.View) error {
					s.logger.Info("Dispatch " + action + "in view " + v.Name())
					return s.dispatch(action)
				},
			)
		}
	}
}
