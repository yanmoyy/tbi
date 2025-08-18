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
	if len(event.Attrs) != 3 {
		return models.TransferEvent{}, fmt.Errorf("invalid attributes")
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
		return models.TransferEvent{}, fmt.Errorf("missing attributes")
	}

	valueInt, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return models.TransferEvent{}, fmt.Errorf("invalid value")
	}

	type attrs struct {
		from  string
		to    string
		value string
	}
	validFuncs := map[models.EventFunc]func(attrs) bool{
		models.Mint:     func(a attrs) bool { return a.from == "" && utils.IsBech32(a.to) },
		models.Burn:     func(a attrs) bool { return utils.IsBech32(a.from) && a.to == "" },
		models.Transfer: func(a attrs) bool { return utils.IsBech32(a.from) && utils.IsBech32(a.to) },
	}

	f, ok := validFuncs[models.EventFunc(event.Func)]
	if !ok {
		return models.TransferEvent{}, fmt.Errorf("invalid function")
	}
	if !f(attrs{from, to, value}) {
		return models.TransferEvent{}, fmt.Errorf("invalid attributes")
	}
	return models.TransferEvent{
		Func:    models.EventFunc(event.Func),
		PkgPath: event.PkgPath,
		From:    from,
		To:      to,
		Value:   valueInt,
	}, nil
}
