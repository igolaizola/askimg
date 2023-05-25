package askimg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type replicateRequest struct {
	Version string         `json:"version"`
	Input   replicateInput `json:"input"`
}

type replicateInput struct {
	// Input image to query or caption
	Image string `json:"image"`
	// Select if you want to generate image captions instead of asking questions
	Caption bool `json:"caption"`
	// Question to ask about this image. Leave blank for captioning
	Question string `json:"question"`
	// Optional - previous questions and answers to be used as context for answering current question
	Context string `json:"context"`
	// Toggles the model using nucleus sampling to generate responses
	UseNucleusSampling bool `json:"use_nucleus_sampling"`
	// Temperature for use with nucleus sampling
	Temperature int `json:"temperature"`
}

type replicateResponse struct {
	ID          string         `json:"id"`
	Version     string         `json:"version"`
	Input       replicateInput `json:"input"`
	Logs        string         `json:"logs"`
	Output      string         `json:"output"`
	Error       string         `json:"error"`
	Status      string         `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	CompletedAt time.Time      `json:"completed_at"`
	Metrics     struct {
		PredictTime float64 `json:"predict_time"`
	} `json:"metrics"`
	URLs struct {
		Cancel string `json:"cancel"`
		Get    string `json:"get"`
	} `json:"urls"`
}

// Version for model https://replicate.com/andreasjansson/blip-2
const replicateVersion = "4b32258c42e9efd4288bb9910bc532a69727f9acd26aa08e175713a0a857a608"

type Config struct {
	Token              string
	Image              string
	Question           string
	Temperature        int
	Timeout            time.Duration
	UseNucleusSampling bool
}

// Ask asks a question about an image
func Ask(ctx context.Context, cfg *Config) (string, error) {
	// Validate config
	if cfg.Token == "" {
		return "", fmt.Errorf("askimg: token is required")
	}
	token := cfg.Token
	if cfg.Image == "" {
		return "", fmt.Errorf("askimg: image is required")
	}
	image := cfg.Image
	temperature := 1
	if cfg.Temperature > 0 {
		temperature = cfg.Temperature
	}

	// Add timeout to context if specified
	if cfg.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cfg.Timeout)
		defer cancel()
	}

	// Check if we are captioning or asking a question
	question := cfg.Question
	caption := question == ""

	// Set timeout

	// Create HTTP client
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create replicate data
	request := replicateRequest{
		Version: replicateVersion,
		Input: replicateInput{
			Image:       image,
			Caption:     caption,
			Question:    question,
			Temperature: temperature,
		},
	}

	// Create request
	body, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("askimg: couldn't marshal request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.replicate.com/v1/predictions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("askimg: couldn't create request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))
	req.Header.Set("Content-Type", "application/json")

	// Launch request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("askimg: couldn't do request: %w", err)
	}
	defer resp.Body.Close()
	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("askimg: couldn't read response body: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Println("askimg:", string(bodyResp))
		return "", fmt.Errorf("askimg: status code not ok: %d", resp.StatusCode)
	}

	// Unmarshal response
	var response replicateResponse
	if err := json.Unmarshal(bodyResp, &response); err != nil {
		return "", fmt.Errorf("askimg: couldn't unmarshal response: %w", err)
	}

	// Create request to get response
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, response.URLs.Get, nil)
	if err != nil {
		return "", fmt.Errorf("askimg: couldn't create request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))
	req.Header.Set("Content-Type", "application/json")

	// Wait for response
	for response.CompletedAt.IsZero() {
		// Wait before launching request
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(500 * time.Millisecond):
		}

		// Launch request to get response
		resp, err = client.Do(req)
		if err != nil {
			return "", fmt.Errorf("askimg: couldn't do request: %w", err)
		}
		defer resp.Body.Close()
		bodyResp, err = io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("askimg: couldn't read response body: %w", err)
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			log.Println("askimg:", string(bodyResp))
			return "", fmt.Errorf("askimg: status code not ok: %d", resp.StatusCode)
		}

		// Unmarshal response
		if err := json.Unmarshal(bodyResp, &response); err != nil {
			log.Println("askimg:", string(bodyResp))
			return "", fmt.Errorf("askimg: couldn't unmarshal response: %w", err)
		}
	}

	// Return response
	return response.Output, nil
}
