package main

import (
	//"fmt"
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
	poll.HandlePollCreatedEvent(pollCreatedEvent)
	poll.HandleVotePollCompletedBecauseOfVoteRecordEvent(completedEvent)
	poll.HandleVotePollFailedBecauseOfVoteRecordEvent(failedEvent)

	assert.Equal(t, pollCreatedEvent.Title, poll._title, "Title错误")
	assert.Equal(t, pollCreatedEvent.Choices, poll._choices, "Choices错误")
	assert.Equal(t, []Choice{"Java", "Golang"}, poll._userChoices["voter1"], "投票记录信息错误")
}

func TestCreatePollCommand(t *testing.T) {
	title := "title1"
	choices := []Choice{"PHP", "Java", "Golang", "Haskell", "Node.js"}
	command := &CreatePollCommand{Title: title, Choices: choices}
	events := []es.Event{&PollCreatedEvent{Title: title, Choices: choices}}

	assert.Equal(t, events, new(Poll).ProcessCreatePollCommand(command), "创建调查题目返回的事件有误")
}

func VotePollBecauseOfVoteRecordCommand2Completed(t *testing.T) {
	title := "title1"
	choices := []Choice{"PHP", "Java", "Golang", "Haskell", "Node.js"}
	voteDetails := VoteDetails{
		VoteRecordId: es.NewGuid(),
		PollId:       es.NewGuid(),
		Voter:        "adj",
		Choices:      []Choice{"PHP", "Java", "Golang"},
		Time:         time.Now(),
	}
	poll := &Poll{_title: title, _choices: choices}
	command := &VotePollBecauseOfVoteRecordCommand{VoteDetails: voteDetails}
	events := []es.Event{&VotePollCompletedBecauseOfVoteRecordEvent{VoteDetails: voteDetails}}

	assert.Equal(t, events, poll.ProcessVotePollBecauseOfVoteRecordCommand(command), "调查题目返回投票成功事件有误")
}

func VotePollBecauseOfVoteRecordCommand2Failured(t *testing.T) {
	title := "title1"
	choices := []Choice{"PHP", "Java", "Golang", "Haskell", "Node.js"}
	errorVoteDetails := VoteDetails{
		VoteRecordId: es.NewGuid(),
		PollId:       es.NewGuid(),
		Voter:        "adj",
		Choices:      []Choice{"PHP", "Java", ".Net"},
		Time:         time.Now(),
	}
	poll := &Poll{_title: title, _choices: choices}
	command := &VotePollBecauseOfVoteRecordCommand{VoteDetails: errorVoteDetails}
	events := []es.Event{&VotePollFailedBecauseOfVoteRecordEvent{VoteDetails: errorVoteDetails}}

	assert.Equal(t, events, poll.ProcessVotePollBecauseOfVoteRecordCommand(command), "调查题目返回投票失败事件有误")
}

func TestPollCommand_Panic(t *testing.T) {
	assert.Panics(t,
		func() {
			new(Poll).ProcessVotePollBecauseOfVoteRecordCommand(&VotePollBecauseOfVoteRecordCommand{})
		},
		"调查题目内容为空，不能进行投票",
	)
}
