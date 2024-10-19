package songocui

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/sirupsen/logrus"
	"github.com/jroimartin/gocui"
)

// Subscriber defines the interface for Songocui event subscribers.
type Subscriber interface {
	On(string) error
}

// Songocui is a wrapper for gocui.
type Songocui struct {
	g           *gocui.Gui
	panels      []*Panel
	subscribers []Subscriber
	logger      *logrus.Logger
}

// NewWithLogger creates a new Songocui instance with a logger.
func NewWithLogger(logger *logrus.Logger) *Songocui {
	return &Songocui{logger: logger}
}

// Configure sets up keybindings and panels based on configuration files and sets the default view focus.
func (s *Songocui) Configure(pathConfViews, pathConfKeybinds, defaultFocus string) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln("Error creating GUI:", err)
	}
	s.g = g

	s.g.SetManagerFunc(func(gg *gocui.Gui) error {
		s.CreateViews()
		if gg.CurrentView() == nil {
			if _, err := gg.SetCurrentView(defaultFocus); err != nil {
				log.Panicln("Error setting default focus view:", err)
			}
		}
		return nil
	})

	maxX, maxY := s.g.Size()
	s.panels = s.loadPanels(pathConfViews, maxX, maxY)
	s.CreateKeybinds(s.loadKeybinds(pathConfKeybinds))
}

// Boot starts the main GUI loop and dispatches the "Launch" event.
func (s *Songocui) Boot() {
	s.dispatch("Launch")
	defer s.g.Close()

	if err := s.g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln("Main loop error:", err)
	}
}

// loadPanels loads the panel configurations from a JSON file and scales their coordinates.
func (s *Songocui) loadPanels(path string, maxX, maxY int) []*Panel {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicln("Error reading panel config file:", err)
	}

	var panels []*Panel
	if err := json.Unmarshal(data, &panels); err != nil {
		log.Panicln("Error unmarshalling panel config:", err)
	}

	for _, panel := range panels {
		panel.Coordinate.Scale(maxX, maxY)
	}

	return panels
}

// RegisterSubscribers registers a list of subscribers for event dispatching.
func (s *Songocui) RegisterSubscribers(subscribers []Subscriber) {
	s.subscribers = subscribers
}

// dispatch sends an event to all registered subscribers.
func (s *Songocui) dispatch(eventName string) error {
	for _, sub := range s.subscribers {
		if err := sub.On(eventName); err != nil {
			return err
		}
	}
	return nil
}

// loadKeybinds loads keybinding configurations from a JSON file.
func (s *Songocui) loadKeybinds(path string) []confViewsKeybind {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicln("Error reading keybind config file:", err)
	}

	var viewsKeybinds []confViewsKeybind
	if err := json.Unmarshal(data, &viewsKeybinds); err != nil {
		log.Panicln("Error unmarshalling keybind config:", err)
	}

	return viewsKeybinds
}

// CreateKeybinds sets up keybindings for the views based on the loaded configurations.
func (s *Songocui) CreateKeybinds(viewsKeybinds []confViewsKeybind) {
	for _, ckb := range viewsKeybinds {
		for _, kb := range ckb.Keybinds {
			viewName, key, action := ckb.ViewName, kb.Key, kb.Action
			s.logger.Infof("Registering keybind for %s: %s -> %s", viewName, key, action)
			if err := s.g.SetKeybinding(viewName, KeyStrToCode(key), gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
				s.logger.Infof("Dispatching action '%s' in view '%s'", action, v.Name())
				return s.dispatch(action)
			}); err != nil {
				log.Panicln("Error setting keybinding:", err)
			}
		}
	}
}

// confViewsKeybind holds the view name and associated keybindings.
type confViewsKeybind struct {
	ViewName string
	Keybinds []confKeybind
}

// confKeybind holds the key-action pair for keybinding configuration.
type confKeybind struct {
	Key    string
	Action string
}
