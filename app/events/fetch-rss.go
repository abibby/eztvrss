package events

import (
	"github.com/abibby/salusa/event"
	"github.com/abibby/salusa/event/cron"
)

type FetchRSSEvent struct {
	event.EventLogger
	cron.CronEvent
}

var _ event.Event = (*FetchRSSEvent)(nil)

func (e *FetchRSSEvent) Type() event.EventType {
	return "eztvrss:fetch-rss"
}
