package main

import (
	es "github.com/sunrongya/eventsourcing"
)

type PollService struct {
	es.Service
}

func NewPollService(store es.EventStore) *PollService {
	service := &PollService{
		Service: es.NewService(store, NewPoll),
	}
	return service
}

func (p *PollService) CreatePoll(title string, choices []Choice) es.Guid {
	guid := es.NewGuid()
	c := &CreatePollCommand{
		WithGuid: es.WithGuid{guid},
		Title:    title,
		Choices:  choices,
	}
	p.PublishCommand(c)
	return guid
}
