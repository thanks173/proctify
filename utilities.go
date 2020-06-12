package proctify

func subtraction(s1 []int, s2 []int) []int {
	subset := make([]int, 0, 100)
	hash := make(map[int]bool)
	for _, n := range s2 {
		hash[n] = true
	}
	for _, n := range s1 {
		if !hash[n] {
			subset = append(subset, n)
		}
	}
	return removeDuplicates(subset)
}

func intersection(s1 []int, s2 []int) []int {
	subset := make([]int, 0, 100)
	hash := make(map[int]bool)
	for _, n := range s1 {
		hash[n] = true
	}
	for _, n := range s2 {
		if hash[n] {
			subset = append(subset, n)
		}
	}
	return removeDuplicates(subset)
}

func removeDuplicates(elements []int) (result []int) {
	encountered := make(map[int]bool)
	for _, element := range elements {
		if !encountered[element] {
			result = append(result, element)
			encountered[element] = true
		}
	}
	return
}
