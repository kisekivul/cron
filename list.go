package cron

import (
	"sync"
	"time"
)

// TaskList list of task
type TaskList struct {
	sync.Map
}

// List list of task
func (tl *TaskList) List() []TaskFace {
	var (
		l = make([]TaskFace, 0)
	)

	tl.Range(func(key, value interface{}) bool {
		l = append(l, value.(TaskFace))
		return true
	})
	return l
}

// Dict dict of task
func (tl *TaskList) Dict() map[string]TaskFace {
	var (
		m = make(map[string]TaskFace)
	)

	tl.Range(func(key, value interface{}) bool {
		m[key.(string)] = value.(TaskFace)
		return true
	})
	return m
}

// Len length of list
func (tl *TaskList) Len() int {
	return len(tl.List())
}

// Add add task
func (tl *TaskList) Add(name string, task TaskFace) {
	task.SetNext(time.Now().Local())
	tl.Store(name, task)
}

// Next set task next runing
func (tl *TaskList) Next(t time.Time) {
	tl.Range(func(key, value interface{}) bool {
		value.(TaskFace).SetNext(t)
		return true
	})
}
