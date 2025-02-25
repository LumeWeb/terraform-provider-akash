package client

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-akash/akash/client/cli"
	"terraform-provider-akash/akash/client/types"
)

type Seqs struct {
	Dseq string
	Gseq string
	Oseq string
}

func (ak *AkashClient) GetDeployments(owner string) ([]types.DeploymentId, error) {
	panic("Not implemented")
}

func (ak *AkashClient) GetDeployment(dseq string, owner string) (types.Deployment, error) {
	cmd := cli.AkashCli(ak).Query().Deployment().Get().SetOwner(owner).SetDseq(dseq).
		SetChainId(ak.Config.ChainId).SetNode(ak.Config.Node).OutputJson()

	deployment := types.Deployment{}
	err := cmd.DecodeJson(&deployment)
	if err != nil {
		return types.Deployment{}, err
	}

	return deployment, nil
}

func (ak *AkashClient) CreateDeployment(manifestLocation string) (Seqs, error) {

	tflog.Info(ak.ctx, "Creating deployment")
	// Create deployment using the file created with the SDL
	attributes, err := transactionCreateDeployment(ak, manifestLocation)
	if err != nil {
		tflog.Error(ak.ctx, "Failed creating deployment")
		tflog.Debug(ak.ctx, fmt.Sprintf("%s", err))
		return Seqs{}, err
	}

	dseq, _ := attributes.Get("dseq")
	gseq, _ := attributes.Get("gseq")
	oseq, _ := attributes.Get("oseq")

	tflog.Info(ak.ctx, fmt.Sprintf("Deployment created with DSEQ=%s GSEQ=%s OSEQ=%s", dseq, gseq, oseq))

	return Seqs{dseq, gseq, oseq}, nil
}

// Perform the transaction to create the deployment and return either the DSEQ or an error.
func transactionCreateDeployment(ak *AkashClient, manifestLocation string) (types.TransactionEventAttributes, error) {
	cmd := cli.AkashCli(ak).Tx().Deployment().Create().Manifest(manifestLocation).
		DefaultGas().AutoAccept().SetFrom(ak.Config.KeyName).
		SetKeyringBackend(ak.Config.KeyringBackend).
		SetChainId(ak.Config.ChainId).
		SetNode(ak.Config.Node).
		SetNote(ak.transactionNote)

	if ak.Config.DepositorAccount != "" {
		cmd = cmd.SetDepositorAccount(ak.Config.DepositorAccount)
	}
	
	if ak.Config.FeeAccount != "" {
		cmd = cmd.SetFeeAccount(ak.Config.FeeAccount)
	}

	cmd = cmd.OutputJson()

	var transaction types.Transaction

	err := ak.WaitForTransaction(func() (string, error) {
		if err := cmd.DecodeJson(&transaction); err != nil {
			return "", err
		}
		return transaction.TxHash, nil
	})
	if err != nil {
		return nil, err
	}

	if transaction.Failed() {
		return nil, fmt.Errorf("transaction failed: %s", transaction.RawLog)
	}

	// Query the transaction to get its events
	queryCmd := cli.AkashCli(ak).Query().QueryTx().SetHash(transaction.TxHash).
		SetNode(ak.Config.Node).
		OutputJson()

	if err := queryCmd.DecodeJson(&transaction); err != nil {
		return nil, err
	}

	if len(transaction.Logs) == 0 {
		return nil, errors.New("no transaction logs found")
	}

	return transaction.Logs[0].Events[0].Attributes, nil
}

func (ak *AkashClient) DeleteDeployment(dseq string, owner string) error {
	cmd := cli.AkashCli(ak).Tx().Deployment().Close().
		SetDseq(dseq).SetOwner(owner).SetFrom(ak.Config.KeyName).
		DefaultGas().SetChainId(ak.Config.ChainId).
		SetKeyringBackend(ak.Config.KeyringBackend).
		SetNode(ak.Config.Node).
		SetNote(ak.transactionNote).AutoAccept().OutputJson()

	var transaction types.Transaction
	err := ak.WaitForTransaction(func() (string, error) {
		if err := cmd.DecodeJson(&transaction); err != nil {
			return "", err
		}
		return transaction.TxHash, nil
	})
	if err != nil {
		return err
	}

	tflog.Debug(ak.ctx, fmt.Sprintf("Response: %s", transaction.RawLog))

	return nil
}

func (ak *AkashClient) UpdateDeployment(dseq string, manifestLocation string) error {
	cmd := cli.AkashCli(ak).Tx().Deployment().Update().Manifest(manifestLocation).
		SetDseq(dseq).SetFrom(ak.Config.KeyName).
		SetNode(ak.Config.Node).
		SetKeyringBackend(ak.Config.KeyringBackend).
		SetChainId(ak.Config.ChainId).
		SetNote(ak.transactionNote).
		GasAuto().SetGasAdjustment(1.5).SetGasPrices().SetSignMode("amino-json").AutoAccept().OutputJson()

	var transaction types.Transaction
	err := ak.WaitForTransaction(func() (string, error) {
		if err := cmd.DecodeJson(&transaction); err != nil {
			return "", err
		}
		return transaction.TxHash, nil
	})

	if err != nil {
		return err
	}

	tflog.Debug(ak.ctx, fmt.Sprintf("Response: %s", transaction.RawLog))

	return nil
}
