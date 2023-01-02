package databox

func BoolToBytes(num bool) [1]byte {
	bytes := [1]byte{}
	if num {
		bytes[0] = 1
	} else {
		bytes[0] = 0
	}
	return bytes
}

func BytesToBool(bytes [1]byte) bool {
	if bytes[0] == 1 {
		return true
	}
	return false
}
