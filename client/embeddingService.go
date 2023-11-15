package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"math"
	"net/http"

	"github.com/gui0923/openai-embedding-client/bean"
	"github.com/pkoukk/tiktoken-go"
	"golang.org/x/time/rate"
)

type EmbeddingClient struct {
	maxTokensPer   int
	maxInputNumPer int
	numberLimiter  rate.Limiter
	tokensLimiter  rate.Limiter
	tiktoken       tiktoken.Tiktoken
	openaiService  openAIEmbeddingProcessService
	azureService   azureOpenAIEmbeddingProcessService
}

func NewEmbeddingClient(maxTokens int, maxInputNum int) *EmbeddingClient {
	return NewLimterEmbeddingClient(maxInputNum, maxInputNum, math.MaxInt32, math.MaxInt32)
}

func NewLimterEmbeddingClient(maxTokens int, maxInputNum int, requestNumberPerMinute int, requestTokensPerMinute int) *EmbeddingClient {
	openaiService := &openAIEmbeddingProcessService{}
	azureService := &azureOpenAIEmbeddingProcessService{}
	t, err := tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		panic(err)
	}
	c := &EmbeddingClient{
		maxTokensPer:   maxTokens,
		maxInputNumPer: maxInputNum,
		tiktoken:       *t,
		openaiService:  *openaiService,
		azureService:   *azureService,
	}
	c.tokensLimiter = *rate.NewLimiter(rate.Limit(requestTokensPerMinute/60.0), requestTokensPerMinute)
	c.numberLimiter = *rate.NewLimiter(rate.Limit(requestNumberPerMinute/60.0), requestNumberPerMinute)
	return c
}

func (client *EmbeddingClient) EmbeddingRequest(request *bean.EmbeddingRequest) (bean.EmbeddingResult, error) {
	inputs := request.Input
	inputMap := make(map[string]interface{}, 0)
	for _, v := range inputs {
		inputMap[v] = nil
	}
	request.Input = make([]string, 0, len(inputMap))
	for k := range inputMap {
		request.Input = append(request.Input, k)
	}
	if len(request.Input) > client.maxInputNumPer {
		return bean.EmbeddingResult{}, errors.New("exceeded maximum input size limit")
	}
	var inputTotalTokens = 0
	for _, v := range request.Input {
		inputTotalTokens += len(client.tiktoken.Encode(v, nil, nil))
	}
	if inputTotalTokens >= client.maxTokensPer {
		return bean.EmbeddingResult{}, errors.New("exceeded maximum tokens limit")
	}
	err := client.numberLimiter.Wait(context.Background())
	if err != nil {
		return bean.EmbeddingResult{}, err
	}
	err2 := client.tokensLimiter.WaitN(context.Background(), inputTotalTokens)
	if err2 != nil {
		return bean.EmbeddingResult{}, err2
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
			res.Input = request.Input
		}
		return res, err
	} else {
		res, err := client.azureService.ConvertEmbeddingResult(string(body))
		if err == nil {
			res.Input = request.Input
		}
		return res, err
	}
}
