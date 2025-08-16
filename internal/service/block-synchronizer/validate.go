package synchronizer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/yanmoyy/tbi/internal/indexer"
)

func validateTransferEvent(event indexer.GnoEvent) error {
	if event.Type != "Transfer" {
		return fmt.Errorf("invalid event type: %s", event.Type)
	}

	var from, to, value string
	attrMap := make(map[string]string)
	for _, attr := range event.Attrs {
		attrMap[attr.Key] = attr.Value
	}

	from, fromExists := attrMap["from"]
	to, toExists := attrMap["to"]
	value, valueExists := attrMap["value"]
	if !fromExists || !toExists || !valueExists {
		return fmt.Errorf("missing attributes")
	}

	if !isNumeric(value) {
		return fmt.Errorf("value should be numeric: %s", value)
	}

	type attrs struct {
		from  string
		to    string
		value string
	}
	validFuncs := map[string]func(attrs) bool{
		"Mint":     func(a attrs) bool { return a.from == "" && isBech32(a.to) },
		"Burn":     func(a attrs) bool { return isBech32(a.from) && a.to == "" },
		"Transfer": func(a attrs) bool { return isBech32(a.from) && isBech32(a.to) },
	}

	f, ok := validFuncs[event.Func]
	if !ok {
		return fmt.Errorf("invalid function: %s", event.Func)
	}
	if !f(attrs{from, to, value}) {
		return fmt.Errorf("invalid attributes")
	}
	return nil
}

func isNumeric(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}

func isBech32(s string) bool {
	if len(s) < 8 || len(s) > 90 {
		return false
	}

	sepIndex := strings.LastIndex(s, "1")
	if sepIndex < 1 || sepIndex == len(s)-1 {
		return false
	}

	hrp := s[:sepIndex] // human readable part
	data := s[sepIndex+1:]

	if len(hrp) < 1 || len(hrp) > 83 {
		return false
	}
	for _, c := range hrp {
		if (c < 'a' || c > 'z') && (c < '0' || c > '9') {
			return false
		}
	}

	if len(data) < 6 {
		return false
	}
	for _, c := range data {
		if !strings.ContainsRune("qpzry9x8gf2tvdw0s3jn54khce6mua7l", c) {
			return false
		}
	}

	return hrp == "g"
}
