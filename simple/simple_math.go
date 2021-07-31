package simple

import "math"

func MathMin(inx int, nums ...int) (min int) {
	min = inx

	for _, num := range nums {
		if min > num {
			min = num
		}
	}

	return
}

func MathMinUint32(inx uint32, nums ...uint32) (min uint32) {
	min = inx

	for _, num := range nums {
		if min > num {
			min = num
		}
	}

	return
}

func MathMax(inx int, nums ...int) (max int) {
	max = inx

	for _, num := range nums {
		if max < num {
			max = num
		}
	}

	return
}

func Round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}
