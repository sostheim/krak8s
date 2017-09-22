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

import "testing"
import "time"

func TestNewTask(t *testing.T) {
	iterations := 10000
	taskid := NewTask().ID
	for i := 0; i < iterations; i++ {
		task := NewTask()
		if task.ID != taskid+1 {
			t.Errorf("TestNewTask() have: task.ID: %d, want taskid: %d", task.ID, taskid+1)
		}
		taskid = task.ID
	}
}

func TestSubmitStartDoneDeleteStatus(t *testing.T) {
	task1 := NewTask()
	depth := Submit(task1)
	if depth != 1 {
		t.Errorf("TestSubmit() have depth: %d, want depth: 1", depth)
	}
	time.Sleep(500 * time.Millisecond)
	Started() // task 1
	task2 := NewTask()
	depth = Submit(task2)
	if depth != 2 {
		t.Errorf("TestSubmit() have depth: %d, want depth: 2", depth)
	}
	if Status(task1.ID) != Running {
		t.Errorf("TestSubmit() have status: %d, want: %d", Status(task1.ID), Running)
	}
	if Status(task2.ID) != Queued {
		t.Errorf("TestSubmit() have status: %d, want: %d", Status(task2.ID), Queued)
	}
	if Status(9000) != NotFound {
		t.Errorf("TestSubmit() have status: %d, want: %d", Status(9000), NotFound)
	}
	time.Sleep(500 * time.Millisecond)
	Done() // task 1
	task3 := NewTask()
	depth = Submit(task3)
	if depth != 2 {
		t.Errorf("TestSubmit() have depth: %d, want depth: 2", depth)
	}
	if Status(task1.ID) != NotFound {
		t.Errorf("TestSubmit() have status: %d, want: %d", Status(task1.ID), NotFound)
	}
	if Status(task2.ID) != Queued {
		t.Errorf("TestSubmit() have status: %d, want: %d", Status(task2.ID), Running)
	}
	if Status(task3.ID) != Queued {
		t.Errorf("TestSubmit() have status: %d, want: %d", Status(task3.ID), Queued)
	}
	Started() // task 2
	if Status(task2.ID) != Running {
		t.Errorf("TestSubmit() have status: %d, want: %d", Status(task2.ID), Running)
	}
	time.Sleep(500 * time.Millisecond)
	Done() // task 2
	if Status(task2.ID) != NotFound {
		t.Errorf("TestSubmit() have status: %d, want: %d", Status(task2.ID), Running)
	}
	depth = TaskCount()
	if depth != 1 {
		t.Errorf("TestSubmit() have depth: %d, want depth: 1", depth)
	}
	Started() // task 3
	if Status(task3.ID) != Running {
		t.Errorf("TestSubmit() have status: %d, want: %d", Status(task2.ID), Running)
	}
	time.Sleep(500 * time.Millisecond)
	Done() // task 3
	if Status(task3.ID) != NotFound {
		t.Errorf("TestSubmit() have status: %d, want: %d", Status(task2.ID), Running)
	}
	depth = TaskCount()
	if depth != 0 {
		t.Errorf("TestSubmit() have depth: %d, want depth: 0", depth)
	}

}
