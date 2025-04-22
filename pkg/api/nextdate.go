package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateLayout = "20060102"

func isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}

func afterNow(date, now time.Time) bool {
	return date.After(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()))
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if dstart == "" {
		return "", fmt.Errorf("empty start date")
	}
	if repeat == "" {
		return "", fmt.Errorf("empty repeat rule")
	}

	startDate, err := time.Parse(dateLayout, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid date format")
	}

	parts := strings.Split(repeat, " ")
	if len(parts) < 1 {
		return "", fmt.Errorf("invalid repeat format")
	}

	switch parts[0] {
	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf("invalid yearly repeat format")
		}
		day := startDate.Day()
		month := startDate.Month()

		for {
			startDate = startDate.AddDate(1, 0, 0)
			// если 29 февраля — проверим, високосный ли год
			if month == time.February && day == 29 {
				// Если год невисокосный, корректируем дату на 1 марта.
				if !isLeapYear(startDate.Year()) {
					startDate = time.Date(startDate.Year(), time.March, 1, 0, 0, 0, 0, startDate.Location())
				} else {
					startDate = time.Date(startDate.Year(), time.February, 29, 0, 0, 0, 0, startDate.Location())
				}
			} else {
				startDate = time.Date(startDate.Year(), month, day, 0, 0, 0, 0, startDate.Location())
			}

			if afterNow(startDate, now) {
				break
			}
		}
		return startDate.Format(dateLayout), nil

	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid daily repeat format")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days <= 0 {
			return "", fmt.Errorf("invalid day count")
		}
		// Если days слишком велико (например, больше 365), это может быть ошибкой
		if days > 365 {
			return "", fmt.Errorf("too large day count for repeat")
		}
		for {
			startDate = startDate.AddDate(0, 0, days)
			if afterNow(startDate, now) {
				break
			}
		}
		return startDate.Format(dateLayout), nil

	case "m":
		if len(parts) == 2 {
			days := strings.Split(parts[1], ",")
			return nextMonthlyByDay(startDate, now, days)
		} else if len(parts) == 3 {
			days := strings.Split(parts[1], ",")
			months := strings.Split(parts[2], ",")
			return nextMonthlyByDayAndMonth(startDate, now, days, months)
		}
		return "", fmt.Errorf("invalid monthly repeat format")

	case "w":
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid weekly repeat format")
		}
		days := strings.Split(parts[1], ",")
		return nextWeekly(startDate, now, days)

	default:
		return "", fmt.Errorf("unknown repeat prefix")
	}
}

func nextMonthlyByDay(startDate, now time.Time, days []string) (string, error) {
	var resDate time.Time
	for _, d := range days {
		day, err := strconv.Atoi(d)
		if err != nil || day < -2 || day == 0 || day > 31 {
			return "", fmt.Errorf("invalid day in repeat")
		}
		date := startDate
		for {
			if day > 0 {
				maxDay := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, date.Location()).Day()
				if day > maxDay {
					date = date.AddDate(0, 1, 0)
					continue
				}
				date = time.Date(date.Year(), date.Month(), day, 0, 0, 0, 0, date.Location())
			} else {
				date = time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, date.Location())
				date = date.AddDate(0, 0, day+1)
			}
			if afterNow(date, now) {
				break
			}
			date = date.AddDate(0, 1, 0)
		}
		if resDate.IsZero() || date.Before(resDate) {
			resDate = date
		}
	}
	return resDate.Format(dateLayout), nil
}

func nextMonthlyByDayAndMonth(startDate, now time.Time, days, months []string) (string, error) {
	var resDate time.Time
	for _, d := range days {
		day, err := strconv.Atoi(d)
		if err != nil || day < -2 || day == 0 || day > 31 {
			return "", fmt.Errorf("invalid day in repeat")
		}
		for _, m := range months {
			month, err := strconv.Atoi(m)
			if err != nil || month < 1 || month > 12 {
				return "", fmt.Errorf("invalid month in repeat")
			}
			date := startDate
			for {
				if day > 0 {
					date = time.Date(date.Year(), time.Month(month), day, 0, 0, 0, 0, date.Location())
				} else {
					date = time.Date(date.Year(), time.Month(month)+1, 0, 0, 0, 0, 0, date.Location())
					date = date.AddDate(0, 0, day+1)
				}
				if afterNow(date, now) {
					break
				}
				date = date.AddDate(1, 0, 0)
			}
			if resDate.IsZero() || date.Before(resDate) {
				resDate = date
			}
		}
	}
	return resDate.Format(dateLayout), nil
}

func nextWeekly(startDate, now time.Time, days []string) (string, error) {
	var resDate time.Time
	for _, d := range days {
		day, err := strconv.Atoi(d)
		if err != nil || day < 1 || day > 7 {
			return "", fmt.Errorf("invalid weekday in repeat")
		}
		date := startDate
		offset := (7 + day - int(date.Weekday())) % 7
		if offset == 0 {
			offset = 7
		}
		date = date.AddDate(0, 0, offset)
		for !afterNow(date, now) {
			date = date.AddDate(0, 0, 7)
		}
		if resDate.IsZero() || date.Before(resDate) {
			resDate = date
		}
	}
	return resDate.Format(dateLayout), nil
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(dateLayout, nowStr)
		if err != nil {
			http.Error(w, "invalid now date", http.StatusBadRequest)
			return
		}
	}

	next, err := NextDate(now, dateStr, repeatStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, next)
}
