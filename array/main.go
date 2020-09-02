package array

/**
 * 数组去重 去空
 */
func RemoveDuplicatesAndEmpty(a []string) (ret []string) {
	a_len := len(a)
	for i := 0; i < a_len; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}

/**
 * 数组去重 去空
 */
func RemoveEmpty(a []string) (ret []string) {
	a_len := len(a)
	for i := 0; i < a_len; i++ {
		if len(a[i]) == 0 || a[i] == " " {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}
