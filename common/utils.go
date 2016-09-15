package common

import (
	"fmt"
	"net/mail"
)

//Checks whether string has special characters
func HasSpecials(text string) bool {
	for i := 0; i < len(text); i++ {
		switch c := text[i]; c {
		case '(', ')', '<', '>', '[', ']', ':', ';', '@', '\\', ',', '.', '"':
			return true
		}
	}

	return false
}

//Check whether a string is present in string array
func SearchString(strArray []string, value string) bool {
	for _, val := range strArray {
		if val == value {
			return true
		}
	}
	return false
}

//Add an address to an address array
func AddAddress(addrArray []string, addr string) []string {
	if SearchString(addrArray, addr) {
		return addrArray
	}

	return append(addrArray, addr)
}

func ParseAddress(addr string) (string, error) {
	parsedAddr, err := mail.ParseAddress(addr)
	if err != nil {
		return "", fmt.Errorf("m-mail: invalid address %q: %v", addr, err)
	}
	return parsedAddr.Address, nil
}
