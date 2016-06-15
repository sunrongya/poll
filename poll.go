package main

import (
	"fmt"
	es "github.com/sunrongya/eventsourcing"
)

type Poll struct {
	es.BaseAggregate
	_title       string
	_choices     []Choice
	_userChoices map[Voter][]Choice
}

var _ es.Aggregate = (*Poll)(nil)

func (p *Poll) ApplyEvents(events []es.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *PollCreatedEvent:
			p._title, p._choices = e.Title, e.Choices
		case *VotePollCompletedBecauseOfVoteRecordEvent:
			p._userChoices[e.Voter] = e.Choices
		case *VotePollFailedBecauseOfVoteRecordEvent:
		default:
			panic(fmt.Errorf("Unknown event %#v", e))
		}
	}
	p.SetVersion(len(events))
}

func (p *Poll) ProcessCommand(command es.Command) []es.Event {
	var event es.Event
	switch c := command.(type) {
	case *CreatePollCommand:
		event = p.processCreatePollCommand(c)
	case *VotePollBecauseOfVoteRecordCommand:
		event = p.processVotePollBecauseOfVoteRecordCommand(c)
	default:
		panic(fmt.Errorf("Unknown command %#v", c))
	}
	event.SetGuid(command.GetGuid())
	return []es.Event{event}
}

func (p *Poll) processCreatePollCommand(command *CreatePollCommand) es.Event {
	return &PollCreatedEvent{Title: command.Title, Choices: command.Choices}
}

func (p *Poll) processVotePollBecauseOfVoteRecordCommand(command *VotePollBecauseOfVoteRecordCommand) es.Event {
	if p._title == "" || len(p._choices) == 0 {
		panic(fmt.Errorf("poll aggegate error"))
	}
	if _, ok := p._userChoices[command.Voter]; ok || !p.isContains(command.Choices) {
		return &VotePollFailedBecauseOfVoteRecordEvent{VoteDetails: command.VoteDetails}

	}
	return &VotePollCompletedBecauseOfVoteRecordEvent{VoteDetails: command.VoteDetails}
}

// TODO 后面再来优化
func (p *Poll) isContains(choices []Choice) bool {
	if len(choices) == 0 {
		return false
	}
	for _, choice := range choices {
		contain := false
		for _, v := range p._choices {
			if v == choice {
				contain = true
				break
			}
		}
		if !contain {
			return false
		}
	}
	return true
}

func NewPoll() es.Aggregate {
	return &Poll{
		_userChoices: make(map[Voter][]Choice),
	}
}
