package utils

type Number interface {
	int64 | float64 | uint64 | uint32
}

func PercentOf[T Number](part T, total T) float64 {
	return (float64(part) * float64(100)) / float64(total)
}
