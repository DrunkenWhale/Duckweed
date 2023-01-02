package databox

func StringToBytes(num string) [128]byte {
	bytes := [128]byte{}
	copy(bytes[:], num)
	return bytes
}

func BytesToString(bytes [128]byte) string {
	return string(bytes[:])
}
