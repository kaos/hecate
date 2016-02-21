package main

func Max (a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min (a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Abs (v int) int {
	if v < 0 {
		return -v
	}
	return v
}
