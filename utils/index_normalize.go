package utils

func ToPositiveIndex(index, len int) int {
	if index < 0 {
		if index+len > 0 {
			return index + len
		}
		return 0
	}
	return index
}

func GetPositiveStartEndIndexes(start, end, len int) (int, int) {
	if end > len-1 {
		end = len - 1
	}
	return ToPositiveIndex(start, len), ToPositiveIndex(end, len)
}
