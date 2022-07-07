package function

func CheckInNumberSlice[T uint64 | int32](a T, s []T) bool {
	for _, val := range s {
		if a == val {
			return true
		}
	}
	return false
}

//DelEleInSlice 只支持 元素不重复的情景
func DelEleInSlice[T uint64 | int32](a T, old []T) (new []T) {
	for i, val := range old {
		if a == val {
			new = append(old[:i], old[i+1:]...)
			return
		}
	}
	return old
}
