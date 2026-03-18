package main

//Algorithm 1 — Centralized Mutual Exclusion
type LeaseArgs struct {
	ClientID string
}

type LeaseReply struct {
	Granted bool
}

type MapArgs struct {
	Chunk []string
}

type MapReply struct {
	Counts map[string]int
}

type ReduceArgs struct {
	PartialResults []map[string]int
}

type ReduceReply struct {
	FinalResult map[string]int
}

type PingArgs struct{}
type PingReply struct {
	Status string
}

type ElectionArgs struct {
	CandidateID int
}

type ElectionReply struct {
	OK bool
}
type WriteArgs struct {
	FileData []byte
}

type WriteReply struct {
	Success bool
}
type SnapshotMarker struct {
	SnapshotID int
}
