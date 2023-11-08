package bean

type Embedding struct {
	Object    string    `json:"object"`
	Embedding []float32 `json:"embedding"`
	Index     int       `json:"index"`
}

type EmbeddingRequest struct {
	Model    string      `json:"model"`
	Input    []string    `json:"input"`
	Type     int8        `json:"type"`
	ApiKey   string      `json:"api_key"`
	Endpoint string      `json:"endpoint"`
	ProxyCfg ProxyConfig `json:"proxy_config"`
}

type ProxyConfig struct {
	NeedProxy bool   `json:"need_proxy"`
	Address   string `json:"address"`
}

type EmbeddingResult struct {
	Model          string      `json:"model"`
	Object         string      `json:"object"`
	Input          []string    `json:"input"`
	Data           []Embedding `json:"data"`
	HttpStatusCode int         `json:"status_code"`
	ResponseBody   string      `json:"response_body"`
	Usage          Usage       `json:"usage"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
