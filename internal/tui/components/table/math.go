package table

func clamp(v, low, high int) int {
	return min(max(v, low), high)
}
