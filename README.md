### Distributed Log Processing System using RPC, MapReduce and Distributed Algorithms

---

### Overview

This project demonstrates a **distributed log processing system built in Go (Golang)** that simulates key concepts of modern distributed systems such as:

- Distributed storage  
- Parallel data processing  
- Fault tolerance  
- Leader election  
- Consensus mechanisms  

The system processes **web server access logs across multiple machines** using a simplified **MapReduce-style architecture**.

Instead of processing logs on a single machine, the workload is **distributed across DataNodes**, coordinated by a **NameNode**, and triggered by a **Client**.

---

### System Architecture

The system consists of three main components:

| Component | Role |
|-----------|------|
| Client | Uploads logs and triggers processing |
| NameNode | Coordinator that distributes tasks |
| DataNodes | Worker nodes that process log chunks |

#### Architecture Flow

             CLIENT
                |
                | RPC Request
                v
           NAMENODE
    (Coordinator + Reducer)
                |
      -----------------------
      |                     |
      v                     v
  DATANODE 1           DATANODE 2
  (Mapper)              (Mapper)
      |                     |
      -----------+-----------
                  |
                  v
             REDUCER
           (NameNode)
                  |
                  v
            FINAL RESULT

---

### Project Objective

The goal of this project is to simulate **distributed log analysis** where:

- Web access logs are split into chunks
- Each chunk is processed on different machines
- Results are aggregated to compute product access statistics

#### Example Log Entry
10.0.0.1 GET /product/101 200
10.0.0.2 GET /product/102 200
10.0.0.3 GET /product/101 200


#### Example Output
/product/101 : 2
/product/102 : 1


This shows how many times each product page was accessed.

---

### Technologies Used

- Go (Golang)
- RPC (Remote Procedure Calls)
- TCP Networking
- Distributed Algorithms
- MapReduce processing model

---

### Key Distributed System Concepts Implemented

#### 1. RPC-based Distributed Communication

Nodes communicate using **Go RPC**, allowing one machine to execute functions on another machine over the network.

Example RPC call:
client.Call("DataNode.ProcessChunk", &mapArgs, &mapReply)


This allows the NameNode to send log chunks to DataNodes for processing.

---

#### 2. MapReduce Style Processing

The log processing follows a simplified **MapReduce pipeline**.

##### Map Phase (DataNodes)

Each DataNode processes a chunk of logs and produces intermediate results.

Example mapper output:
/product/101 1
/product/102 1
/product/101 1


##### Reduce Phase (NameNode)

The NameNode aggregates mapper outputs to compute final counts.

Example reduce output:
/product/101 : 2
/product/102 : 1


---

#### 3. Centralized Mutual Exclusion

Only one client can modify the system at a time using a **lease mechanism**.

This prevents concurrent write conflicts.

Implemented using:
RequestLease()
ReleaseLease()


---

#### 4. Quorum-based Consensus

The system commits processing only if responses from a **majority of DataNodes** are received.

Example:
QUORUM ACHIEVED
Processing committed


This ensures reliability even if some nodes fail.

---

#### 5. Heartbeat Failure Detection

The NameNode continuously checks whether DataNodes are alive.
Example:
Node ALIVE: 10.198.182.149
Node DOWN: 10.198.182.127


This enables detection of node failures in the cluster.

---

#### 6. Log Replication

Log files are replicated across DataNodes for **fault tolerance**.

If one node fails, another replica can serve the data.

---

#### 7. Bully Leader Election Algorithm

The system implements the **Bully algorithm** for leader election.

If the current leader fails:

1. Nodes detect the failure  
2. An election is initiated  
3. The node with the highest ID becomes the new leader  

Example output:
Heartbeat timeout detected
Leader NameNode is down
Starting Bully Election...
This node has higher ID
Becoming new leader


---

#### 8. Chandy-Lamport Distributed Snapshot

The system triggers a **global snapshot** to capture the current state of distributed nodes.
Starting Chandy-Lamport Snapshot

This technique is used to analyze **consistent system states in distributed environments**.

---

### Concurrent Execution with Goroutines

The system uses **Goroutines** for concurrency.

Examples:
go heartbeat(datanodes)
go rpc.ServeConn(conn)


#### Benefits

- Parallel processing
- Non-blocking network operations
- Efficient system monitoring

---

### How the System Works

#### Step 1 – Client starts processing

The client sends an RPC request to the NameNode.
StartProcessing()

#### Step 2 – NameNode splits log file

The log file is divided into chunks.
Chunk1 → DataNode1
Chunk2 → DataNode2


#### Step 3 – DataNodes run mapper

Each DataNode processes its chunk and returns counts.

#### Step 4 – NameNode performs reduce

The NameNode aggregates mapper results to produce the final output.

#### Step 5 – Logs replicated

The processed logs are replicated across nodes.

#### Step 6 – Snapshot initiated

A distributed snapshot is triggered for system state capture.

---

### Advantages of the System

- Parallel processing of large log datasets
- Fault tolerance through replication
- Leader election for high availability
- Distributed coordination using RPC
- Scalable architecture

---

### Example Output
Client connected. Starting log processing...

Connected to DataNode: 10.198.182.149
Connected to DataNode: 10.198.182.127

QUORUM ACHIEVED
Processing committed

FINAL PRODUCT ACCESS COUNT

/product/101 : 2
/product/102 : 1

Logs replicated to DataNodes
Starting Chandy-Lamport Snapshot



