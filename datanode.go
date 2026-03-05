package main

import (
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strings"
)

type DataNode struct {
	ID          int
	currentLogs []string
}
type SnapshotState struct {
	NodeID int
	Logs   []string
}

func (d *DataNode) ReceiveMarker(marker SnapshotMarker, reply *bool) error {

	fmt.Println("Snapshot marker received")

	state := SnapshotState{
		NodeID: 1,
		Logs:   d.currentLogs,
	}

	fmt.Println("Local snapshot recorded:", state)

	*reply = true

	return nil
}

func (dn *DataNode) StoreLog(args *WriteArgs, reply *WriteReply) error {
	file, err := os.Create("replica.log")

	if err != nil {
		reply.Success = false
		return err
	}

	defer file.Close()

	file.Write(args.FileData)

	fmt.Println("Log replicated on DataNode")

	reply.Success = true

	return nil
}

func (dn *DataNode) ProcessChunk(args *MapArgs, reply *MapReply) error {

	result := make(map[string]int)

	for _, line := range args.Chunk {

		parts := strings.Fields(line)

		if len(parts) >= 3 && strings.Contains(parts[2], "/product/") {

			product := parts[2]
			result[product]++
		}
	}

	fmt.Println("Processed chunk. Result:", result)

	reply.Counts = result
	return nil
}

func (dn *DataNode) Ping(args *PingArgs, reply *PingReply) error {
	reply.Status = "ALIVE"
	return nil
}

func main() {

	rpc.Register(new(DataNode))

	listener, _ := net.Listen("tcp", ":8001")

	fmt.Println("DataNode running...")

	rpc.Accept(listener)
}
