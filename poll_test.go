package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	es "github.com/sunrongya/eventsourcing"
	"testing"
	"time"
)

func TestPollRestore(t *testing.T) {
	poll := &Poll{_userChoices: make(map[Voter][]Choice)}
	pollCreatedEvent := &PollCreatedEvent{
		Title:   "喜欢哪几种语言？",
		Choices: []Choice{"PHP", "Java", "Golang", "Haskell", "Node.js"},
	}
	completedEvent := &VotePollCompletedBecauseOfVoteRecordEvent{
		VoteDetails: VoteDetails{
			Voter:   "voter1",
			Choices: []Choice{"Java", "Golang"},
			Time:    time.Now(),
		},
	}
	failedEvent := &VotePollFailedBecauseOfVoteRecordEvent{
		VoteDetails: VoteDetails{
			Voter:   "voter2",
			Choices: []Choice{"PHP", "Golang", "Haskell"},
			Time:    time.Now(),
		},
	}
	poll.ApplyEvents([]es.Event{
		pollCreatedEvent,
		completedEvent,
		failedEvent,
	})
	assert.Equal(t, 3, poll.Version(), "version error")
	assert.Equal(t, pollCreatedEvent.Title, poll._title, "Title错误")
	assert.Equal(t, pollCreatedEvent.Choices, poll._choices, "Choices错误")
	assert.Equal(t, []Choice{"Java", "Golang"}, poll._userChoices["voter1"], "投票记录信息错误")
}

func TestPollRestoreForErrorEvent(t *testing.T) {
	assert.Panics(t, func() {
		NewPoll().ApplyEvents([]es.Event{&struct{ es.WithGuid }{}})
	}, "restore error event must panic error")
}

func TestCheckPollApplyEvents(t *testing.T) {
	events := []es.Event{
		&PollCreatedEvent{},
		&VotePollCompletedBecauseOfVoteRecordEvent{},
		&VotePollFailedBecauseOfVoteRecordEvent{},
	}
	assert.NotPanics(t, func() { NewPoll().ApplyEvents(events) }, "Check Process All Event")
}

func TestPollCommand(t *testing.T) {
	guid := es.NewGuid()
	title := "title1"
	choices := []Choice{"PHP", "Java", "Golang", "Haskell", "Node.js"}
	voteDetails := VoteDetails{
		VoteRecordId: es.NewGuid(),
		PollId:       guid,
		Voter:        "adj",
		Choices:      []Choice{"PHP", "Java", "Golang"},
		Time:         time.Now(),
	}
	errorVoteDetails := VoteDetails{
		VoteRecordId: es.NewGuid(),
		PollId:       guid,
		Voter:        "adj",
		Choices:      []Choice{"PHP", "Java", ".Net"},
		Time:         time.Now(),
	}

	tests := []struct {
		poll    *Poll
		command es.Command
		event   es.Event
	}{
		{
			&Poll{},
			&CreatePollCommand{WithGuid: es.WithGuid{Guid: guid}, Title: title, Choices: choices},
			&PollCreatedEvent{WithGuid: es.WithGuid{Guid: guid}, Title: title, Choices: choices},
		},
		{
			&Poll{_title: title, _choices: choices},
			&VotePollBecauseOfVoteRecordCommand{WithGuid: es.WithGuid{Guid: guid}, VoteDetails: voteDetails},
			&VotePollCompletedBecauseOfVoteRecordEvent{WithGuid: es.WithGuid{Guid: guid}, VoteDetails: voteDetails},
		},
		{
			&Poll{_title: title, _choices: choices},
			&VotePollBecauseOfVoteRecordCommand{WithGuid: es.WithGuid{Guid: guid}, VoteDetails: errorVoteDetails},
			&VotePollFailedBecauseOfVoteRecordEvent{WithGuid: es.WithGuid{Guid: guid}, VoteDetails: errorVoteDetails},
		},
	}

	for _, v := range tests {
		assert.Equal(t, []es.Event{v.event}, v.poll.ProcessCommand(v.command))
	}
}

func TestPollCommand_Panic(t *testing.T) {
	tests := []struct {
		poll    *Poll
		command es.Command
	}{
		{
			&Poll{},
			&struct{ es.WithGuid }{},
		},
		{
			&Poll{},
			&VotePollBecauseOfVoteRecordCommand{},
		},
	}

	for _, v := range tests {
		assert.Panics(t, func() { v.poll.ProcessCommand(v.command) }, fmt.Sprintf("test panics error: command:%v", v.command))
	}
}
