package common

func IsInArray(arr []string, item string) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return true
}
