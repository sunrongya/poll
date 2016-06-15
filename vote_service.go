package main

import (
	es "github.com/sunrongya/eventsourcing"
	"time"
)

type VoteService struct {
	es.Service
}

func NewVoteService(store es.EventStore) *VoteService {
	service := &VoteService{
		Service: es.NewService(store, NewVoteRecord),
	}
	return service
}

func (p *VoteService) VotePoll(pollId es.Guid, voter Voter, choices []Choice, created_time time.Time) es.Guid {
	guid := es.NewGuid()
	c := &CreateVoteRecordCommand{
		WithGuid: es.WithGuid{guid},
		VoteDetails: VoteDetails{
			VoteRecordId: guid,
			PollId:       pollId,
			Voter:        voter,
			Choices:      choices,
			Time:         created_time,
		},
	}
	p.PublishCommand(c)
	return guid
}
