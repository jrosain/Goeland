package global

type Map[T Basic[T], U Basic[U]] struct {
	keys   List[T]
	values List[U]
}

func (cm *Map[T, U]) Get(otherKey T) U {
	for i, key := range cm.keys {
		if key.Equals(otherKey) {
			return cm.values[i]
		}
	}
	var zero U
	return zero
}

func (cm *Map[T, U]) GetExists(otherKey T) (U, bool) {
	for i, key := range cm.keys {
		if key.Equals(otherKey) {
			return cm.values[i], true
		}
	}
	var zero U
	return zero, false
}

func (cm *Map[T, U]) Exists(otherKey T) bool {
	_, result := cm.GetExists(otherKey)
	return result
}

func (cm *Map[T, U]) Set(otherKey T, otherValue U) {
	for i, key := range cm.keys {
		if key.Equals(otherKey) {
			cm.values[i] = otherValue
			return
		}
	}

	cm.keys = *cm.keys.Append(otherKey)
	cm.values = *cm.values.Append(otherValue)
}

func (cm *Map[T, U]) Length() int {
	return len(cm.keys)
}

func (cm *Map[T, U]) Clear() {
	cm.keys = List[T]{}
	cm.values = List[U]{}
}

func (cm *Map[T, U]) InsertInto(other *Map[T, U]) {
	for i := range cm.keys {
		other.Set(cm.keys[i], cm.values[i])
	}
}

func (cm *Map[T, U]) Keys() List[T] {
	return cm.keys
}

func (cm *Map[T, U]) Values() List[U] {
	return cm.values
}

func (cm *Map[T, U]) ToString() string {
	str := ""

	for i, key := range cm.keys {
		str += key.ToString() + " -> " + cm.values[i].ToString() + "\n"
	}

	return str[:len(str)-1]
}

func (cm *Map[T, U]) Equals(other any) bool {
	if typed, ok := other.(*Map[T, U]); ok {
		if len(typed.keys) == len(cm.keys) {
			for i, key := range cm.keys {
				value, exists := typed.GetExists(key)

				if !exists || !cm.values[i].Equals(value) {
					return false
				}
			}

			return true
		}
	}

	return false
}

func (cm *Map[T, U]) Copy() *Map[T, U] {
	newCm := new(Map[T, U])
	newCm.keys = *newCm.keys.Append(cm.keys...)
	newCm.values = *newCm.values.Append(cm.values...)
	return newCm
}
