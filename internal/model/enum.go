package model

import "time"

// ExpireOption defines the available expiration durations
//
// 1 = 6 Hours
// 2 = 12 Hours
// 3 = 1 Day
// 4 = 7 Days
type ExpireOption int

const (
	Expire6Hours  ExpireOption = 1
	Expire12Hours ExpireOption = 2
	Expire1Day    ExpireOption = 3
	Expire7Days   ExpireOption = 4
)

func (e ExpireOption) Duration() time.Duration {
	switch e {
	case Expire6Hours:
		return 6 * time.Hour
	case Expire12Hours:
		return 12 * time.Hour
	case Expire1Day:
		return 24 * time.Hour
	case Expire7Days:
		return 7 * 24 * time.Hour
	default:
		return 0
	}
}

func (e ExpireOption) IsValid() bool {
	switch e {
	case Expire6Hours, Expire12Hours, Expire1Day, Expire7Days:
		return true
	default:
		return false
	}
}

func (e ExpireOption) String() string {
	switch e {
	case Expire6Hours:
		return "6 hours"
	case Expire12Hours:
		return "12 hours"
	case Expire1Day:
		return "1 day"
	case Expire7Days:
		return "7 days"
	default:
		return "invalid"
	}
}
