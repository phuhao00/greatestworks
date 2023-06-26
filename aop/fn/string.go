package fn

import (
	"strconv"
	"strings"
)

func SplitStringToInt32Slice(src string, sep string) []int32 {
	strSlice := strings.Split(src, sep)
	int32Slice := make([]int32, 0, len(strSlice))
	for _, item := range strSlice {
		value, err := strconv.ParseInt(item, 10, 32)
		if err != nil {
			continue
		}
		int32Slice = append(int32Slice, int32(value))
	}

	return int32Slice
}

func SplitStringToUint32Slice(src string, sep string) []uint32 {
	strSlice := strings.Split(src, sep)
	var uint32Slice []uint32
	for _, item := range strSlice {
		value, err := strconv.ParseUint(item, 10, 32)
		if err != nil {
			continue
		}
		uint32Slice = append(uint32Slice, uint32(value))
	}

	return uint32Slice
}
