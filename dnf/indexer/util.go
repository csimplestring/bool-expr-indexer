package indexer

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func search(start, end int, f func(int) bool) int {
	// Define f(-1) == false and f(n) == true.
	// Invariant: f(i-1) == false, f(j) == true.
	if start > end {
		panic("start must be not greater than end")
	}

	i, j := start, end
	for i < j {
		h := int(uint(i-start+j-start)>>1) + start // avoid overflow when computing h
		// i ≤ h < j
		if !f(h) {
			i = h + 1 // preserves f(i-1) == false
		} else {
			j = h // preserves f(j) == true
		}
	}
	// i == j, f(i-1) == false, and f(j) (= f(i)) == true  =>  answer is i.
	return i
}
