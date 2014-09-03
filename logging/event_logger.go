package logging

import (
	"fmt"
	"time"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	durationfmt "github.com/cloudfoundry/bosh-micro-cli/durationfmt"

	bmui "github.com/cloudfoundry/bosh-micro-cli/ui"
)

type Event struct {
	Time    time.Time
	Stage   string
	Total   int
	Task    string
	State   EventState
	Index   int
	Message string
}

type EventState string

const (
	Started  EventState = "started"
	Finished EventState = "finished"
	Failed   EventState = "failed"
)

type EventLogger interface {
	AddEvent(event Event) error
}

func NewEventLogger(ui bmui.UI) EventLogger {
	return &eventLogger{
		ui:           ui,
		startedTasks: make(map[string]time.Time),
	}
}

type eventLogger struct {
	ui           bmui.UI
	startedTasks map[string]time.Time
}

func (e *eventLogger) AddEvent(event Event) error {
	key := fmt.Sprintf("%s > %s.", event.Stage, event.Task)
	switch event.State {
	case Started:
		if event.Index == 1 {
			e.ui.Sayln(fmt.Sprintf("Started %s", event.Stage))
		}
		e.ui.Say(fmt.Sprintf("Started %s", key))
		e.startedTasks[key] = event.Time
	case Finished:
		duration := event.Time.Sub(e.startedTasks[key])
		e.ui.Sayln(fmt.Sprintf(" Done (%s)", durationfmt.Format(duration)))
		if event.Index == event.Total {
			e.ui.Sayln(fmt.Sprintf("Done %s", event.Stage))
		}
	case Failed:
		duration := event.Time.Sub(e.startedTasks[key])
		e.ui.Sayln(fmt.Sprintf(" Failed '%s' (%s)", event.Message, durationfmt.Format(duration)))
	default:
		return bosherr.New("Unsupported event state `%s'", event.State)
	}
	return nil
}