package databox

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestIntBox(t *testing.T) {
	ins := []int64{
		-(1 << 63),
		-1145141919810,
		1,
		0,
		114514,
		1919810,
		(1 << 63) - 1,
	}
	for i := 0; i < len(ins); i++ {
		n := ins[i]
		bytes := IntToBytes(n)
		num := BytesToInt(bytes)
		assert.Equal(t, n, num)
	}
}
