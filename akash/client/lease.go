package client

import (
	"encoding/json"
	"strings"
	"terraform-provider-akash/akash/client/cli"
	"terraform-provider-akash/akash/client/types"
)

func (ak *AkashClient) CreateLease(seqs Seqs, provider string) (string, error) {
	cmd := cli.AkashCli(ak).Tx().Market().Lease().Create().
		SetDseq(seqs.Dseq).SetGseq(seqs.Gseq).SetOseq(seqs.Oseq).
		SetProvider(provider).SetOwner(ak.Config.AccountAddress).SetFrom(ak.Config.KeyName).
		DefaultGas().SetChainId(ak.Config.ChainId).SetKeyringBackend(ak.Config.KeyringBackend).
		SetNote(ak.transactionNote).AutoAccept().SetNode(ak.Config.Node).OutputJson()

	var out []byte
	var err error
	var transaction types.Transaction

	if err = ak.WaitForTransaction(func() (string, error) {
		out, err = cmd.Raw()
		if err != nil {
			return "", err
		}

		err = json.Unmarshal(out, &transaction)
		if err != nil {
			return "", err
		}

		err = json.NewDecoder(strings.NewReader(string(out))).Decode(&transaction)
		if err != nil {
			return "", err
		}

		return transaction.TxHash, nil
	}); err != nil {
		return "", err
	}

	return string(out), nil
}
