package main

import (
	es "github.com/sunrongya/eventsourcing"
)

type RPoll struct {
	Id         es.Guid
	Title      string
	Choices    []Choice
	ChoiceStat map[Choice]int
}

type PollProjector struct {
	_repository es.ReadRepository
}

func NewPollProjector(repository es.ReadRepository) *PollProjector {
	return &PollProjector{_repository: repository}
}

func (this *PollProjector) HandlePollCreatedEvent(event *PollCreatedEvent) {
	poll := &RPoll{
		Id:         event.GetGuid(),
		Title:      event.Title,
		Choices:    event.Choices,
		ChoiceStat: make(map[Choice]int),
	}
	this._repository.Save(poll.Id, poll)
}

func (this *PollProjector) HandleVotePollCompletedBecauseOfVoteRecordEvent(event *VotePollCompletedBecauseOfVoteRecordEvent) {
	this.do(event.GetGuid(), func(poll *RPoll) {
		for _, choice := range event.Choices {
			poll.ChoiceStat[choice] += 1
		}
	})
}

func (this *PollProjector) do(id es.Guid, assignRPollFn func(*RPoll)) {
	i, err := this._repository.Find(id)
	if err != nil {
		return
	}
	poll := i.(*RPoll)
	assignRPollFn(poll)
	this._repository.Save(id, poll)
}
