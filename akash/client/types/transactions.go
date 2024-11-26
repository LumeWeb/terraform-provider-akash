package types

import "errors"

type Transaction struct {
	Height string           `json:"height"`
	TxHash string          `json:"txhash"`
	Logs   []TransactionLog `json:"logs"`
	RawLog string           `json:"raw_log"`
}

type TransactionLog struct {
	Events []TransactionEvent `json:"events"`
}

type TransactionEvent struct {
	Type       string                     `json:"type"`
	Attributes TransactionEventAttributes `json:"attributes"`
}

type TransactionEventAttribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type TransactionEventAttributes []TransactionEventAttribute

func (a TransactionEventAttributes) Get(key string) (string, error) {
	for i, attr := range a {
		if attr.Key == key {
			return a[i].Value, nil
		}
	}

	return "", errors.New("attribute not found")
}

func (t Transaction) Failed() bool {
	return len(t.Logs) == 0
}

func (t Transaction) GetEventsByType(eventType string) []TransactionEvent {
	if len(t.Logs) == 0 {
		return nil
	}
	
	var matches []TransactionEvent
	for _, log := range t.Logs {
		for _, event := range log.Events {
			if event.Type == eventType {
				matches = append(matches, event)
			}
		}
	}
	return matches
}
