package bean

type Embedding struct {
	Object    string    `json:"object"`
	Embedding []float32 `json:"embedding"`
	Index     int       `json:"index"`
}

type EmbeddingRequest struct {
	Input []string `json:"input"`
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
