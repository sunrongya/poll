package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	es "github.com/sunrongya/eventsourcing"
	"testing"
)

func TestGoodsReadModel(t *testing.T) {
	readRepository := es.NewMemoryReadRepository()
	pollProjector := NewPollProjector(readRepository)

	// 调查题目创建
	createdEvents := []*PollCreatedEvent{
		&PollCreatedEvent{WithGuid: es.WithGuid{es.NewGuid()}, Title: "mangos", Choices: []Choice{"1", "2", "3", "4"}},
		&PollCreatedEvent{WithGuid: es.WithGuid{es.NewGuid()}, Title: "apple", Choices: []Choice{"5", "6", "7", "8"}},
	}

	for _, event := range createdEvents {
		pollProjector.HandlePollCreatedEvent(event)
	}

	// 调查题目创建验证
	for _, event := range createdEvents {
		i, err := readRepository.Find(event.GetGuid())
		assert.NoError(t, err, fmt.Sprintf("读取调查题目创建[%s]信息错误", event.Title))
		poll := i.(*RPoll)

		assert.Equal(t, event.GetGuid(), poll.Id, "ID 不相等")
		assert.Equal(t, event.Title, poll.Title, "Title 不相等")
		assert.Equal(t, event.Choices, poll.Choices, "Price 不相等")
		assert.Equal(t, 0, len(poll.ChoiceStat), "还没有选项信息")
	}

	// 投票
	votedEvents := []*VotePollCompletedBecauseOfVoteRecordEvent{
		&VotePollCompletedBecauseOfVoteRecordEvent{
			WithGuid:    es.WithGuid{createdEvents[1].GetGuid()},
			VoteDetails: VoteDetails{Choices: []Choice{"5", "6", "7"}},
		},
		&VotePollCompletedBecauseOfVoteRecordEvent{
			WithGuid:    es.WithGuid{createdEvents[1].GetGuid()},
			VoteDetails: VoteDetails{Choices: []Choice{"5", "6", "8"}},
		},
		&VotePollCompletedBecauseOfVoteRecordEvent{
			WithGuid:    es.WithGuid{createdEvents[1].GetGuid()},
			VoteDetails: VoteDetails{Choices: []Choice{"5", "7"}},
		},
	}
	for _, event := range votedEvents {
		pollProjector.HandleVotePollCompletedBecauseOfVoteRecordEvent(event)
	}

	// 投票统计验证
	stat := map[Choice]int{"5": 3, "6": 2, "7": 2, "8": 1}
	for i, event := range createdEvents {
		iPoll, err := readRepository.Find(event.GetGuid())
		assert.NoError(t, err, fmt.Sprintf("读取调查题目创建[%s]信息错误", event.Title))
		poll := iPoll.(*RPoll)

		assert.Equal(t, event.GetGuid(), poll.Id, "ID 不相等")
		assert.Equal(t, event.Title, poll.Title, "Title 不相等")
		assert.Equal(t, event.Choices, poll.Choices, "Price 不相等")
		dstStat := stat
		if i == 0 {
			dstStat = map[Choice]int{}
		}
		assert.Equal(t, dstStat, poll.ChoiceStat, "调查统计信息错误")
	}
}
