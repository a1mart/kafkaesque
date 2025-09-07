package main

import (
	"fmt"
	"time"
)

// Process represents a process in the system
type Process struct {
	ID            int // Process ID
	BurstTime     int // Time required for execution
	RemainingTime int // Remaining time to execute
}

// RoundRobinScheduler schedules processes using the round robin algorithm
type RoundRobinScheduler struct {
	Queue   []*Process
	Quantum int // Time slice for each process
}

// NewRoundRobinScheduler creates a new RoundRobinScheduler
func NewRoundRobinScheduler(quantum int) *RoundRobinScheduler {
	return &RoundRobinScheduler{
		Queue:   make([]*Process, 0),
		Quantum: quantum,
	}
}

// AddProcess adds a process to the queue
func (rr *RoundRobinScheduler) AddProcess(p *Process) {
	rr.Queue = append(rr.Queue, p)
}

// Schedule runs the round robin scheduling
func (rr *RoundRobinScheduler) Schedule() {
	for len(rr.Queue) > 0 {
		// Get the first process in the queue
		process := rr.Queue[0]
		// Remove it from the queue
		rr.Queue = rr.Queue[1:]

		// Print the process starting to execute
		fmt.Printf("Executing Process %d (Remaining time: %d)\n", process.ID, process.RemainingTime)

		// Process executes for the time quantum or until it finishes
		if process.RemainingTime > rr.Quantum {
			process.RemainingTime -= rr.Quantum
			// Re-add the process to the end of the queue if it's not finished
			rr.Queue = append(rr.Queue, process)
			fmt.Printf("Process %d not finished, remaining time: %d\n", process.ID, process.RemainingTime)
		} else {
			// Process finishes
			fmt.Printf("Process %d completed!\n", process.ID)
			process.RemainingTime = 0
		}

		// Simulate the process running by sleeping for a short time
		time.Sleep(time.Millisecond * 500) // Simulate CPU time slice
	}
}

func main() {
	// Create a round robin scheduler with a time quantum of 3 units
	scheduler := NewRoundRobinScheduler(3)

	// Add some processes with burst times
	scheduler.AddProcess(&Process{ID: 1, BurstTime: 5, RemainingTime: 5})
	scheduler.AddProcess(&Process{ID: 2, BurstTime: 7, RemainingTime: 7})
	scheduler.AddProcess(&Process{ID: 3, BurstTime: 2, RemainingTime: 2})

	// Start the round robin scheduling
	scheduler.Schedule()
}
