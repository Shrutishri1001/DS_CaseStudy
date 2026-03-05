package main

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
