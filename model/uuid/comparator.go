package uuid

func (u UUID) Equal(other UUID) bool {
	for i := 0; i < uuidLen; i++ {
		if u[i] != other[i] {
			return false
		}
	}
	return true
}
