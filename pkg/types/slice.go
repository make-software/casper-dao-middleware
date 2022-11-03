package types

func ContainsSlice[T comparable](collection, slice []T) bool {
	if len(slice) > len(collection) {
		return false
	}
	for _, elem := range slice {
		if !Contains(collection, elem) {
			return false
		}
	}
	return true
}

func Contains[T comparable](collection []T, element T) bool {
	for _, value := range collection {
		if value == element {
			return true
		}
	}
	return false
}
