package main

import (
	"bufio"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"sync"
	"time"
)

// DataNodes in the cluster
var datanodes = []string{
	"10.198.182.149:8001",
	"10.198.182.127:8001",
}

// ALGORITHM 1 — Centralized Mutual Exclusion
type NameNode struct {
	mu           sync.Mutex
	currentLease string
}
type SnapshotState struct {
	NodeID int
	Logs   []string
}

// RPC function called by client
func (n *NameNode) StartProcessing(args *PingArgs, reply *PingReply) error {

	fmt.Println("Client connected. Starting log processing...")

	chunks := splitFile()

	nodes := datanodes

	results := []map[string]int{}

	success := 0

	for i, node := range nodes {

		client, err := rpc.Dial("tcp", node)

		if err != nil {
			fmt.Println("Could not connect to DataNode:", node)
			continue
		}

		fmt.Println("Connected to DataNode:", node)
		var mapReply MapReply
		mapArgs := MapArgs{
			Chunk: chunks[i],
		}

		err = client.Call("DataNode.ProcessChunk", &mapArgs, &mapReply)
		if err != nil {
			fmt.Println("RPC Error:", err)
		}
		results = append(results, mapReply.Counts)
		success++
		client.Close()
	}

	// Algorithm 2 — Quorum-Based Consensus
	if success >= 2 {

		fmt.Println("QUORUM ACHIEVED")
		fmt.Println("Processing committed")

	} else {

		fmt.Println("QUORUM FAILED")
	}

	final := make(map[string]int)

	for _, r := range results {

		for k, v := range r {

			final[k] += v
		}
	}

	fmt.Println("\nFINAL PRODUCT ACCESS COUNT")
	fmt.Println("---------------------------")

	for k, v := range final {
		fmt.Println(k, ":", v)
	}
	data, err := os.ReadFile("web_access.log")

	if err == nil {
		replicateLogs(&WriteArgs{
			FileData: data,
		})
		fmt.Println("Logs replicated to DataNodes")
	}

	/* Start Global Snapshot */
	startSnapshot(datanodes)

	reply.Status = "Processing Complete"

	return nil
}

// Function to split log file
func splitFile() [][]string {

	file, err := os.Open("web_access.log")

	if err != nil {
		fmt.Println("Error opening log file")
		return nil
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	mid := len(lines) / 2

	return [][]string{
		lines[:mid],
		lines[mid:],
	}
}

// Lease Request
func (nn *NameNode) RequestLease(args *LeaseArgs, reply *LeaseReply) error {

	nn.mu.Lock()
	defer nn.mu.Unlock()

	if nn.currentLease == "" {

		nn.currentLease = args.ClientID
		reply.Granted = true

		fmt.Println("Lease granted to", args.ClientID)

	} else {

		reply.Granted = false
		fmt.Println("Client waiting:", args.ClientID)

	}

	return nil
}

// Lease Release
func (nn *NameNode) ReleaseLease(args *LeaseArgs, reply *LeaseReply) error {

	nn.mu.Lock()
	defer nn.mu.Unlock()

	if nn.currentLease == args.ClientID {

		nn.currentLease = ""

		fmt.Println("Lease released by", args.ClientID)

	}

	return nil
}

// Algorithm 3 — Heartbeat Failure Detection
func heartbeat(nodes []string) {

	for {

		for _, node := range nodes {
			conn, err := net.DialTimeout("tcp", node, 1*time.Second)
			if err != nil {
				fmt.Println("Node DOWN:", node)
			} else {
				fmt.Println("Node ALIVE:", node)
				conn.Close()
			}
		}

		time.Sleep(3 * time.Second)
	}
}

// Algorithm 4 — NameNode Replication
func replicateLogs(args *WriteArgs) {

	for _, node := range datanodes {

		go func(addr string) {

			client, err := rpc.Dial("tcp", addr)

			if err != nil {
				fmt.Println("Replication failed to:", addr)
				return
			}

			defer client.Close()

			var reply WriteReply

			err = client.Call("DataNode.StoreLog", args, &reply)

			if err != nil {
				fmt.Println("Replication RPC error:", err)
			}

		}(node)
	}
}

// Algorithm 5 — Bully Leader Election
func leaderElection(masterID int, leaderID int, peerAddress string) {

	for {

		conn, err := net.DialTimeout("tcp", peerAddress, 1*time.Second)

		if err != nil {

			fmt.Println("\nHeartbeat timeout detected")
			fmt.Println("Leader NameNode is down")

			fmt.Println("Starting Bully Election...")

			if masterID > leaderID {

				fmt.Println("This node has higher ID")
				fmt.Println("Becoming new leader")

			}

		} else {

			conn.Close()
		}

		time.Sleep(2 * time.Second)
	}
}

func main() {

	namenode := new(NameNode)
	rpc.Register(namenode)
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	fmt.Println("NameNode running on port 9000")
	go heartbeat(datanodes)

	/* Start leader election monitor */
	go leaderElection(1, 0, "10.198.182.149:9000")

	for {

		conn, err := listener.Accept()

		if err != nil {
			continue
		}

		go rpc.ServeConn(conn)
	}
}

// snapshot algorithm
func startSnapshot(nodes []string) {
	fmt.Println("\nStarting Chandy-Lamport Snapshot")
	for _, node := range nodes {
		client, err := rpc.Dial("tcp", node)
		if err != nil {
			continue
		}
		var reply bool
		client.Call("DataNode.ReceiveMarker", SnapshotMarker{SnapshotID: 1}, &reply)
		client.Close()
	}
}
