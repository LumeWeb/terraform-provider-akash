package client

import (
	"context"
	"errors"
	"fmt"
	"terraform-provider-akash/akash/client/cli"
	"terraform-provider-akash/akash/client/types"
	"time"
)

type AkashClient struct {
	ctx             context.Context
	Config          AkashProviderConfiguration
	transactionNote string
	txQueue         chan types.TxRequest
}

type AkashProviderConfiguration struct {
	KeyName        string
	KeyringBackend string
	AccountAddress string
	Net            string
	Version        string
	ChainId        string
	Node           string
	Home           string
	Path           string
	ProvidersApi   string
}

func (ak *AkashClient) GetContext() context.Context {
	return ak.ctx
}

func (ak *AkashClient) GetPath() string {
	return ak.Config.Path
}

func (ak *AkashClient) SetGlobalTransactionNote(note string) {
	ak.transactionNote = note
}

func New(ctx context.Context, configuration AkashProviderConfiguration) *AkashClient {
	client := &AkashClient{
		ctx:     ctx,
		Config:  configuration,
		txQueue: make(chan types.TxRequest, 100),
	}

	go client.txWorker()

	return client
}

func (ak *AkashClient) txWorker() {
	for req := range ak.txQueue {
		txId, err := req.Handler()
		if err != nil {
			req.Result <- types.TxResult{Err: err}
		}

		err = ak.waitForTx(txId)
		req.Result <- types.TxResult{Err: err}
		close(req.Result)
	}
}

func (ak *AkashClient) waitForTx(txHash string) error {
	timer := time.NewTimer(90 * time.Second)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			return errors.New("timeout waiting for transaction confirmation")
		default:
			cmd := cli.AkashCli(ak).Query().Tx().SetHash(txHash).
				SetNode(ak.Config.Node).
				OutputJson()

			var txResponse types.Transaction
			if err := cmd.DecodeJson(&txResponse); err != nil {
				time.Sleep(500 * time.Millisecond)
				continue
			}

			if txResponse.Failed() {
				return fmt.Errorf("transaction failed: %s", txResponse.RawLog)
			}

			return nil
		}
	}
}

func (ak *AkashClient) WaitForTransaction(handler types.TxHandler) error {
	resultChan := make(chan types.TxResult, 1)
	ak.txQueue <- types.TxRequest{
		Handler: handler,
		Result:  resultChan,
	}

	timer := time.NewTimer(30 * time.Second)
	defer timer.Stop()

	select {
	case <-timer.C:
		return errors.New("timeout waiting for transaction confirmation")
	case result := <-resultChan:
		return result.Err
	}
}
