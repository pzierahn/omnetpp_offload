package simple

func MathMin(inx int, nums ...int) (min int) {
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
