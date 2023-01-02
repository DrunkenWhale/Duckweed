package databox

func FloatToBytes(num float64) [8]byte {
	ui := uint64(num)
	bytes := [8]byte{}
	for i := 0; i < 8; i++ {
		bytes[i] = uint8(ui >> (i * 8))
	}
	return bytes
}

func BytesToFloat(bytes [8]byte) float64 {
	ui := uint64(0)
	for i := 0; i < 8; i++ {
		ui <<= 8
		ui += uint64(bytes[7-i])
	}
	return float64(ui)
}
