package main

import (
	es "github.com/sunrongya/eventsourcing"
)

// --------------
// Poll Commands
// --------------

type CreatePollCommand struct {
	es.WithGuid
	Title   string
	Choices []Choice
}

type VotePollBecauseOfVoteRecordCommand struct {
	es.WithGuid
	VoteDetails
}

// --------------------
// VoteRecord Commands
// --------------------
type CreateVoteRecordCommand struct {
	es.WithGuid
	VoteDetails
}

type CompleteVoteRecordCommand struct {
	es.WithGuid
	VoteDetails
}

type FailVoteRecordCommand struct {
	es.WithGuid
	VoteDetails
}
