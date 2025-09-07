package main

import (
	"fmt"
	"sort"
)

type Task struct {
	ID        string // Unique identifier for the task
	BurstTime int    // How long the task will run
	Priority  int    // Priority of the task (only used in Priority Scheduling)
}

func fcfs(tasks []Task) {
	for _, task := range tasks {
		fmt.Printf("Executing task %s for %d units of time.\n", task.ID, task.BurstTime)
	}
}

func sjf(tasks []Task) {
	// Sort tasks by burst time
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].BurstTime < tasks[j].BurstTime
	})

	for _, task := range tasks {
		fmt.Printf("Executing task %s for %d units of time.\n", task.ID, task.BurstTime)
	}
}

func priorityScheduling(tasks []Task) {
	// Sort tasks by priority (higher number means higher priority)
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Priority > tasks[j].Priority
	})

	for _, task := range tasks {
		fmt.Printf("Executing task %s with priority %d for %d units of time.\n", task.ID, task.Priority, task.BurstTime)
	}
}

func roundRobin(tasks []Task, quantum int) {
	queue := append([]Task{}, tasks...) // Copy the slice to avoid modifying the original
	for len(queue) > 0 {
		task := queue[0]
		queue = queue[1:]

		// Determine how long to execute this task
		timeToRun := task.BurstTime
		if timeToRun > quantum {
			timeToRun = quantum
		}

		fmt.Printf("Executing task %s for %d units of time.\n", task.ID, timeToRun)
		task.BurstTime -= timeToRun

		// If the task is not finished, put it back in the queue
		if task.BurstTime > 0 {
			queue = append(queue, task)
		}
	}
}

func multilevelQueueScheduling(tasks []Task) {
	// Split tasks into high-priority and low-priority queues
	highPriorityQueue := []Task{}
	lowPriorityQueue := []Task{}

	for _, task := range tasks {
		if task.Priority > 5 { // Arbitrary threshold for high priority
			highPriorityQueue = append(highPriorityQueue, task)
		} else {
			lowPriorityQueue = append(lowPriorityQueue, task)
		}
	}

	// Process high-priority tasks first (using Round Robin)
	fmt.Println("Executing high-priority tasks:")
	roundRobin(highPriorityQueue, 4) // 4 is the quantum for high priority

	// Then process low-priority tasks
	fmt.Println("Executing low-priority tasks:")
	roundRobin(lowPriorityQueue, 4)
}

func multilevelFeedbackQueue(tasks []Task) {
	// Step 1: Queue tasks in high priority or low priority
	highPriorityQueue := []Task{}
	lowPriorityQueue := []Task{}

	for _, task := range tasks {
		if task.Priority > 5 { // Arbitrary threshold
			highPriorityQueue = append(highPriorityQueue, task)
		} else {
			lowPriorityQueue = append(lowPriorityQueue, task)
		}
	}

	// Step 2: Process tasks in high-priority queue
	fmt.Println("Executing high-priority tasks:")
	roundRobin(highPriorityQueue, 4)

	// Step 3: If there are leftover tasks, put them in low-priority queue
	// and process them with longer time quantum.
	fmt.Println("Executing low-priority tasks:")
	roundRobin(lowPriorityQueue, 6)
}

type TaskWithDeadline struct {
	ID        string
	BurstTime int
	Deadline  int // Lower number means higher priority
}

func edf(tasks []TaskWithDeadline) {
	// Sort tasks by deadline (earliest deadline first)
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Deadline < tasks[j].Deadline
	})

	for _, task := range tasks {
		fmt.Printf("Executing task %s with deadline %d for %d units of time.\n", task.ID, task.Deadline, task.BurstTime)
	}
}

type RealTimeTask struct {
	ID        string
	BurstTime int
	Period    int // Shorter periods get higher priority
}

func rateMonotonicScheduling(tasks []RealTimeTask) {
	// Sort tasks by period (shortest period first)
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Period < tasks[j].Period
	})

	for _, task := range tasks {
		fmt.Printf("Executing real-time task %s with period %d for %d units of time.\n", task.ID, task.Period, task.BurstTime)
	}
}

func main() {
	// Initialize a set of tasks for testing
	tasks := []Task{
		{"T1", 4, 3},
		{"T2", 3, 2},
		{"T3", 5, 1},
	}

	// FCFS
	fmt.Println("First-Come, First-Served (FCFS):")
	fcfs(tasks)
	fmt.Println()

	// SJF
	fmt.Println("Shortest Job First (SJF):")
	sjf(tasks)
	fmt.Println()

	// Priority Scheduling
	fmt.Println("Priority Scheduling:")
	priorityScheduling(tasks)
	fmt.Println()

	// Round Robin
	fmt.Println("Round Robin Scheduling:")
	roundRobin(tasks, 2)
	fmt.Println()

	// Multilevel Queue Scheduling
	fmt.Println("Multilevel Queue Scheduling:")
	multilevelQueueScheduling(tasks)
	fmt.Println()

	// Multilevel Feedback Queue Scheduling
	fmt.Println("Multilevel Feedback Queue Scheduling:")
	multilevelFeedbackQueue(tasks)
	fmt.Println()

	// Initialize tasks for EDF and RMS (different structures)
	edfTasks := []TaskWithDeadline{
		{"T1", 4, 2},
		{"T2", 3, 1},
		{"T3", 5, 3},
	}

	// EDF
	fmt.Println("Earliest Deadline First (EDF):")
	edf(edfTasks)
	fmt.Println()

	// Initialize tasks for Rate Monotonic Scheduling
	rmsTasks := []RealTimeTask{
		{"T1", 4, 2},
		{"T2", 3, 1},
		{"T3", 5, 3},
	}

	// Rate Monotonic Scheduling
	fmt.Println("Rate Monotonic Scheduling (RMS):")
	rateMonotonicScheduling(rmsTasks)
}
