package overtime

import (
	"math"
)

type Rounding struct{}

func NewRounding() *Rounding {
	return &Rounding{}
}

func (r *Rounding) RoundUp(minutes int, roundingMinutes int) int {
	if roundingMinutes <= 0 {
		return minutes
	}
	return int(math.Ceil(float64(minutes)/float64(roundingMinutes))) * roundingMinutes
}

func (r *Rounding) RoundNearest(minutes int, roundingMinutes int) int {
	if roundingMinutes <= 0 {
		return minutes
	}
	return int(math.Round(float64(minutes)/float64(roundingMinutes))) * roundingMinutes
}

func (r *Rounding) RoundDown(minutes int, roundingMinutes int) int {
	if roundingMinutes <= 0 {
		return minutes
	}
	return int(math.Floor(float64(minutes)/float64(roundingMinutes))) * roundingMinutes
}

func (r *Rounding) Apply(minutes int, roundingMinutes int, mode string) int {
	switch mode {
	case "UP":
		return r.RoundUp(minutes, roundingMinutes)
	case "DOWN":
		return r.RoundDown(minutes, roundingMinutes)
	default:
		return r.RoundNearest(minutes, roundingMinutes)
	}
}
