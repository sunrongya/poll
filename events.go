package main

import (
	es "github.com/sunrongya/eventsourcing"
)

// --------------
// Poll Events
// --------------

type PollCreatedEvent struct {
	es.WithGuid
	Title   string
	Choices []Choice
}

type VotePollCompletedBecauseOfVoteRecordEvent struct {
	es.WithGuid
	VoteDetails
}

type VotePollFailedBecauseOfVoteRecordEvent struct {
	es.WithGuid
	VoteDetails
}

// ------------------
// VoteRecord Events
// ------------------
type VoteRecordCreatedEvent struct {
	es.WithGuid
	VoteDetails
}

type VoteRecordCompletedEvent struct {
	es.WithGuid
	VoteDetails
}

type VoteRecordFailedEvent struct {
	es.WithGuid
	VoteDetails
}
