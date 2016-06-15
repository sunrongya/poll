package main

import (
	es "github.com/sunrongya/eventsourcing"
	"github.com/sunrongya/eventsourcing/utiltest"
	"testing"
)

func TestPollServiceDoCreatePoll(t *testing.T) {
	utiltest.TestServicePublishCommand(t, func(service es.Service) es.Command {
		gs := PollService{Service: service}
		title := "title1"
		choices := []Choice{"PHP", "Java", "Golang", "Haskell", "Node.js"}
		guid := gs.CreatePoll(title, choices)
		return &CreatePollCommand{WithGuid: es.WithGuid{guid}, Title: title, Choices: choices}
	})
}
