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

func (this *Poll) ProcessCreatePollCommand(command *CreatePollCommand) []es.Event {
	return []es.Event{&PollCreatedEvent{Title: command.Title, Choices: command.Choices}}
}

func (this *Poll) ProcessVotePollBecauseOfVoteRecordCommand(command *VotePollBecauseOfVoteRecordCommand) []es.Event {
	if this._title == "" || len(this._choices) == 0 {
		panic(fmt.Errorf("poll aggegate error"))
	}
	if _, ok := this._userChoices[command.Voter]; ok || !this.isContains(command.Choices) {
		return []es.Event{
			&VotePollFailedBecauseOfVoteRecordEvent{VoteDetails: command.VoteDetails},
		}
	}
	return []es.Event{
		&VotePollCompletedBecauseOfVoteRecordEvent{VoteDetails: command.VoteDetails},
	}
}

func (this *Poll) HandlePollCreatedEvent(event *PollCreatedEvent) {
	this._title, this._choices = event.Title, event.Choices
}

func (this *Poll) HandleVotePollCompletedBecauseOfVoteRecordEvent(event *VotePollCompletedBecauseOfVoteRecordEvent) {
	this._userChoices[event.Voter] = event.Choices
}

func (this *Poll) HandleVotePollFailedBecauseOfVoteRecordEvent(event *VotePollFailedBecauseOfVoteRecordEvent) {
}

func (this *Poll) isContains(choices []Choice) bool {
	if len(choices) == 0 {
		return false
	}
	for _, choice := range choices {
		contain := false
		for _, v := range this._choices {
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
