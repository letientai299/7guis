package task

import (
	"fmt"
	"math"

	"github.com/maja42/goval"
)

var cellFuncs = map[string]goval.ExpressionFunction{
	"SUM": fnSUM,
	"sum": fnSUM,

	"AVG": fnAVG,
	"avg": fnAVG,

	"MAX": fnMAX,
	"max": fnMAX,

	"MIN": fnMAX,
	"min": fnMIN,
}

func fnMIN(args ...interface{}) (interface{}, error) {
	s := math.MaxFloat64
	for _, a := range args {
		switch x := a.(type) {
		case float64:
			if s > x {
				s = x
			}
		case int:
			if s > float64(x) {
				s = float64(x)
			}
		default:
			return 0, fmt.Errorf("%v is NaN", a)
		}
	}
	return s, nil
}

func fnMAX(args ...interface{}) (interface{}, error) {
	s := -math.MaxFloat64

	for _, a := range args {
		switch x := a.(type) {
		case float64:
			if s < x {
				s = x
			}
		case int:
			if s < float64(x) {
				s = float64(x)
			}
		default:
			return 0, fmt.Errorf("%v is NaN", a)
		}
	}
	return s, nil
}

func fnAVG(args ...interface{}) (interface{}, error) {
	if len(args) == 0 {
		return 0, nil
	}

	s, err := fnSUM(args...)
	if err != nil {
		return nil, err
	}

	return s.(float64) / float64(len(args)), nil
}

func fnSUM(args ...interface{}) (interface{}, error) {
	s := float64(0)
	for _, a := range args {
		switch x := a.(type) {
		case float64:
			s += x
		case int:
			s += float64(x)
		default:
			return float64(0), fmt.Errorf("%v is NaN", a)
		}
	}
	return s, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
