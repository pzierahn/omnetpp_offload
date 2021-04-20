package simple

func CombineBuffers(buffers ...[]byte) (grown []byte) {

	total := 0
	for _, buffer := range buffers {
		total += len(buffer)
	}

	grown = make([]byte, total)

	start := 0
	for _, buffer := range buffers {
		copy(grown[start:], buffer)
		start += len(buffer)
	}

	return
}

func BytesToShorts(bytes []byte) (shorts []int16) {

	shorts = make([]int16, len(bytes)/2)

	for inx := 0; inx < len(shorts); inx++ {
		shorts[inx] = int16(bytes[(inx<<1)+0])<<0 |
			int16(bytes[(inx<<1)+1])<<8
	}

	return
}

func ShortsToBytes(shorts []int16) (bytes []byte) {

	bytes = make([]byte, len(shorts)*2)

	for inx := 0; inx < len(shorts); inx++ {
		bytes[(inx<<1)+0] = byte(shorts[inx] >> 0)
		bytes[(inx<<1)+1] = byte(shorts[inx] >> 8)
	}

	return
}
