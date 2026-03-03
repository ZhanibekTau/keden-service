package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"keden-service/back/internal/configs/structures"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type AIService struct {
	config     *structures.AIConfig
	httpClient *http.Client
}

func NewAIService(cfg *structures.AIConfig) *AIService {
	timeout, err := time.ParseDuration(cfg.Timeout)
	if err != nil {
		timeout = 120 * time.Second
	}

	return &AIService{
		config: cfg,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

type AIResponse struct {
	DocumentType string                 `json:"document_type"`
	Fields       map[string]interface{} `json:"fields"`
	Items        []map[string]interface{} `json:"items"`
}

func (s *AIService) ProcessDocument(ctx context.Context, pdfData []byte) (*AIResponse, error) {
	if s.config.ServiceURL == "" {
		logrus.Warn("AI service URL not configured, using mock response")
		return s.mockResponse(), nil
	}

	body, err := json.Marshal(map[string]interface{}{
		"document": pdfData,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.config.ServiceURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		logrus.Warnf("AI service unavailable, using mock response: %v", err)
		return s.mockResponse(), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("AI service returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var aiResp AIResponse
	if err := json.NewDecoder(resp.Body).Decode(&aiResp); err != nil {
		return nil, fmt.Errorf("failed to decode AI response: %w", err)
	}

	return &aiResp, nil
}

func (s *AIService) mockResponse() *AIResponse {
	return &AIResponse{
		DocumentType: "customs_declaration",
		Fields: map[string]interface{}{
			"declaration_number": "10001010/060226/0001234",
			"date":              time.Now().Format("02.01.2006"),
			"sender":            "ТОО \"Экспорт Компани\"",
			"receiver":          "ТОО \"Импорт Компани\"",
			"country_origin":    "Китай",
			"country_dest":      "Казахстан",
			"currency":          "USD",
			"total_value":       25000.00,
			"customs_value":     27500.00,
		},
		Items: []map[string]interface{}{
			{
				"number":      1,
				"hs_code":     "8471300000",
				"description": "Ноутбуки портативные",
				"quantity":    100,
				"unit":        "шт",
				"weight_net":  250.0,
				"weight_gross": 300.0,
				"value":       15000.00,
				"duty_rate":   "0%",
				"vat_rate":    "12%",
			},
			{
				"number":      2,
				"hs_code":     "8517120000",
				"description": "Смартфоны",
				"quantity":    200,
				"unit":        "шт",
				"weight_net":  40.0,
				"weight_gross": 55.0,
				"value":       10000.00,
				"duty_rate":   "0%",
				"vat_rate":    "12%",
			},
		},
	}
}
