package converter

import (
	"encoding/binary"
	"math"
	"strconv"
)

func Int64ToBytes(value int64) [8]byte {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], uint64(value))
	return buf
}

func Float64ToBytes(value float64) [8]byte {
	bits := math.Float64bits(value)
	var bytes [8]byte
	binary.LittleEndian.PutUint64(bytes[:], bits)
	return bytes
}
func FloatToString(value float64) string {
	return strconv.FormatFloat(value, 'g', 1, 64)
}
func Int64ToString(value int64) string {
	return strconv.FormatInt(value, 64)
}
