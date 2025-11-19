package vo

type RemoteNodeInfo struct {
	Ip      string `json:"ip"`
	ApiPort uint64 `json:"apiPort"`
	WsPort  uint64 `json:"wsPort"`
}
