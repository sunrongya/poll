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

func (p *VoteRecord) ApplyEvents(events []es.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *VoteRecordCreatedEvent:
			p._state, p._voteDetails = Created, e.VoteDetails
		case *VoteRecordCompletedEvent:
			p._state = Completed
		case *VoteRecordFailedEvent:
			p._state = Failured
		default:
			panic(fmt.Errorf("Unknown event %#v", e))
		}
	}
	p.SetVersion(len(events))
}

func (p *VoteRecord) ProcessCommand(command es.Command) []es.Event {
	var event es.Event
	switch c := command.(type) {
	case *CreateVoteRecordCommand:
		event = p.processCreateVoteRecordCommand(c)
	case *CompleteVoteRecordCommand:
		event = p.processCompleteVoteRecordCommand(c)
	case *FailVoteRecordCommand:
		event = p.processFailVoteRecordCommand(c)
	default:
		panic(fmt.Errorf("Unknown command %#v", c))
	}
	event.SetGuid(command.GetGuid())
	return []es.Event{event}
}

func (p *VoteRecord) processCreateVoteRecordCommand(command *CreateVoteRecordCommand) es.Event {
	return &VoteRecordCreatedEvent{VoteDetails: command.VoteDetails}
}

func (p *VoteRecord) processCompleteVoteRecordCommand(command *CompleteVoteRecordCommand) es.Event {
	if p._state != Created {
		panic(fmt.Errorf("Can't process CompleteVoteRecordCommand of state:%s", p._state))
	}
	return &VoteRecordCompletedEvent{VoteDetails: command.VoteDetails}
}

func (p *VoteRecord) processFailVoteRecordCommand(command *FailVoteRecordCommand) es.Event {
	if p._state != Created {
		panic(fmt.Errorf("Can't process FailVoteRecordCommand of state:%s", p._state))
	}
	return &VoteRecordFailedEvent{VoteDetails: command.VoteDetails}
}

func NewVoteRecord() es.Aggregate {
	return &VoteRecord{}
}
