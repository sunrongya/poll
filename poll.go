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

func NewPoll() es.Aggregate {
	return &Poll{
		_userChoices: make(map[Voter][]Choice),
	}
}

func (p *Poll) ProcessCreatePollCommand(command *CreatePollCommand) []es.Event {
	return []es.Event{&PollCreatedEvent{Title: command.Title, Choices: command.Choices}}
}

func (p *Poll) ProcessVotePollBecauseOfVoteRecordCommand(command *VotePollBecauseOfVoteRecordCommand) []es.Event {
	if p._title == "" || len(p._choices) == 0 {
		panic(fmt.Errorf("poll aggegate error"))
	}
	if _, ok := p._userChoices[command.Voter]; ok || !p.isContains(command.Choices) {
		return []es.Event{
			&VotePollFailedBecauseOfVoteRecordEvent{VoteDetails: command.VoteDetails},
		}
	}
	return []es.Event{
		&VotePollCompletedBecauseOfVoteRecordEvent{VoteDetails: command.VoteDetails},
	}
}

func (p *Poll) HandlePollCreatedEvent(event *PollCreatedEvent) {
	p._title, p._choices = event.Title, event.Choices
}

func (p *Poll) HandleVotePollCompletedBecauseOfVoteRecordEvent(event *VotePollCompletedBecauseOfVoteRecordEvent) {
	p._userChoices[event.Voter] = event.Choices
}

func (p *Poll) HandleVotePollFailedBecauseOfVoteRecordEvent(event *VotePollFailedBecauseOfVoteRecordEvent) {
}

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
