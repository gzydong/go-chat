package helper

func InStr(find string, arr []string) bool {
	for _, s := range arr {
		if s == find {
			return true
		}
	}
	return false
}

func InInt(find int, arr []int) bool {
	for _, s := range arr {
		if s == find {
			return true
		}
	}
	return false
}
