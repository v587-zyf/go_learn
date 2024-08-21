package mq

type Config struct {
	SId    int32   `json:"SId,omitempty"`
	SType  string  `json:"SType,omitempty"`
	Url    string  `json:"url,omitempty"`
	Method IMethod `json:"method,omitempty"`
}
