package indexer

import (
	"encoding/json"
	"fmt"

	"github.com/yanmoyy/tbi/internal/models"
)

func convert[T, V any](slice []T, convertor func(T) (V, error)) ([]V, error) {
	result := make([]V, len(slice))
	for i, v := range slice {
		v, err := convertor(v)
		if err != nil {
			return nil, err
		}
		result[i] = v
	}
	return result, nil
}

func blockConvertor(b block) (models.Block, error) {
	return b.toModel(), nil
}

func (b *block) toModel() models.Block {
	return models.Block{
		Hash:     b.Hash,
		Height:   b.Height,
		Time:     b.Time,
		NumTxs:   b.NumTxs,
		TotalTxs: b.TotalTxs,
	}
}

func transactionConvertor(t transaction) (models.Transaction, error) {
	messages, err := convert(t.Messages, transactionMessageConvertor)
	if err != nil {
		return models.Transaction{}, fmt.Errorf("convert messages: %w", err)
	}
	response, err := transactionResponseConvertor(t.Response)
	if err != nil {
		return models.Transaction{}, fmt.Errorf("convert response: %w", err)
	}
	return models.Transaction{
		Index:       t.Index,
		Hash:        t.Hash,
		Success:     t.Success,
		BlockHeight: t.BlockHeight,
		GasWanted:   t.GasWanted,
		GasUsed:     t.GasUsed,
		Memo:        t.Memo,
		GasFee: models.Coin{
			Amount: t.GasFee.Amount,
			Denom:  t.GasFee.Denom,
		},
		Messages: messages,
		Response: response,
	}, nil
}

func transactionMessageConvertor(tm transactionMessage) (models.TransactionMessage, error) {
	var value any

	switch tm.TypeURL {
	case "send":
		value = tm.Value.BankMsgSend
	case "exec":
		value = tm.Value.MsgRun
	case "add_package":
		value = tm.Value.MsgAddPackage
	case "run":
		value = tm.Value.MsgRun
	}
	valueJson, err := json.Marshal(value)
	if err != nil {
		return models.TransactionMessage{}, fmt.Errorf("marshal value: %w", err)
	}
	return models.TransactionMessage{
		Route:   tm.Route,
		TypeURL: tm.TypeURL,
		Value:   valueJson,
	}, nil
}

func transactionResponseConvertor(tr transactionResponse) (models.TransactionResponse, error) {
	events, err := convert(tr.Events, gnoEventConvertor)
	if err != nil {
		return models.TransactionResponse{}, err
	}
	return models.TransactionResponse{
		Log:    tr.Log,
		Info:   tr.Info,
		Error:  tr.Error,
		Data:   tr.Data,
		Events: events,
	}, nil
}

func gnoEventConvertor(re responseEvent) (models.GnoEvent, error) {
	attrs, err := convert(re.GNOEvent.Attrs, gnoEventAttrConvertor)
	if err != nil {
		return models.GnoEvent{}, fmt.Errorf("convert attrs: %w", err)
	}
	return models.GnoEvent{
		Type:    re.GNOEvent.Type,
		Func:    re.GNOEvent.Func,
		PkgPath: re.GNOEvent.PkgPath,
		Attrs:   attrs,
	}, nil
}

func gnoEventAttrConvertor(ga gnoEventAttr) (models.GnoEventAttr, error) {
	return models.GnoEventAttr{
		Key:   ga.Key,
		Value: ga.Value,
	}, nil
}
