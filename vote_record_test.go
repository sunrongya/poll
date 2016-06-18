package main

import (
	//"fmt"
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
	voteRecord.HandleVoteRecordCreatedEvent(&VoteRecordCreatedEvent{VoteDetails: voteDetails})
	voteRecord.HandleVoteRecordCompletedEvent(&VoteRecordCompletedEvent{VoteDetails: voteDetails})

	assert.Equal(t, Completed, voteRecord._state, "state error")
	assert.Equal(t, voteDetails, voteRecord._voteDetails, "vote details error")
}

func TestCreateVoteRecordCommand(t *testing.T) {
	voteDetails := VoteDetails{
		VoteRecordId: es.NewGuid(),
		PollId:       es.NewGuid(),
		Voter:        "adj",
		Choices:      []Choice{"PHP", "Java", "Golang"},
		Time:         time.Now(),
	}
	command := &CreateVoteRecordCommand{VoteDetails: voteDetails}
	events := []es.Event{&VoteRecordCreatedEvent{VoteDetails: voteDetails}}

	assert.Equal(t, events, new(VoteRecord).ProcessCreateVoteRecordCommand(command), "创建投票记录命令返回事件有误")
}

func TestCompleteVoteRecordCommand(t *testing.T) {
	voteDetails := VoteDetails{
		VoteRecordId: es.NewGuid(),
		PollId:       es.NewGuid(),
		Voter:        "adj",
		Choices:      []Choice{"PHP", "Java", "Golang"},
		Time:         time.Now(),
	}
	voteRecord := &VoteRecord{_state: Created}
	command := &CompleteVoteRecordCommand{VoteDetails: voteDetails}
	events := []es.Event{&VoteRecordCompletedEvent{VoteDetails: voteDetails}}

	assert.Equal(t, events, voteRecord.ProcessCompleteVoteRecordCommand(command), "成功创建投票记录命令返回事件有误")
}

func TestFailVoteRecordCommand(t *testing.T) {
	voteDetails := VoteDetails{
		VoteRecordId: es.NewGuid(),
		PollId:       es.NewGuid(),
		Voter:        "adj",
		Choices:      []Choice{"PHP", "Java", "Golang"},
		Time:         time.Now(),
	}
	voteRecord := &VoteRecord{_state: Created}
	command := &FailVoteRecordCommand{VoteDetails: voteDetails}
	events := []es.Event{&VoteRecordFailedEvent{VoteDetails: voteDetails}}

	assert.Equal(t, events, voteRecord.ProcessFailVoteRecordCommand(command), "失败创建投票记录命令返回事件有误")
}

func TestCompleteVoteRecordCommand_Panic(t *testing.T) {
	voteRecords := []*VoteRecord{
		&VoteRecord{},
		&VoteRecord{_state: Completed},
		&VoteRecord{_state: Failured},
	}
	for _, voteRecord := range voteRecords {
		assert.Panics(t, func() {
			voteRecord.ProcessCompleteVoteRecordCommand(&CompleteVoteRecordCommand{})
		},
			"执行CompleteVoteRecordCommand应该抛出异常",
		)
	}
}

func TestFailVoteRecordCommand_Panic(t *testing.T) {
	voteRecords := []*VoteRecord{
		&VoteRecord{},
		&VoteRecord{_state: Completed},
		&VoteRecord{_state: Failured},
	}
	for _, voteRecord := range voteRecords {
		assert.Panics(t, func() {
			voteRecord.ProcessFailVoteRecordCommand(&FailVoteRecordCommand{})
		},
			"执行FailVoteRecordCommand应该抛出异常",
		)
	}
}
