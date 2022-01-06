package cron

import (
	"sort"
)

// Dict sort map for taskface
type Dict struct {
	Keys []string
	Vals []TaskFace
}

// NewMapSorter create new taskface map
func NewMapSorter(m map[string]TaskFace) *Dict {
	var (
		ms = &Dict{
			Keys: make([]string, 0, len(m)),
			Vals: make([]TaskFace, 0, len(m)),
		}
	)

	for k, v := range m {
		ms.Keys = append(ms.Keys, k)
		ms.Vals = append(ms.Vals, v)
	}
	return ms
}

// Sort sort taskface map
func (ms *Dict) Sort() {
	sort.Sort(ms)
}

func (ms *Dict) Len() int { return len(ms.Keys) }

func (ms *Dict) Less(i, j int) bool {
	if ms.Vals[i].GetNext().IsZero() {
		return false
	}
	if ms.Vals[j].GetNext().IsZero() {
		return true
	}
	return ms.Vals[i].GetNext().Before(ms.Vals[j].GetNext())
}

func (ms *Dict) Swap(i, j int) {
	ms.Vals[i], ms.Vals[j] = ms.Vals[j], ms.Vals[i]
	ms.Keys[i], ms.Keys[j] = ms.Keys[j], ms.Keys[i]
}
