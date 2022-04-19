package sliceutil

func InStr(find string, arr []string) bool {
	for _, value := range arr {
		if value == find {
			return true
		}
	}
	return false
}

func InInt(find int, arr []int) bool {
	for _, value := range arr {
		if value == find {
			return true
		}
	}
	return false
}

func InInt64(find int64, arr []int64) bool {
	for _, value := range arr {
		if value == find {
			return true
		}
	}
	return false
}
