package common

import (
	"fmt"
	"net/mail"
	"time"
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
func AddStrToUniqueList(addrArray []string, addr string) []string {
	if SearchString(addrArray, addr) {
		return addrArray
	}

	return append(addrArray, addr)
}

//Check for a valid email address
func ParseAddress(addr string) (string, error) {
	parsedAddr, err := mail.ParseAddress(addr)
	if err != nil {
		return "", fmt.Errorf("m-mail: invalid address %q: %v", addr, err)
	}

	return parsedAddr.Address, nil
}

//Return formated address for host and port
func HostPortAddr(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}

// FormatDate formats a date as a valid RFC 5322 date.
func FormatDate(date time.Time) string {
	return date.Format(time.RFC1123Z)
}
