package util

type HashSet[K comparable, V any] struct {
	data map[K]V
}

func (i *HashSet[K, V]) TryAdd(key K, val V) bool {
	_, found := i.data[key]

	if found {
		return false
	}

	i.data[key] = val

	return true
}
func (i *HashSet[K, V]) TryUpdate(key K, val V) bool {
	_, found := i.data[key]
	if !found {
		return false
	}

	i.data[key] = val

	return true
}
func (i *HashSet[K, V]) TryGet(key K, val *V) bool {
	data, found := i.data[key]
	if !found {
		return false
	}

	*val = data

	return true
}

// func (i *HashSet[K, V]) Remove() bool {

// }
