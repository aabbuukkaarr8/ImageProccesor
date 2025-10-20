package handler

type Request struct {
	Action string            `json:"action"`
	Params map[string]string `json:"params"`
}
