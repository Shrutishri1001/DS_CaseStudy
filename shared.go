package main

type MapArgs struct {
	Chunk []string
}
type WriteArgs struct {
	FileData []byte
}

type WriteReply struct {
	Success bool
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
type SnapshotMarker struct {
	SnapshotID int
}
