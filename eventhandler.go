package main

import (
	es "github.com/sunrongya/eventsourcing"
)

type EventHandler struct {
	_pollChan chan<- es.Command
	_voteChan chan<- es.Command
}

func (this *EventHandler) HandleVoteRecordCreatedEvent(event *VoteRecordCreatedEvent) {
	this._pollChan <- &VotePollBecauseOfVoteRecordCommand{
		WithGuid:    es.WithGuid{Guid: event.PollId},
		VoteDetails: event.VoteDetails,
	}
}

func (this *EventHandler) HandleVotePollCompletedBecauseOfVoteRecordEvent(event *VotePollCompletedBecauseOfVoteRecordEvent) {
	this._voteChan <- &CompleteVoteRecordCommand{
		WithGuid:    es.WithGuid{Guid: event.VoteRecordId},
		VoteDetails: event.VoteDetails,
	}
}

func (this *EventHandler) HandleVotePollFailedBecauseOfVoteRecordEvent(event *VotePollFailedBecauseOfVoteRecordEvent) {
	this._voteChan <- &FailVoteRecordCommand{
		WithGuid:    es.WithGuid{Guid: event.VoteRecordId},
		VoteDetails: event.VoteDetails,
	}
}

func NewEventHandler(pollChan, voteChan chan<- es.Command) *EventHandler {
	return &EventHandler{
		_pollChan: pollChan,
		_voteChan: voteChan,
	}
}
