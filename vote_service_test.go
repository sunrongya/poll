package main

import (
	es "github.com/sunrongya/eventsourcing"
	"github.com/sunrongya/eventsourcing/utiltest"
	"testing"
	"time"
)

func TestVoteServiceDoCreatePoll(t *testing.T) {
	utiltest.TestServicePublishCommand(t, func(service es.Service) es.Command {
		gs := VoteService{Service: service}
		pollId := es.NewGuid()
		voter := Voter("sry")
		choices := []Choice{"Golang", "Haskell", "Node.js"}
		now := time.Now()

		guid := gs.VotePoll(pollId, voter, choices, now)
		return &CreateVoteRecordCommand{
			WithGuid: es.WithGuid{guid},
			VoteDetails: VoteDetails{
				VoteRecordId: guid,
				PollId:       pollId,
				Voter:        voter,
				Choices:      choices,
				Time:         now,
			},
		}
	})
}
