package alg

func MakeExecs() map[string]int{
	execs := make(map[string]int)
	for _, k := range Labels {
		execs[k] = 0
	}
	return execs
}