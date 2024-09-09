package timewheel

/*
   @Author: orbit-w
   @File: utils
   @2024 9月 周一 07:40
*/

func truncate(x, m int64) int64 {
	if m <= 0 {
		return x
	}
	return x - x%m
}
