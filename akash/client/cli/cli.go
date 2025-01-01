package cli

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"os"
)

func envExists(key string) bool {
	_, exists := os.LookupEnv(key)
	return exists
}

type AkashCommand struct {
	client      AkashCliClient
	ctx         context.Context
	Content     []string
	blockHeight string
}

type AkashCliClient interface {
	GetContext() context.Context
	GetPath() string
}

func AkashCli(client AkashCliClient) AkashCommand {
	path := client.GetPath()
	if path == "" {
		path = "provider-services"
	}

	return AkashCommand{
		client:  client,
		ctx:     client.GetContext(),
		Content: []string{path},
	}
}

func (c AkashCommand) Tx() AkashCommand {
	return c.append("tx")
}

func (c AkashCommand) Query() AkashCommand {
	return c.append("query")
}

func (c AkashCommand) SetHash(hash string) AkashCommand {
	return c.append(hash)
}

func (c AkashCommand) QueryTx() AkashCommand {
	return c.append("tx")
}

func (c AkashCommand) Deployment() AkashCommand {
	return c.append("deployment")
}

func (c AkashCommand) Get() AkashCommand {
	return c.append("get")
}

func (c AkashCommand) Create() AkashCommand {
	return c.append("create")
}

func (c AkashCommand) Update() AkashCommand {
	return c.append("update")
}

func (c AkashCommand) LeaseStatus() AkashCommand {
	return c.append("lease-status")
}

func (c AkashCommand) SendManifest(path string) AkashCommand {
	return c.append("send-manifest").append(path)
}

func (c AkashCommand) Close() AkashCommand {
	return c.append("close")
}

func (c AkashCommand) Market() AkashCommand {
	return c.append("market")
}

func (c AkashCommand) Provider() AkashCommand {
	return c.append("provider")
}

func (c AkashCommand) Node() AkashCommand {
	return c.append("node")
}

func (c AkashCommand) Status() AkashCommand {
	return c.append("status")
}

func (c AkashCommand) Bid() AkashCommand {
	return c.append("bid")
}

func (c AkashCommand) List() AkashCommand {
	return c.append("list")
}

func (c AkashCommand) Lease() AkashCommand {
	return c.append("lease")
}

func (c AkashCommand) Manifest(path string) AkashCommand {
	return c.append(path)
}

/** OPTIONS **/

func (c AkashCommand) SetDseq(dseq string) AkashCommand {
	return c.append("--dseq").append(dseq)
}

func (c AkashCommand) SetOseq(oseq string) AkashCommand {
	return c.append("--oseq").append(oseq)
}

func (c AkashCommand) SetGseq(gseq string) AkashCommand {
	return c.append("--gseq").append(gseq)
}

func (c AkashCommand) SetProvider(provider string) AkashCommand {
	return c.append("--provider").append(provider)
}

func (c AkashCommand) SetHome(home string) AkashCommand {
	return c.append("--home").append(home)
}

func (c AkashCommand) SetOwner(owner string) AkashCommand {
	return c.append("--owner").append(owner)
}

func (c AkashCommand) SetFees(amount int64) AkashCommand {
	return c.append("--fees").append(fmt.Sprintf("%duakt", amount))
}

func (c AkashCommand) SetFrom(key string) AkashCommand {
	return c.append("--from").append(key)
}

func (c AkashCommand) GasAuto() AkashCommand {
	if !envExists("AKASH_GAS") {
		return c.append("--gas=auto")
	}
	return c
}
func (c AkashCommand) SetGasAdjustment(adjustment float32) AkashCommand {
	if !envExists("AKASH_GAS_ADJUSTMENT") {
		return c.append(fmt.Sprintf("--gas-adjustment=%2f", adjustment))
	}
	return c
}

func (c AkashCommand) SetGasPrices() AkashCommand {
	if !envExists("AKASH_GAS_PRICES") {
		return c.append("--gas-prices=0.025uakt")
	}
	return c
}

func (c AkashCommand) SetChainId(chainId string) AkashCommand {
	if !envExists("AKASH_CHAIN_ID") {
		return c.append("--chain-id").append(chainId)
	}
	return c
}

func (c AkashCommand) SetNode(node string) AkashCommand {
	if !envExists("AKASH_NODE") {
		return c.append("--node").append(node)
	}
	return c
}

func (c AkashCommand) SetKeyringBackend(keyringBackend string) AkashCommand {
	if !envExists("AKASH_KEYRING_BACKEND") {
		return c.append("--keyring-backend").append(keyringBackend)
	}
	return c
}

func (c AkashCommand) SetNote(note string) AkashCommand {
	return c.append(fmt.Sprintf("--note=\"%s\"", note))
}

func (c AkashCommand) SetDepositorAccount(account string) AkashCommand {
	return c.append("--depositor-account").append(account)
}

func (c AkashCommand) SetSignMode(mode string) AkashCommand {
	if !envExists("AKASH_SIGN_MODE") {
		supportedModes := map[string]bool{
			"default":    true,
			"amino-json": true,
		}

		if _, ok := supportedModes[mode]; !ok {
			tflog.Error(c.ctx, fmt.Sprintf("Mode '%s' not supported", mode))
			return c
		}

		return c.append("--sign-mode").append(mode)
	}
	return c
}

func (c AkashCommand) AutoAccept() AkashCommand {
	return c.append("-y")
}

func (c AkashCommand) OutputJson() AkashCommand {
	return c.append("-o").append("json")
}

func (c AkashCommand) Headless() []string {
	return c.Content[1:]
}

func (c AkashCommand) append(str string) AkashCommand {
	c.Content = append(c.Content, str)
	return c
}

func (c *AkashCommand) SetBlockHeight(height string) {
	c.blockHeight = height
}

func (c AkashCommand) GetBlockHeight() string {
	return c.blockHeight
}
