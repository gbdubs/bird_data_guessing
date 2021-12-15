package bird_data_guessing

import "strconv"

func atoiOrFail(s string) int {
	i, e := strconv.Atoi(s)
	if e != nil {
		panic(e)
	}
	return i
}

func floatOrFail(s string) float64 {
	f, e := strconv.ParseFloat(s, 64)
	if e != nil {
		panic(e)
	}
	return f
}
