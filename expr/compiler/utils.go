package compiler

import "encoding/binary"

func encode(input uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, input)

	return b
}
