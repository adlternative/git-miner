package utils

import "encoding/binary"

func Htonl(val uint32) uint32 {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, val)
	return binary.LittleEndian.Uint32(buf)
}

func Ntohl(networkVal uint32) uint32 {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, networkVal)
	return binary.BigEndian.Uint32(buf)
}
