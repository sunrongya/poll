package main

import (
	"fmt"
	es "github.com/sunrongya/eventsourcing"
	"time"
)

type State string

const (
	Created   = State("Created")
	Completed = State("Completed")
	Failured  = State("Failured")
)

type VoteRecord struct {
	es.BaseAggregate
	_voteDetails VoteDetails
	_state       State
}

type VoteDetails struct {
	VoteRecordId es.Guid
	PollId       es.Guid
	Voter        Voter
	Choices      []Choice
	Time         time.Time
}

var _ es.Aggregate = (*VoteRecord)(nil)

func NewVoteRecord() es.Aggregate {
	return &VoteRecord{}
}

func (this *VoteRecord) ProcessCreateVoteRecordCommand(command *CreateVoteRecordCommand) []es.Event {
	return []es.Event{
		&VoteRecordCreatedEvent{VoteDetails: command.VoteDetails},
	}
}

func (this *VoteRecord) ProcessCompleteVoteRecordCommand(command *CompleteVoteRecordCommand) []es.Event {
	if this._state != Created {
		panic(fmt.Errorf("Can't process CompleteVoteRecordCommand of state:%s", this._state))
	}
	return []es.Event{
		&VoteRecordCompletedEvent{VoteDetails: command.VoteDetails},
	}
}

func (this *VoteRecord) ProcessFailVoteRecordCommand(command *FailVoteRecordCommand) []es.Event {
	if this._state != Created {
		panic(fmt.Errorf("Can't process FailVoteRecordCommand of state:%s", this._state))
	}
	return []es.Event{
		&VoteRecordFailedEvent{VoteDetails: command.VoteDetails},
	}
}

func (this *VoteRecord) HandleVoteRecordCreatedEvent(event *VoteRecordCreatedEvent) {
	this._state, this._voteDetails = Created, event.VoteDetails
}

func (this *VoteRecord) HandleVoteRecordCompletedEvent(event *VoteRecordCompletedEvent) {
	this._state = Completed
}

func (this *VoteRecord) HandleVoteRecordFailedEvent(event *VoteRecordFailedEvent) {
	this._state = Failured
}
