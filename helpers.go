package main

// HandleError : Handles an `error != nil`
func HandleError(err error) bool {
	if err != nil {
		panic(err)
	}

	return true
}

// Unique : reduces the redundancies in a slice
func Unique(arr []string) []string {
	uniq := make([]string, 0)
	check := make(map[string]bool)

	for _, item := range arr {
		if v := check[item]; !v {
			check[item] = true
			uniq = append(uniq, item)
		}
	}

	return uniq
}

// Count : counts occurences of a word
func Count(arr []string, search string) int {
	c := 0

	for _, x := range arr {
		if x == search {
			c++
		}
	}

	return c
}

func Pow(a, n int) int {
	out := 1

	for i := 1; i < n; i++ {
		out *= a
	}

	return out
}

// Abs : returns the absolute value of an int
func Abs(a int) int {
	if a >= 0 {
		return a
	}
	return -a
}
