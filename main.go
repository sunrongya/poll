package main

import (
	"fmt"
	ES "github.com/sunrongya/eventsourcing"
	"sync"
	"time"
)

func main() {
	var store = ES.NewInMemStore()
	wg := sync.WaitGroup{}
	wg.Add(1)

	ps := NewPollService(store)
	vs := NewVoteService(store)
	eventbus := ES.NewInternalEventBus(store)

	// 注册EventHandler/读模型Handler
	eh := NewEventHandler(ps.CommandChannel(), vs.CommandChannel())
	eventbus.RegisterHandlers(eh)

	go eventbus.HandleEvents()
	go ps.HandleCommands()
	go vs.HandleCommands()

	// 执行命令
	fmt.Printf("- 创建调查题目1\tOK\n")
	poll1 := ps.CreatePoll("喜欢哪几种语言？", []Choice{"PHP", "Java", "Golang", "Haskell", "Node.js"})
	fmt.Printf("- 创建调查题目2\tOK\n")
	poll2 := ps.CreatePoll("请选择你喜欢的数字？", []Choice{"1", "2", "3", "4", "5", "6", "7", "8", "9"})
	fmt.Printf("- 投票成功\tOK\n")
	vote1 := vs.VotePoll(poll1, "sry", []Choice{"Golang", "Haskell", "Node.js"}, time.Now())
	fmt.Printf("- 不能重复投票\tOK\n")
	vote2 := vs.VotePoll(poll1, "sry", []Choice{"PHP", "Haskell", "Node.js"}, time.Now())
	fmt.Printf("- 投票成功\tOK\n")
	vote3 := vs.VotePoll(poll1, "abc", []Choice{"PHP", "Java", "Golang"}, time.Now())
	fmt.Printf("- 投票成功\tOK\n")
	vote4 := vs.VotePoll(poll2, "sry", []Choice{"2", "3", "4"}, time.Now())
	fmt.Printf("- 投票失败：投票选项不能为空\tOK\n")
	vote5 := vs.VotePoll(poll2, "ccd", []Choice{}, time.Now())

	// 验证
	//wait and print
	go func() {
		time.Sleep(200 * time.Millisecond)
		printEvents(store.GetEvents(0, 100))
		fmt.Printf("-----------------\nAggregates:\n\n")
		fmt.Printf("%v\n------------------\n", ps.RestoreAggregate(poll1))
		fmt.Printf("%v\n------------------\n", ps.RestoreAggregate(poll2))
		fmt.Printf("%v\n------------------\n", vs.RestoreAggregate(vote1))
		fmt.Printf("%v\n------------------\n", vs.RestoreAggregate(vote2))
		fmt.Printf("%v\n------------------\n", vs.RestoreAggregate(vote3))
		fmt.Printf("%v\n------------------\n", vs.RestoreAggregate(vote4))
		fmt.Printf("%v\n------------------\n", vs.RestoreAggregate(vote5))

		wg.Done()
	}()

	wg.Wait()
}

func printEvents(events []ES.Event) {
	fmt.Printf("-----------------\nEvents after all operations:\n\n")
	for i, e := range events {
		fmt.Printf("%v: %T\n", i, e)
	}
}
