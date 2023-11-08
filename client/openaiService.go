package client

import (
	"encoding/json"

	"github.com/gui0923/openai-embedding-client/bean"
	"github.com/gui0923/openai-embedding-client/constant"
)

type EmbeddingProcessService interface {
	GenerateRequest(request *bean.EmbeddingRequest) (string, map[string]string, map[string]interface{})

	ConvertEmbeddingResult(responseBody string) bean.EmbeddingResult
}

type OpenAIEmbeddingProcessService struct {
}

type AzureOpenAIEmbeddingProcessService struct {
	OpenAIEmbeddingProcessService
}

func (service *OpenAIEmbeddingProcessService) GenerateRequest(request *bean.EmbeddingRequest) (string, map[string]string, map[string]interface{}) {
	header := make(map[string]string, 0)
	header[constant.AUTHORIZATION_KEY] = constant.BEARER_PREFIX + request.ApiKey
	header[constant.CONTENT_TYPE] = constant.JSON_TYPE

	requestBody := make(map[string]interface{}, 0)
	requestBody[constant.MODEL_KEY] = request.Model
	if len(request.Input) == 1 {
		requestBody[constant.INPUT_KEY] = request.Input[0]
	} else {
		requestBody[constant.INPUT_KEY] = request.Input
	}
	return request.Endpoint, header, requestBody
}

func (service *AzureOpenAIEmbeddingProcessService) GenerateRequest(request *bean.EmbeddingRequest) (string, map[string]string, map[string]interface{}) {
	header := make(map[string]string, 0)
	header[constant.API_KEY] = request.ApiKey
	header[constant.CONTENT_TYPE] = constant.JSON_TYPE

	requestBody := make(map[string]interface{}, 0)
	if len(request.Input) == 1 {
		requestBody[constant.INPUT_KEY] = request.Input[0]
	} else {
		requestBody[constant.INPUT_KEY] = request.Input
	}
	return request.Endpoint, header, requestBody
}

func (service *OpenAIEmbeddingProcessService) ConvertEmbeddingResult(responseBody string) (bean.EmbeddingResult, error) {
	res := &bean.EmbeddingResult{}
	err := json.Unmarshal([]byte(responseBody), res)
	return *res, err
}
