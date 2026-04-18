package main

var (
	// leftForwardCodeTable
	leftFCT = map[int]int{
		0001101: 0,
		0011001: 1,
		0010011: 2,
		0111101: 3,
		0100011: 4,
		0110001: 5,
		0101111: 6,
		0111011: 7,
		0110111: 8,
		0001011: 9,
	}
	// rightForwardCodeTable
	rightFCT = map[int]int{
		1110010: 0,
		1100110: 1,
		1101100: 2,
		1000010: 3,
		1011100: 4,
		1001110: 5,
		1010000: 6,
		1000100: 7,
		1001000: 8,
		1110100: 9,
	}

	// leftBackwardCodeTable
	leftBCT = map[int]int{
		1011000: 0,
		1001100: 1,
		1100100: 2,
		1011110: 3,
		1100010: 4,
		1000110: 5,
		1111010: 6,
		1101110: 7,
		1110110: 8,
		1101000: 9,
	}

	// rightBbackwardCodeTable
	rightBCT = map[int]int{
		0100111: 0,
		0110011: 1,
		0011011: 2,
		0100001: 3,
		0011101: 4,
		0111001: 5,
		0000101: 6,
		0010001: 7,
		0001001: 8,
		0010111: 9,
	}
)

type UPCcode struct {
	leftSideCodes  []int
	rightSideCodes []int
	remainder      int
	checkerL       []int
	checkerR       []int
	checkerC       []int
}

func Construct(code []int, remainder int) UPCcode {
	if len(code) != 95 {
		return
	}

	return UPCcode{
		leftSideCodes:  code[3:46],
		rightSideCodes: code[51:93],
		checkerL:       code[:3],
		checkerR:       code[92:],
		checkerC:       code[46:51],
		remainder:      remainder,
	}
}

func GenereteUPC() UPCcode {
	return UPCcode{}
}

// return true if code secseed check and int = 1 if it's even and 0 if it's odd.
func (code UPCcode) checkParity() (bool, int) {
	return false, -1
}

// return true if code secsessfuly goes through
func (code UPCcode) checkRemainder() bool {
	return false
}
