package dbx

func IntSlice2ISlice(set []int) []interface{} {
	res := make([]interface{}, len(set))
	for i := 0; i < len(set); i++ {
		res[i] = set[i]
	}

	return res
}

func StrSliceToISlice(set []string) []interface{} {
	res := make([]interface{}, len(set))
	for i := 0; i < len(set); i++ {
		res[i] = set[i]
	}

	return res
}
