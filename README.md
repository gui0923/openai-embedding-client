# openai-embedding-client
Openai and Azure embedding client

## Install

``` bash
go get github.com/gui0923/openai-embedding-client
```

## Usage

Apply your openai account, such as:
```
"model": "text-embedding-ada-002"
"api_key": "sk-BB*******************************"
"endpoint": https://api.openai.com/v1/embeddings
```
Or apply your azure account of microsoft, such as

```
"api_key":"**************************"
"endpoint": https://********.openai.azure.com/openai/deployments/******-embedding-v1/embeddings?api-version=2022-12-01
```

### Making requests

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/gui0923/openai-embedding-client/bean"
	"github.com/gui0923/openai-embedding-client/client"
)

func main() {
	embeddingClient := client.NewEmbeddingClient(8191, 16)
	request := &bean.EmbeddingRequest{
		Endpoint: "https://********.openai.azure.com/openai/deployments/******-embedding-v1/embeddings?api-version=2022-12-01",
		Type:     1,
		ApiKey:   "******************************",
		Input:    []string{"china", "Your text string goes here"},
	}
	embeddingResult, err := embeddingClient.EmbeddingRequest(request)
	if err != nil {
		panic(err)
	}
	b, _ := json.Marshal(embeddingResult)
	fmt.Println(string(b))
}
```