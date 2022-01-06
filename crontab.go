package cron

import (
	"strconv"
	"strings"
)

const (
	ASTERISK = "*"
	BLANK    = " "
	SLASH    = "/"
)

func Crontab(i string) string {
	var (
		strs []string
		temp = []string{}
	)

	for _, v := range strings.Split(i, BLANK) {
		if strings.Compare(v, "") != 0 {
			strs = append(strs, strings.TrimSpace(v))
		}
	}
	// convert
	switch len(strs) {
	case 1:
		t, _ := strconv.Atoi(strs[0])
		if t == 0 {
			return ""
		}
		t = t * 60
		// interval kind
		var (
			day    = t / 3600 / 24
			hour   = t / 3600 % 24
			minute = t / 60 % 60
			second = t % 3600
		)

		switch {
		case day > 0:
			temp = append(temp, "0", "0", "0")
			temp = append(temp, strings.Join([]string{ASTERISK, strconv.Itoa(day)}, SLASH))
		case hour > 0:
			temp = append(temp, "0", "0")
			temp = append(temp, strings.Join([]string{ASTERISK, strconv.Itoa(hour)}, SLASH))
		case minute > 0:
			temp = append(temp, "0")
			temp = append(temp, strings.Join([]string{ASTERISK, strconv.Itoa(minute)}, SLASH))
		case second > 0:
			temp = append(temp, strings.Join([]string{ASTERISK, strconv.Itoa(second)}, SLASH))
		}
		// fix
		for i := 6 - len(temp); i > 0; i-- {
			temp = append(temp, ASTERISK)
		}
		return strings.Join(temp, BLANK)
	case 5:
		temp = append([]string{"0"}, temp...)
		fallthrough
	case 6:
		temp = append(temp, strs...)
		// convert month
		temp[4] = convertMonth(temp[4])
		// convert weekday
		temp[5] = convertWeek(temp[5])
		return strings.Join(temp, BLANK)
	}
	return ""
}

func convertMonth(m string) string {
	m = strings.ToLower(m)
	for name, month := range months.names {
		m = strings.Replace(m, name, strconv.Itoa(int(month)), -1)
	}
	return m
}

func convertWeek(w string) string {
	w = strings.ToLower(w)
	for name, day := range weeks.names {
		w = strings.Replace(w, name, strconv.Itoa(int(day)), -1)
	}
	return w
}
