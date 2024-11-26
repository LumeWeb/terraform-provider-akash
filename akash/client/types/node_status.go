package types

type NodeStatus struct {
	SyncInfo struct {
		LatestBlockHeight string `json:"latest_block_height"`
		LatestBlockTime   string `json:"latest_block_time"` 
		CatchingUp        bool   `json:"catching_up"`
	} `json:"SyncInfo"`
}
