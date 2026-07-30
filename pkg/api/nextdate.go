package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repeat rule is empty")
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", err
	}

	parts := strings.Split(repeat, " ")

	switch parts[0] {
	case "d":
		if len(parts) < 2 {
			return "", errors.New("interval not specified for rule d")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", err
		}
		if days <= 0 || days > 400 {
			return "", errors.New("invalid interval for rule d")
		}
		for {
			date = date.AddDate(0, 0, days)
			if afterNow(date, now) {
				break
			}
		}

	case "y":
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}

	case "w":
		if len(parts) < 2 {
			return "", errors.New("weekdays not specified for rule w")
		}
		var weekdays [8]bool
		for _, ds := range strings.Split(parts[1], ",") {
			d, err := strconv.Atoi(ds)
			if err != nil || d < 1 || d > 7 {
				return "", errors.New("invalid weekday value")
			}
			weekdays[d] = true
		}
		for {
			date = date.AddDate(0, 0, 1)
			wd := int(date.Weekday())
			if wd == 0 {
				wd = 7
			}
			if weekdays[wd] && afterNow(date, now) {
				break
			}
		}

	case "m":
		if len(parts) < 2 {
			return "", errors.New("days of month not specified for rule m")
		}
		var days [32]bool
		var negDays []int
		for _, ds := range strings.Split(parts[1], ",") {
			d, err := strconv.Atoi(ds)
			if err != nil || d == 0 || d < -2 || d > 31 {
				return "", errors.New("invalid day of month")
			}
			if d < 0 {
				negDays = append(negDays, d)
			} else {
				days[d] = true
			}
		}

		var months [13]bool
		hasMonths := false
		if len(parts) >= 3 {
			hasMonths = true
			for _, ms := range strings.Split(parts[2], ",") {
				m, err := strconv.Atoi(ms)
				if err != nil || m < 1 || m > 12 {
					return "", errors.New("invalid month")
				}
				months[m] = true
			}
		}

		for {
			date = date.AddDate(0, 0, 1)
			if hasMonths && !months[int(date.Month())] {
				continue
			}
			matched := days[date.Day()]
			if !matched && len(negDays) > 0 {
				lastDay := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, date.Location()).Day()
				for _, nd := range negDays {
					if date.Day() == lastDay+nd+1 {
						matched = true
						break
					}
				}
			}
			if matched && afterNow(date, now) {
				break
			}
		}

	default:
		return "", errors.New("unsupported repeat format")
	}

	return date.Format(dateFormat), nil
}

// afterNow сравнивает только даты, без учёта времени суток
func afterNow(date, now time.Time) bool {
	y1, m1, d1 := date.Date()
	y2, m2, d2 := now.Date()
	d1n := time.Date(y1, m1, d1, 0, 0, 0, 0, time.UTC)
	d2n := time.Date(y2, m2, d2, 0, 0, 0, 0, time.UTC)
	return d1n.After(d2n)
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(dateFormat, nowStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	next, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(next))
}
