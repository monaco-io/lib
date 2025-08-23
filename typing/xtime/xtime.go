package xtime

import (
	"strconv"
	"strings"
	"time"
)

type Uint interface {
	~uint64 | ~uint32 | ~uint16 | ~uint8 | ~uint
}

var hour8 = (time.Second * 60 * 60 * 8).Seconds()

var Shanghai = time.FixedZone("Asia/Shanghai", int(hour8))

func init() {
	SetLocation(Shanghai)
}

func SetLocation(location *time.Location) {
	time.Local = location
}

func AsMillisecond[T Uint](ms T) time.Duration {
	return time.Duration(ms) * time.Millisecond
}

func AsSecond[T Uint](s T) time.Duration {
	return time.Duration(s) * time.Second
}

func AsMinute[T Uint](m T) time.Duration {
	return time.Duration(m) * time.Minute
}

func AsHour[T Uint](h T) time.Duration {
	return time.Duration(h) * time.Hour
}

func DateCN(t time.Time) string {
	return t.In(Shanghai).Format(time.DateOnly)
}

func DateTimeCN(t time.Time) string {
	return t.In(Shanghai).Format(time.DateTime)
}

func ParseLocal(t string, format ...string) time.Time {
	// Helper function to try parsing and convert to local time
	tryParseAndConvert := func(layout, value string) (time.Time, bool) {
		// For formats with timezone info, use time.Parse to respect the timezone
		if hasTimezone(layout) {
			if parsed, err := parseTimeWithTimezone(layout, value); err == nil {
				return parsed.In(time.Local), true
			}
		}
		// For formats without timezone info, use ParseInLocation with Local
		if parsed, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			return parsed, true
		}
		return time.Time{}, false
	}

	// Handle special case for unix timestamp
	if len(format) > 0 {
		for _, f := range format {
			if f == "unix" {
				// Try to parse as unix timestamp
				if unix := parseUnixTimestamp(t); unix > 0 {
					return time.Unix(unix, 0).In(time.Local)
				}
				continue
			}
			if parsed, ok := tryParseAndConvert(f, t); ok {
				return parsed
			}
		}
	}

	switch len(t) {
	case len(time.TimeOnly):
		if parsed, ok := tryParseAndConvert(time.TimeOnly, t); ok {
			return parsed
		}
	case len(time.Layout):
		if parsed, ok := tryParseAndConvert(time.Layout, t); ok {
			return parsed
		}
	case len(time.ANSIC):
		if parsed, ok := tryParseAndConvert(time.ANSIC, t); ok {
			return parsed
		}
	case len(time.UnixDate):
		if parsed, ok := tryParseAndConvert(time.UnixDate, t); ok {
			return parsed
		}
	case len(time.RFC822Z):
		if parsed, ok := tryParseAndConvert(time.RFC822Z, t); ok {
			return parsed
		}
	case len(time.RFC1123):
		if parsed, ok := tryParseAndConvert(time.RFC1123, t); ok {
			return parsed
		}
	case len(time.RFC1123Z):
		if parsed, ok := tryParseAndConvert(time.RFC1123Z, t); ok {
			return parsed
		}
	case len(time.RFC3339Nano):
		if parsed, ok := tryParseAndConvert(time.RFC3339Nano, t); ok {
			return parsed
		}
	case len(time.Kitchen):
		if parsed, ok := tryParseAndConvert(time.Kitchen, t); ok {
			return parsed
		}
	case len(time.Stamp):
		if parsed, ok := tryParseAndConvert(time.Stamp, t); ok {
			return parsed
		}
	case len(time.StampMicro):
		if parsed, ok := tryParseAndConvert(time.StampMicro, t); ok {
			return parsed
		}
	case 10:
		// Handle DateOnly format (2006-01-02) or unix timestamp
		if parsed, ok := tryParseAndConvert(time.DateOnly, t); ok {
			return parsed
		}
		// Try unix timestamp for 10-digit strings
		if unix := parseUnixTimestamp(t); unix > 0 {
			return time.Unix(unix, 0).In(time.Local)
		}
	case 19:
		if parsed, ok := tryParseAndConvert(time.DateTime, t); ok {
			return parsed
		}
		if parsed, ok := tryParseAndConvert(time.RFC822, t); ok {
			return parsed
		}
		if parsed, ok := tryParseAndConvert(time.StampMilli, t); ok {
			return parsed
		}
	case 20:
		if parsed, ok := tryParseAndConvert(time.RFC3339, t); ok {
			return parsed
		}
	case 25:
		if parsed, ok := tryParseAndConvert(time.RFC3339, t); ok {
			return parsed
		}
		if parsed, ok := tryParseAndConvert(time.StampNano, t); ok {
			return parsed
		}
	case 30:
		if parsed, ok := tryParseAndConvert(time.RFC850, t); ok {
			return parsed
		}
		if parsed, ok := tryParseAndConvert(time.RubyDate, t); ok {
			return parsed
		}
	}
	return time.Time{}
}

func DateUTC(t time.Time) string {
	return t.UTC().Format(time.DateOnly)
}

func DateTimeUTC(t time.Time) string {
	return t.UTC().Format(time.DateTime)
}

// hasTimezone checks if the layout contains timezone information
func hasTimezone(layout string) bool {
	// Check for timezone indicators in common formats
	return strings.Contains(layout, "Z") ||
		strings.Contains(layout, "MST") ||
		strings.Contains(layout, "-07") ||
		strings.Contains(layout, "+07") ||
		strings.Contains(layout, "-0700") ||
		strings.Contains(layout, "+0700") ||
		strings.Contains(layout, "-07:00") ||
		strings.Contains(layout, "+07:00") ||
		layout == time.RFC822 ||
		layout == time.RFC822Z ||
		layout == time.RFC850 ||
		layout == time.RFC1123 ||
		layout == time.RFC1123Z ||
		layout == time.RFC3339 ||
		layout == time.RFC3339Nano
}

// parseTimeWithTimezone handles timezone-aware parsing with special cases for common timezone abbreviations
func parseTimeWithTimezone(layout, value string) (time.Time, error) {
	// First try standard parsing
	if parsed, err := time.Parse(layout, value); err == nil {
		// Check if the parsed timezone has a zero offset but the value contains a known timezone abbreviation
		_, offset := parsed.Zone()
		if offset == 0 && strings.Contains(value, "MST") {
			// Handle MST as UTC-7
			// Re-parse assuming the timezone is actually -7 hours
			if layout == time.RFC1123 {
				// For RFC1123, replace MST with a proper offset and re-parse
				correctedValue := strings.Replace(value, "MST", "-0700", 1)
				if corrected, err := time.Parse(time.RFC1123Z, correctedValue); err == nil {
					return corrected, nil
				}
			}
			// Fallback: manually adjust the time by -7 hours for MST
			return parsed.Add(-7 * time.Hour), nil
		}
		return parsed, nil
	} else {
		return time.Time{}, err
	}
}

// parseUnixTimestamp tries to parse a string as unix timestamp
func parseUnixTimestamp(s string) int64 {
	if unix, err := strconv.ParseInt(s, 10, 64); err == nil {
		return unix
	}
	return 0
}
