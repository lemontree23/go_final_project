package scheduler

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const TimeFormat = "20060102"

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("repeat is required")
	}

	start_date, err := time.Parse(TimeFormat, date)
	if err != nil {
		return "", fmt.Errorf("failed to parse task date: %w", err)
	}

	parts := strings.Split(repeat, " ")

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid repeat daily")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days <= 0 || days > 400 {
			return "", fmt.Errorf("invalid repeat daily interval")
		}
		next_date := start_date.AddDate(0, 0, days)
		for !next_date.After(now) {
			next_date = next_date.AddDate(0, 0, days)
		}
		return next_date.Format(TimeFormat), nil
	case "y":
		next_date := start_date.AddDate(1, 0, 0)
		for !next_date.After(now) {
			next_date = next_date.AddDate(1, 0, 0)
		}
		return next_date.Format("20060102"), nil

	default:
		return "", errors.New("unsupported or invalid repeat rule")
	}
}
