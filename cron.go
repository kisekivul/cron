package cron

import (
	"time"
)

var (
	// TASKS all tasks
	TASKS   *TaskList
	stopped chan bool
	changed chan bool
	running bool
)

const (
	// Set the top bit if a star was included in the expression.
	starBit = 1 << 63
)

func init() {
	TASKS = new(TaskList)
	stopped = make(chan bool)
	changed = make(chan bool)
}

// NewTask add new task with name, time and func
func NewTask(tname string, spec string, f TaskFunc) *Task {
	task := &Task{
		Taskname: tname,
		DoFunc:   f,
		ErrLimit: 100,
		SpecStr:  spec,
	}
	task.SetCron(spec)
	return task
}

// StartTask start all tasks
func StartTask() {
	if running {
		//If already started， no need to start another goroutine.
		return
	}
	running = true
	// start running
	go start()
}

func start() {
	var (
		now = time.Now().Local()
	)
	TASKS.Next(now)

	for {
		var (
			dict      = TASKS.Dict()
			sortList  = NewMapSorter(dict)
			effective time.Time
		)
		sortList.Sort()

		if len(dict) == 0 || sortList.Vals[0].GetNext().IsZero() {
			// If there are no entries yet, just sleep - it still handles new entries
			// and stopped requests.
			effective = now.AddDate(10, 0, 0)
		} else {
			effective = sortList.Vals[0].GetNext()
		}

		select {
		case now = <-time.After(effective.Sub(now)):
			// Run every entry whose next time was this effective time.
			for _, e := range sortList.Vals {
				if e.GetNext() != effective {
					break
				}
				go e.Run()
				e.SetPrev(e.GetNext())
				e.SetNext(effective)
			}
			continue
		case <-changed:
			now = time.Now().Local()
			TASKS.Next(now)
			continue
		case <-stopped:
			return
		}
	}
}

// StopTask stopped all tasks
func StopTask() {
	if running {
		running = false
		stopped <- true
	}
}

// AddTask add task with name
func AddTask(name string, task TaskFace) {
	TASKS.Add(name, task)
	if running {
		changed <- true
	}
}

// DeleteTask delete task with name
func DeleteTask(name string) {
	TASKS.Delete(name)
	if running {
		changed <- true
	}
}
