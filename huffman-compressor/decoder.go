package main

import "strings"

func uint32ToBitString(code uint32, length uint8) string {
	var result strings.Builder
	for i := range length {
		if code&(1<<(31-i)) != 0 {
			result.WriteByte('1')
		} else {
			result.WriteByte('0')
		}
	}
	return result.String()
}
