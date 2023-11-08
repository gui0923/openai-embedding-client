package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gui0923/openai-embedding-client/bean"
	"github.com/pkoukk/tiktoken-go"
)

type embeddingClient struct {
	maxTokens     int
	maxInputNum   int
	tiktoken      tiktoken.Tiktoken
	openaiService openAIEmbeddingProcessService
	azureService  azureOpenAIEmbeddingProcessService
}

func NewEmbeddingClient(maxTokens int, maxInputNum int) *embeddingClient {
	openaiService := &openAIEmbeddingProcessService{}
	azureService := &azureOpenAIEmbeddingProcessService{}
	t, err := tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		panic(err)
	}
	return &embeddingClient{
		maxTokens:     maxTokens,
		maxInputNum:   maxInputNum,
		tiktoken:      *t,
		openaiService: *openaiService,
		azureService:  *azureService,
	}
}

func (client *embeddingClient) EmbeddingRequest(request *bean.EmbeddingRequest) (bean.EmbeddingResult, error) {
	inputs := request.Input
	inputMap := make(map[string]interface{}, 0)
	for _, v := range inputs {
		inputMap[v] = nil
	}
	keys := make([]string, 0, len(inputMap))
	for k := range inputMap {
		keys = append(keys, k)
	}
	request.Input = keys
	fmt.Printf("keys: %v\n", keys)
	if len(keys) > client.maxInputNum {
		return bean.EmbeddingResult{}, errors.New("exceeded maximum input size limit.")
	}
	var inputTotalTokens = 0
	for _, v := range keys {
		inputTotalTokens += len(client.tiktoken.Encode(v, nil, nil))
	}
	if inputTotalTokens >= client.maxTokens {
		return bean.EmbeddingResult{}, errors.New("exceeded maximum tokens limit.")
	}
	var url string
	var header map[string]string
	var content map[string]interface{}
	if request.Type == 0 {
		a, b, d := client.openaiService.GenerateRequest(request)
		url = a
		header = b
		content = d
	} else {
		a, b, d := client.azureService.GenerateRequest(request)
		url = a
		header = b
		content = d
	}
	b, err := json.Marshal(content)
	if err != nil {
		return bean.EmbeddingResult{}, err
	}
	r, err := http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil {
		return bean.EmbeddingResult{}, err
	}
	for k, v := range header {
		r.Header.Add(k, v)
	}

	httpclient := &http.Client{}
	response, err := httpclient.Do(r)
	if err != nil {
		return bean.EmbeddingResult{}, err
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	if request.Type == 0 {
		res, err := client.openaiService.ConvertEmbeddingResult(string(body))
		if err == nil {
			res.Input = keys
		}
		return res, err
	} else {
		res, err := client.azureService.ConvertEmbeddingResult(string(body))
		if err == nil {
			res.Input = keys
		}
		return res, err
	}
}
