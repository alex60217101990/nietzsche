package helpers

import (
	"strconv"
	"time"

	"github.com/alex60217101990/nietzsche/external/configs"
)

func TimeToTimePtr(t time.Time) *time.Time {
	return &t
}

func TimePtrToTime(t *time.Time) (emptyTime time.Time) {
	if t != nil {
		return *t
	}
	return emptyTime
}

func TimeoutSecond(seconds interface{}) time.Duration {
	switch s := seconds.(type) {
	case *int:
		return time.Duration(*s) * time.Second
	case int:
		return time.Duration(s) * time.Second
	case *string:
		if i32, err := strconv.Atoi(*s); err != nil {
			return time.Duration(configs.Conf.Timeouts.DefaultTimeout) * time.Second
		} else {
			return time.Duration(i32) * time.Second
		}
	case string:
		if i32, err := strconv.Atoi(s); err != nil {
			return time.Duration(configs.Conf.Timeouts.DefaultTimeout) * time.Second
		} else {
			return time.Duration(i32) * time.Second
		}
	default:
		return time.Duration(configs.Conf.Timeouts.DefaultTimeout) * time.Second
	}
}
