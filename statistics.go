package cron

import (
	"fmt"
	"sync"
	"time"
)

// Statistics struct
type Statistics struct {
	RequestURL        string
	RequestController string
	RequestNum        int64
	MinTime           time.Duration
	MaxTime           time.Duration
	TotalTime         time.Duration
}

// URLMap contains several statistics struct to log different data
type URLMap struct {
	lock        sync.RWMutex
	LengthLimit int //limit the urlmap's length if it's equal to 0 there's no limit
	urlmap      map[string]map[string]*Statistics
}

// AddStatistics add statistics task.
// it needs request method, request url, request controller and statistics time duration
func (m *URLMap) AddStatistics(requestMethod, requestURL, requestController string, requesttime time.Duration) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if method, ok := m.urlmap[requestURL]; ok {
		if s, ok := method[requestMethod]; ok {
			s.RequestNum++
			if s.MaxTime < requesttime {
				s.MaxTime = requesttime
			}
			if s.MinTime > requesttime {
				s.MinTime = requesttime
			}
			s.TotalTime += requesttime
		} else {
			nb := &Statistics{
				RequestURL:        requestURL,
				RequestController: requestController,
				RequestNum:        1,
				MinTime:           requesttime,
				MaxTime:           requesttime,
				TotalTime:         requesttime,
			}
			m.urlmap[requestURL][requestMethod] = nb
		}
	} else {
		if m.LengthLimit > 0 && m.LengthLimit <= len(m.urlmap) {
			return
		}

		var (
			methodmap = make(map[string]*Statistics)
		)
		methodmap[requestMethod] = &Statistics{
			RequestURL:        requestURL,
			RequestController: requestController,
			RequestNum:        1,
			MinTime:           requesttime,
			MaxTime:           requesttime,
			TotalTime:         requesttime,
		}
		m.urlmap[requestURL] = methodmap
	}
}

// GetMap put url statistics result in io.Writer
func (m *URLMap) GetMap() map[string]interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()

	var (
		resultLists [][]string
		fields      = []string{"requestUrl", "method", "times", "used", "max used", "min used", "avg used"}
		content     = make(map[string]interface{})
	)
	content["Fields"] = fields

	for k, v := range m.urlmap {
		for kk, vv := range v {
			result := []string{
				fmt.Sprintf("% -50s", k),
				fmt.Sprintf("% -10s", kk),
				fmt.Sprintf("% -16d", vv.RequestNum),
				fmt.Sprintf("%d", vv.TotalTime),
				fmt.Sprintf("% -16s", toS(vv.TotalTime)),
				fmt.Sprintf("%d", vv.MaxTime),
				fmt.Sprintf("% -16s", toS(vv.MaxTime)),
				fmt.Sprintf("%d", vv.MinTime),
				fmt.Sprintf("% -16s", toS(vv.MinTime)),
				fmt.Sprintf("%d", time.Duration(int64(vv.TotalTime)/vv.RequestNum)),
				fmt.Sprintf("% -16s", toS(time.Duration(int64(vv.TotalTime)/vv.RequestNum))),
			}
			resultLists = append(resultLists, result)
		}
	}
	content["Data"] = resultLists
	return content
}

// GetMapData return all mapdata
func (m *URLMap) GetMapData() []map[string]interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()

	var (
		resultLists []map[string]interface{}
	)

	for k, v := range m.urlmap {
		for kk, vv := range v {
			resultLists = append(resultLists,
				map[string]interface{}{
					"request_url": k,
					"method":      kk,
					"times":       vv.RequestNum,
					"total_time":  toS(vv.TotalTime),
					"max_time":    toS(vv.MaxTime),
					"min_time":    toS(vv.MinTime),
					"avg_time":    toS(time.Duration(int64(vv.TotalTime) / vv.RequestNum)),
				},
			)
		}
	}
	return resultLists
}

// StatisticsMap hosld global statistics data map
var (
	StatisticsMap *URLMap
)

func init() {
	StatisticsMap = &URLMap{
		urlmap: make(map[string]map[string]*Statistics),
	}
}
