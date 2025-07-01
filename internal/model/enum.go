package model

import "time"

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
