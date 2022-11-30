package databox

func IntToBytes(num int64) [8]byte {
	ui := uint64(num)
	bytes := [8]byte{}
	for i := 0; i < 8; i++ {
		// 超过的部分会在强转时被抛弃
		bytes[i] = uint8(ui >> (i * 8))
	}
	return bytes
}

func BytesToInt(bytes [8]byte) int64 {
	ui := uint64(0)
	for i := 0; i < 8; i++ {
		ui <<= 8
		ui += uint64(bytes[7-i])
	}
	return int64(ui)
}
