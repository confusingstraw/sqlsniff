package main

// little endian
// https://dev.mysql.com/doc/internals/en/integer.html#packet-Protocol::FixedLengthInteger
func ToFixedLengthInt3(value []byte) int {
	return int(uint(value[0]) | uint(value[1])<<8 | uint(value[2])<<16)
}
