package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	es "github.com/sunrongya/eventsourcing"
	"testing"
	"time"
)

func TestVoteRecordRestore(t *testing.T) {
	voteDetails := VoteDetails{
		VoteRecordId: es.NewGuid(),
		PollId:       es.NewGuid(),
		Voter:        "voter1",
		Choices:      []Choice{"Java", "Golang"},
		Time:         time.Now(),
	}
	voteRecord := &VoteRecord{}
	voteRecord.ApplyEvents([]es.Event{
		&VoteRecordCreatedEvent{VoteDetails: voteDetails},
		&VoteRecordCompletedEvent{VoteDetails: voteDetails},
	})
	assert.Equal(t, 2, voteRecord.Version(), "version error")
	assert.Equal(t, Completed, voteRecord._state, "state error")
	assert.Equal(t, voteDetails, voteRecord._voteDetails, "vote details error")
}

func TestVoteRecordRestoreForErrorEvent(t *testing.T) {
	assert.Panics(t, func() {
		NewVoteRecord().ApplyEvents([]es.Event{&struct{ es.WithGuid }{}})
	}, "restore error event must panic error")
}

func TestCheckVoteRecordApplyEvents(t *testing.T) {
	events := []es.Event{
		&VoteRecordCreatedEvent{},
		&VoteRecordCompletedEvent{},
		&VoteRecordFailedEvent{},
	}
	assert.NotPanics(t, func() { NewVoteRecord().ApplyEvents(events) }, "Check Process All Event")
}

func TestVoteRecordCommand(t *testing.T) {
	voteDetails := VoteDetails{
		VoteRecordId: es.NewGuid(),
		PollId:       es.NewGuid(),
		Voter:        "adj",
		Choices:      []Choice{"PHP", "Java", "Golang"},
		Time:         time.Now(),
	}

	tests := []struct {
		voteRecord *VoteRecord
		command    es.Command
		event      es.Event
	}{
		{
			&VoteRecord{},
			&CreateVoteRecordCommand{WithGuid: es.WithGuid{Guid: voteDetails.VoteRecordId}, VoteDetails: voteDetails},
			&VoteRecordCreatedEvent{WithGuid: es.WithGuid{Guid: voteDetails.VoteRecordId}, VoteDetails: voteDetails},
		},
		{
			&VoteRecord{_state: Created},
			&CompleteVoteRecordCommand{WithGuid: es.WithGuid{Guid: voteDetails.VoteRecordId}, VoteDetails: voteDetails},
			&VoteRecordCompletedEvent{WithGuid: es.WithGuid{Guid: voteDetails.VoteRecordId}, VoteDetails: voteDetails},
		},
		{
			&VoteRecord{_state: Created},
			&FailVoteRecordCommand{WithGuid: es.WithGuid{Guid: voteDetails.VoteRecordId}, VoteDetails: voteDetails},
			&VoteRecordFailedEvent{WithGuid: es.WithGuid{Guid: voteDetails.VoteRecordId}, VoteDetails: voteDetails},
		},
	}

	for _, v := range tests {
		assert.Equal(t, []es.Event{v.event}, v.voteRecord.ProcessCommand(v.command))
	}
}

func TestVoteRecordCommand_Panic(t *testing.T) {
	tests := []struct {
		voteRecord *VoteRecord
		command    es.Command
	}{
		{
			&VoteRecord{},
			&struct{ es.WithGuid }{},
		},
		{
			&VoteRecord{},
			&CompleteVoteRecordCommand{},
		},
		{
			&VoteRecord{_state: Completed},
			&CompleteVoteRecordCommand{},
		},
		{
			&VoteRecord{_state: Failured},
			&CompleteVoteRecordCommand{},
		},
		{
			&VoteRecord{},
			&FailVoteRecordCommand{},
		},
		{
			&VoteRecord{_state: Completed},
			&FailVoteRecordCommand{},
		},
		{
			&VoteRecord{_state: Failured},
			&FailVoteRecordCommand{},
		},
	}

	for _, v := range tests {
		assert.Panics(t, func() { v.voteRecord.ProcessCommand(v.command) }, fmt.Sprintf("test panics error: command:%v", v.command))
	}
}
