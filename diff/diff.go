package diff

// Of returns the elements present in A but not in B
func Of(a, b []string) []string {
	var result []string
out:
	for _, x := range a {
		for _, y := range b {
			if x == y {
				continue out
			}
		}
		result = append(result, x)
	}
	return result
}
