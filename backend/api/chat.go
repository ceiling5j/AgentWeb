package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const deepSeekAPIURL = "https://api.deepseek.com/chat/completions"

type chatRequest struct {
	Message string `json:"message"`
}

type deepSeekMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type deepSeekRequest struct {
	Model       string            `json:"model"`
	Messages    []deepSeekMessage `json:"messages"`
	Temperature float64           `json:"temperature,omitempty"`
	Stream      bool              `json:"stream"`
}

type deepSeekResponse struct {
	Choices []struct {
		Message deepSeekMessage `json:"message"`
	} `json:"choices"`
}

type chatResponse struct {
	Reply string `json:"reply"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	setCORS(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "only POST is supported", http.StatusMethodNotAllowed)
		return
	}

	var req chatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json payload", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Message) == "" {
		http.Error(w, "message is required", http.StatusBadRequest)
		return
	}

	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		http.Error(w, "DEEPSEEK_API_KEY is not configured", http.StatusInternalServerError)
		return
	}

	payload := deepSeekRequest{
		Model: "deepseek-chat",
		Messages: []deepSeekMessage{
			{Role: "system", Content: "You are a helpful online AI agent."},
			{Role: "user", Content: req.Message},
		},
		Temperature: 0.7,
		Stream:      false,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "failed to build request", http.StatusInternalServerError)
		return
	}

	httpReq, err := http.NewRequest(http.MethodPost, deepSeekAPIURL, bytes.NewBuffer(body))
	if err != nil {
		http.Error(w, "failed to create upstream request", http.StatusInternalServerError)
		return
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		http.Error(w, "failed to reach DeepSeek API", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "failed to read upstream response", http.StatusBadGateway)
		return
	}

	if resp.StatusCode >= http.StatusBadRequest {
		w.WriteHeader(http.StatusBadGateway)
		w.Write(respBody)
		return
	}

	var deepResp deepSeekResponse
	if err := json.Unmarshal(respBody, &deepResp); err != nil {
		http.Error(w, "failed to parse upstream response", http.StatusBadGateway)
		return
	}

	if len(deepResp.Choices) == 0 {
		http.Error(w, "upstream returned no choices", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chatResponse{Reply: deepResp.Choices[0].Message.Content})
}

func setCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
}
