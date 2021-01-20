package util

import (
	"math"
	"math/rand"
	"net"
)

func RandUDPAddr() net.UDPAddr {
	ip := []byte{byte(rand.Int()), byte(rand.Int()), byte(rand.Int()), byte(rand.Int())}

	return net.UDPAddr{
		IP:   ip,
		Port: rand.Int() % math.MaxUint16,
		Zone: RandString(rand.Intn(127)),
	}
}

func RandString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(32 + rand.Intn(94))
	}

	return string(bytes)
}

func RandBool() bool {
	return rand.Intn(2) == 1
}

func RandIntSlice(maxLength int) []int {
	ints := make([]int, maxLength)
	for i := 0; i < maxLength; i++ {
		ints[i] = rand.Int();	
	}

	return ints
}