package client

import (
	"encoding/json"

	"github.com/gui0923/openai-embedding-client/bean"
	"github.com/gui0923/openai-embedding-client/constant"
)

type embeddingProcessService interface {
	GenerateRequest(request *bean.EmbeddingRequest) (string, map[string]string, map[string]interface{})

	ConvertEmbeddingResult(responseBody string) bean.EmbeddingResult
}

type openAIEmbeddingProcessService struct {
}

type azureOpenAIEmbeddingProcessService struct {
	openAIEmbeddingProcessService
}

func (service *openAIEmbeddingProcessService) GenerateRequest(request *bean.EmbeddingRequest, config *EmbeddingClientConfig) (string, map[string]string, map[string]interface{}) {
	header := make(map[string]string, 0)
	header[constant.AUTHORIZATION_KEY] = constant.BEARER_PREFIX + config.ApiKey
	header[constant.CONTENT_TYPE] = constant.JSON_TYPE

	requestBody := make(map[string]interface{})
	requestBody[constant.MODEL_KEY] = config.Model
	if len(request.Input) == 1 {
		requestBody[constant.INPUT_KEY] = request.Input[0]
	} else {
		requestBody[constant.INPUT_KEY] = request.Input
	}
	return config.Endpoint, header, requestBody
}

func (service *azureOpenAIEmbeddingProcessService) GenerateRequest(request *bean.EmbeddingRequest, config *EmbeddingClientConfig) (string, map[string]string, map[string]interface{}) {
	header := make(map[string]string)
	header[constant.API_KEY] = config.ApiKey
	header[constant.CONTENT_TYPE] = constant.JSON_TYPE

	requestBody := make(map[string]interface{}, 0)
	if len(request.Input) == 1 {
		requestBody[constant.INPUT_KEY] = request.Input[0]
	} else {
		requestBody[constant.INPUT_KEY] = request.Input
	}
	return config.Endpoint, header, requestBody
}

func (service *openAIEmbeddingProcessService) ConvertEmbeddingResult(responseBody string) (bean.EmbeddingResult, error) {
	res := &bean.EmbeddingResult{}
	err := json.Unmarshal([]byte(responseBody), res)
	return *res, err
}
