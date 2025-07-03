package embedding

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type OllamaClient struct {
	ApiBase string
	Model   string
}

type ollamaEmbedRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type ollamaEmbedResponse struct {
	Model      string    `json:"model"`
	Embeddings []float32 `json:"embeddings"`
}

func NewOllamaClient(url string, model string) *OllamaClient {
	return &OllamaClient{
		ApiBase: strings.TrimRight(url, "/"),
		Model:   model,
	}
}

func (c *OllamaClient) Embed(ctx context.Context, input string) ([]float32, error) {
	url := fmt.Sprintf("%s/api/embed", c.ApiBase)

	reqBody, err := json.Marshal(ollamaEmbedRequest{Model: c.Model, Input: input})
	if err != nil {
		return []float32{}, err
	}

	respBytes, err := c.sendRequest(ctx, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	var respObj ollamaEmbedResponse
	err = json.Unmarshal(respBytes, &respObj)
	if err != nil {
		return []float32{}, fmt.Errorf("error decoding response: %v, body: %s", err, string(respBytes))
	}

	return respObj.Embeddings, nil // need check error from server
}

func (c *OllamaClient) sendRequest(ctx context.Context, url string, body []byte) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
