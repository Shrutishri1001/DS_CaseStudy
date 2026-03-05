package main

import (
	"fmt"
	"net/rpc"
	"os"
)

func main() {

	// STEP 1: Create log file
	logs := []string{
		"10.0.0.1 GET /product/101 200",
		"10.0.0.2 GET /product/102 200",
		"10.0.0.3 GET /product/101 200",
		"10.0.0.4 GET /product/103 200",
		"10.0.0.5 GET /product/101 200",
	}

	file, err := os.Create("web_access.log")

	if err != nil {
		fmt.Println("Error creating log file")
		return
	}

	for _, line := range logs {
		file.WriteString(line + "\n")
	}

	file.Close()

	fmt.Println("Log file created successfully")

	// STEP 2: Connect to NameNode
	masterIP := "10.198.182.146:9000" // CHANGE to your NameNode IP

	client, err := rpc.Dial("tcp", masterIP)

	if err != nil {
		fmt.Println("Cannot connect to NameNode:", err)
		return
	}

	fmt.Println("Connected to NameNode")

	// STEP 3: Call RPC function
	var reply PingReply

	err = client.Call("NameNode.StartProcessing", &PingArgs{}, &reply)

	if err != nil {
		fmt.Println("RPC call failed:", err)
		return
	}

	fmt.Println("Server reply:", reply.Status)

	client.Close()
}
