package types

type TxHandler func() (string, error)

type TxRequest struct {
	Handler TxHandler
	Result  chan TxResult
}

type TxResult struct {
	Err error
}
