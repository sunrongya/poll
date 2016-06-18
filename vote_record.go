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

func (p *VoteRecord) ProcessCreateVoteRecordCommand(command *CreateVoteRecordCommand) []es.Event {
	return []es.Event{
		&VoteRecordCreatedEvent{VoteDetails: command.VoteDetails},
	}
}

func (p *VoteRecord) ProcessCompleteVoteRecordCommand(command *CompleteVoteRecordCommand) []es.Event {
	if p._state != Created {
		panic(fmt.Errorf("Can't process CompleteVoteRecordCommand of state:%s", p._state))
	}
	return []es.Event{
		&VoteRecordCompletedEvent{VoteDetails: command.VoteDetails},
	}
}

func (p *VoteRecord) ProcessFailVoteRecordCommand(command *FailVoteRecordCommand) []es.Event {
	if p._state != Created {
		panic(fmt.Errorf("Can't process FailVoteRecordCommand of state:%s", p._state))
	}
	return []es.Event{
		&VoteRecordFailedEvent{VoteDetails: command.VoteDetails},
	}
}

func (p *VoteRecord) HandleVoteRecordCreatedEvent(event *VoteRecordCreatedEvent) {
	p._state, p._voteDetails = Created, event.VoteDetails
}

func (p *VoteRecord) HandleVoteRecordCompletedEvent(event *VoteRecordCompletedEvent) {
	p._state = Completed
}

func (p *VoteRecord) HandleVoteRecordFailedEvent(event *VoteRecordFailedEvent) {
	p._state = Failured
}
