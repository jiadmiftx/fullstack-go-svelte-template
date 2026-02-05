package common_utils

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

func ParseDuration(durationStr string) (time.Duration, error) {
	re := regexp.MustCompile(`(\d+)([dhm])`)
	matches := re.FindAllStringSubmatch(durationStr, -1)

	if len(matches) == 0 {
		return 0, fmt.Errorf("invalid duration format: %s (expected format: 1h, 30m, 1d, etc.)", durationStr)
	}

	var totalDuration time.Duration
	for _, match := range matches {
		value, err := strconv.Atoi(match[1])
		if err != nil {
			return 0, fmt.Errorf("invalid number in duration: %s", match[1])
		}

		unit := match[2]
		switch unit {
		case "m":
			totalDuration += time.Duration(value) * time.Minute
		case "h":
			totalDuration += time.Duration(value) * time.Hour
		case "d":
			totalDuration += time.Duration(value) * 24 * time.Hour
		default:
			return 0, fmt.Errorf("invalid duration unit: %s (supported: m, h, d)", unit)
		}
	}

	return totalDuration, nil
}

func FormatDuration(d time.Duration) string {
	if d >= 24*time.Hour {
		days := int(d / (24 * time.Hour))
		remainder := d % (24 * time.Hour)
		if remainder == 0 {
			return fmt.Sprintf("%dd", days)
		}
		hours := int(remainder / time.Hour)
		return fmt.Sprintf("%dd%dh", days, hours)
	}
	if d >= time.Hour {
		hours := int(d / time.Hour)
		remainder := d % time.Hour
		if remainder == 0 {
			return fmt.Sprintf("%dh", hours)
		}
		minutes := int(remainder / time.Minute)
		return fmt.Sprintf("%dh%dm", hours, minutes)
	}
	minutes := int(d / time.Minute)
	return fmt.Sprintf("%dm", minutes)
}

func DefaultDuration(diff time.Duration) time.Duration {
	switch {
	case diff <= time.Hour:
		return time.Minute
	case diff <= 24*time.Hour:
		return time.Minute * 10
	case diff <= time.Hour*24*7:
		return time.Hour
	case diff > time.Hour*24*27 && diff <= time.Hour*24*365:
		return time.Hour * 24
	case diff <= time.Hour*24*400*10:
		return time.Hour * 24 * 30
	default:
		return time.Minute
	}
}
