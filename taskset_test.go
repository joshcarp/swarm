package swarm

import "testing"

func TestWeighingTaskSetWithSingleTask(t *testing.T) {
	ts := NewWeighingTaskSet()

	taskAIsRun := false
	taskA := &Task{
		Namef:   "A",
		Weightf: 1,
		Fn: func() {
			taskAIsRun = true
		},
	}
	ts.AddTask(taskA)

	if ts.GetTask(0).Namef != "A" {
		t.Error("Expecting A, but got ", ts.GetTask(0).Namef)
	}
	if ts.GetTask(1) != nil {
		t.Error("Out of bound, should return nil")
	}
	if ts.GetTask(-1) != nil {
		t.Error("Out of bound, should return nil")
	}

	ts.Run()

	if !taskAIsRun {
		t.Error("Task A should be run")
	}
}

func TestWeighingTaskSetWithTwoTasks(t *testing.T) {
	ts := NewWeighingTaskSet()
	taskA := &Task{
		Namef:   "A",
		Weightf: 1,
	}
	taskB := &Task{
		Namef:   "B",
		Weightf: 2,
	}
	ts.AddTask(taskA)
	ts.AddTask(taskB)

	if ts.GetTask(0).Namef != "A" {
		t.Error("Expecting A, but got ", ts.GetTask(0).Namef)
	}
	if ts.GetTask(1).Namef != "B" {
		t.Error("Expecting B, but got ", ts.GetTask(1).Namef)
	}
}

func TestWeighingTaskSetGetTaskWithThreeTasks(t *testing.T) {
	ts := NewWeighingTaskSet()
	taskA := &Task{
		Namef:   "A",
		Weightf: 1,
	}
	taskB := &Task{
		Namef:   "B",
		Weightf: 2,
	}
	taskC := &Task{
		Namef:   "C",
		Weightf: 3,
	}
	ts.AddTask(taskA)
	ts.AddTask(taskB)
	ts.AddTask(taskC)

	if ts.GetTask(0).Namef != "A" {
		t.Error("Expecting A, but got ", ts.GetTask(0).Namef)
	}
	if ts.GetTask(1).Namef != "B" {
		t.Error("Expecting B, but got ", ts.GetTask(1).Namef)
	}
	if ts.GetTask(2).Namef != "B" {
		t.Error("Expecting B, but got ", ts.GetTask(2).Namef)
	}
	if ts.GetTask(3).Namef != "C" {
		t.Error("Expecting C, but got ", ts.GetTask(3).Namef)
	}
	if ts.GetTask(4).Namef != "C" {
		t.Error("Expecting C, but got ", ts.GetTask(4).Namef)
	}
	if ts.GetTask(5).Namef != "C" {
		t.Error("Expecting C, but got ", ts.GetTask(5).Namef)
	}
}
