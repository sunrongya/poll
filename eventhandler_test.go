package main

import (
	"github.com/stretchr/testify/assert"
	es "github.com/sunrongya/eventsourcing"
	"testing"
	"time"
)

func testHandleEvent(t *testing.T, methodName string, doVotePollHandle func(chan es.Command, VoteDetails) es.Command) {
	voteDetails := VoteDetails{
		VoteRecordId: es.NewGuid(),
		PollId:       es.NewGuid(),
		Voter:        "adj",
		Choices:      []Choice{"PHP", "Java", "Golang"},
		Time:         time.Now(),
	}
	ch := make(chan es.Command)
	command := doVotePollHandle(ch, voteDetails)
	select {
	case c := <-ch:
		assert.Equal(t, c, command, methodName)
	case <-time.After(1 * time.Second):
		t.Error(methodName)
	}
}

func TestHandleVoteRecordCreatedEvent(t *testing.T) {
	testHandleEvent(t, "TestHandleVoteRecordCreatedEvent",
		func(pollCH chan es.Command, details VoteDetails) es.Command {
			handler := NewEventHandler(pollCH, nil)
			go handler.HandleVoteRecordCreatedEvent(
				&VoteRecordCreatedEvent{
					WithGuid:    es.WithGuid{details.VoteRecordId},
					VoteDetails: details,
				},
			)
			return &VotePollBecauseOfVoteRecordCommand{WithGuid: es.WithGuid{details.PollId}, VoteDetails: details}
		},
	)
}

func TestHandleVotePollCompletedBecauseOfVoteRecordEvent(t *testing.T) {
	testHandleEvent(t, "TestHandleVotePollCompletedBecauseOfVoteRecordEvent",
		func(voteCH chan es.Command, details VoteDetails) es.Command {
			handler := NewEventHandler(nil, voteCH)
			go handler.HandleVotePollCompletedBecauseOfVoteRecordEvent(
				&VotePollCompletedBecauseOfVoteRecordEvent{
					WithGuid:    es.WithGuid{details.PollId},
					VoteDetails: details,
				},
			)
			return &CompleteVoteRecordCommand{WithGuid: es.WithGuid{details.VoteRecordId}, VoteDetails: details}
		},
	)
}

func TestHandleVotePollFailedBecauseOfVoteRecordEvent(t *testing.T) {
	testHandleEvent(t, "TestHandleVotePollFailedBecauseOfVoteRecordEvent",
		func(voteCH chan es.Command, details VoteDetails) es.Command {
			handler := NewEventHandler(nil, voteCH)
			go handler.HandleVotePollFailedBecauseOfVoteRecordEvent(
				&VotePollFailedBecauseOfVoteRecordEvent{
					WithGuid:    es.WithGuid{details.PollId},
					VoteDetails: details,
				},
			)
			return &FailVoteRecordCommand{WithGuid: es.WithGuid{details.VoteRecordId}, VoteDetails: details}
		},
	)
}
