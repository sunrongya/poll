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
	repository es.ReadRepository
}

func NewPollProjector(repository es.ReadRepository) *PollProjector {
	return &PollProjector{repository: repository}
}

func (g *PollProjector) HandlePollCreatedEvent(event *PollCreatedEvent) {
	poll := &RPoll{
		Id:         event.GetGuid(),
		Title:      event.Title,
		Choices:    event.Choices,
		ChoiceStat: make(map[Choice]int),
	}
	g.repository.Save(poll.Id, poll)
}

func (g *PollProjector) HandleVotePollCompletedBecauseOfVoteRecordEvent(event *VotePollCompletedBecauseOfVoteRecordEvent) {
	g.do(event.GetGuid(), func(poll *RPoll) {
		for _, choice := range event.Choices {
			poll.ChoiceStat[choice] += 1
		}
	})
}

func (g *PollProjector) do(id es.Guid, assignRPollFn func(*RPoll)) {
	i, err := g.repository.Find(id)
	if err != nil {
		return
	}
	poll := i.(*RPoll)
	assignRPollFn(poll)
	g.repository.Save(id, poll)
}
