/*
Copyright 2017 Samsung SDSA CNCT

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package queue

import (
	"sync"
	"time"
)

// A dead simple work queue with time stamps for queued and running duration.
// Currently work queue acts as a singleton within the system to force all
// work in to a single queuing mechanism.  It does not define a type, all of
// the queue data is local in scope thereby forcing the use of the methods
// below. All of this could very easily be converted to a type with methods.

// TaskStatus - tasks status
type TaskStatus int

const (
	// NotFound status
	NotFound TaskStatus = iota
	// Queued status - waiting for processing
	Queued
	// Running status - task is being processed
	Running
	// Deleted status - task was removed from queue
	Deleted
)

// Task describes the queued request
type Task struct {
	// ID - Required by several methods for managing the task
	ID     uint
	timeIn time.Time
	timeUp time.Time
}

// NewTask creates an task request for processing.  The Task is not queued
// until Submit() is called for the Task.
func NewTask() *Task {
	taskID++
	return &Task{
		ID: taskID,
	}
}

// All work happens on a single system queue, the only locking is done when the
// underlying slice is being changed, and not for status operations (for now).
var mutex = &sync.Mutex{}
var workQueue []*Task
var taskID uint
var queued time.Duration
var running time.Duration

// TaskCount returns the number of tasks currently queued for processing
func TaskCount() int {
	return len(workQueue)
}

// QueuedDuration returns the duration (Submit to Done) of the last tasks processed
func QueuedDuration() time.Duration {
	return queued
}

// RunningDuration returns the duration (Running to Done) of the last tasks processed
func RunningDuration() time.Duration {
	return running
}

// Submit - enqueue a task for processing - returns the current number of tasks
// in the queue (including the new task just submitted)
func Submit(task *Task) int {
	task.timeIn = time.Now()
	mutex.Lock()
	defer mutex.Unlock()
	workQueue = append(workQueue, task)
	return len(workQueue)
}

// Status - return the status of the task associated with id, or NotFound
func Status(id uint) TaskStatus {
	for _, t := range workQueue {
		if id == t.ID {
			if t.timeUp.IsZero() {
				return Queued
			}
			return Running
		}
	}
	return NotFound
}

func queueDelete(index int) {
	mutex.Lock()
	defer mutex.Unlock()
	// delete ensuring *Task can be garbage collected
	copy(workQueue[index:], workQueue[index+1:])
	workQueue[len(workQueue)-1] = nil
	workQueue = workQueue[:len(workQueue)-1]
}

// Delete - remove a task from the work queue if, and only if, it has not yet
// been started.  Once the task has been started, it has to run to completion.
// Deleting a task from the work queue does not update duration times.
func Delete(id uint) TaskStatus {
	for i, t := range workQueue {
		if id == t.ID {
			if t.timeUp.IsZero() {
				queueDelete(i)
				return Deleted
			}
			return Running
		}
	}
	return NotFound
}

// Started - update the first task in the work queue if, and only if, it
// hasn't already been started, to indicate that it is now running.
func Started() bool {
	// don't take the lock to update the timestamp, not changing queue slice
	if !workQueue[0].timeUp.IsZero() {
		return false
	}
	workQueue[0].timeUp = time.Now()
	return true
}

// Done - remove's only the first task from the work queue if, and only if, it
// has already been Started().  Done delete's the taks and calculates durations
func Done() bool {
	if workQueue[0].timeUp.IsZero() {
		return false
	}
	t := time.Now()
	queued = t.Sub(workQueue[0].timeIn)
	running = t.Sub(workQueue[0].timeUp)
	queueDelete(0)
	return true
}
