package synchronizer

import (
	"fmt"
	"strconv"

	"github.com/yanmoyy/tbi/internal/indexer"
	"github.com/yanmoyy/tbi/internal/models"
	"github.com/yanmoyy/tbi/internal/utils"
)

// validate and convert gno event to event processor event
func processEvent(event indexer.GnoEvent) (models.TransferEvent, error) {
	if len(event.Attrs) != 4 {
		return models.TransferEvent{}, fmt.Errorf("len(attrs) != 4")
	}

	attrMap := make(map[string]string)
	for _, attr := range event.Attrs {
		attrMap[attr.Key] = attr.Value
	}

	for _, key := range []string{"token", "from", "to", "value"} {
		if _, ok := attrMap[key]; !ok {
			return models.TransferEvent{}, fmt.Errorf("missing key: %s", key)
		}
	}

	valueInt, err := strconv.ParseInt(attrMap["value"], 10, 64)
	if err != nil {
		return models.TransferEvent{}, fmt.Errorf("value couldn't be parsed: %w", err)
	}

	tokenPath, err := utils.ExtractGRC20TokenPath(attrMap["token"])
	if err != nil {
		return models.TransferEvent{}, fmt.Errorf("token couldn't be parsed: %w", err)
	}

	type attrs struct {
		from  string
		to    string
		value string
	}

	validFuncs := map[string]func(attrs) bool{
		"Mint":     func(a attrs) bool { return a.from == "" && utils.IsBech32(a.to) },
		"Burn":     func(a attrs) bool { return utils.IsBech32(a.from) && a.to == "" },
		"Transfer": func(a attrs) bool { return utils.IsBech32(a.from) && utils.IsBech32(a.to) },
	}

	f, ok := validFuncs[event.Func]
	if !ok {
		return models.TransferEvent{}, fmt.Errorf("invalid function: %s", event.Func)
	}
	if !f(attrs{attrMap["from"], attrMap["to"], attrMap["value"]}) {
		return models.TransferEvent{}, fmt.Errorf("invalid attributes for function: %s", event.Func)
	}
	return models.TransferEvent{
		Func:      models.EventFunc(event.Func),
		TokenPath: tokenPath,
		From:      attrMap["from"],
		To:        attrMap["to"],
		Value:     valueInt,
	}, nil
}
