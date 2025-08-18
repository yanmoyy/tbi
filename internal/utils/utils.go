package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func IsNumeric(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}

func IsBech32(s string) bool {
	if len(s) < 8 || len(s) > 90 {
		return false
	}

	sepIndex := strings.Index(s, "1")
	if sepIndex < 1 || sepIndex == len(s)-1 {
		return false
	}

	hrp := s[:sepIndex] // human readable part
	data := s[sepIndex+1:]

	if hrp != "g" { // check hrp
		return false
	}

	if len(data) < 6 {
		return false
	}

	// TODO: need to check this is right.
	for _, c := range data {
		// bech32 chars: https://www.cs.utexas.edu/~moore/acl2/manuals/current/manual/index-seo.php?xkey=BITCOIN_____A2BECH32-CHAR-VALS_A2&path=4368/2215/6197/6198
		if !strings.ContainsRune("qpzry9x8gf2tvdw0s3jn54khce6mua7l", c) {
			return false
		}
	}

	return true
}

func ExtractGRC20TokenPath(token string) (string, error) {
	parts := strings.Split(token, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid token format: %s, expected {tokenPath}:{symbol}", token)
	}
	tokenPath := parts[0]
	symbol := parts[1]
	if tokenPath == "" || symbol == "" {
		return "", fmt.Errorf("tokenPath or symbol empty in token: %s", token)
	}
	if !regexp.MustCompile(`^gno\.land/r/[\w-]+/[\w-]+$`).MatchString(tokenPath) {
		return "", fmt.Errorf("invalid tokenPath format: %s", tokenPath)
	}
	return tokenPath, nil
}
