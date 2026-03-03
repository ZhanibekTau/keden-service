package mq

import (
	"context"
	"encoding/json"
	"keden-service/back/internal/services/ai"
	"keden-service/back/internal/services/document"
	"keden-service/back/internal/services/excel"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/sirupsen/logrus"
)

type DocumentProcessorHandler struct {
	docService   *document.DocumentService
	aiService    *ai.AIService
	excelService *excel.ExcelService
}

func NewDocumentProcessorHandler(
	ds *document.DocumentService,
	as *ai.AIService,
	es *excel.ExcelService,
) *DocumentProcessorHandler {
	return &DocumentProcessorHandler{
		docService:   ds,
		aiService:    as,
		excelService: es,
	}
}

type ProcessingMessage struct {
	DocumentID uint `json:"document_id"`
}

func (h *DocumentProcessorHandler) ConsumeMessage(ctx context.Context, msg *message.Message) error {
	var body ProcessingMessage
	if err := json.Unmarshal(msg.Payload, &body); err != nil {
		logrus.Errorf("Failed to unmarshal message: %v", err)
		return err
	}

	logrus.Infof("Processing document ID: %d", body.DocumentID)

	if err := h.docService.MarkDocumentProcessing(ctx, body.DocumentID); err != nil {
		logrus.Errorf("Failed to mark document as processing: %v", err)
		return err
	}

	aiResp, err := h.aiService.ProcessDocument(ctx, nil)
	if err != nil {
		errMsg := "AI processing failed: " + err.Error()
		logrus.Error(errMsg)
		_ = h.docService.MarkDocumentError(ctx, body.DocumentID, errMsg)
		return err
	}

	excelData, err := h.excelService.GenerateFromAIResponse(aiResp)
	if err != nil {
		errMsg := "Excel generation failed: " + err.Error()
		logrus.Error(errMsg)
		_ = h.docService.MarkDocumentError(ctx, body.DocumentID, errMsg)
		return err
	}

	aiJSON, _ := json.Marshal(aiResp)
	if err := h.docService.UpdateDocumentAfterAI(ctx, body.DocumentID, string(aiJSON), excelData); err != nil {
		errMsg := "Failed to save results: " + err.Error()
		logrus.Error(errMsg)
		_ = h.docService.MarkDocumentError(ctx, body.DocumentID, errMsg)
		return err
	}

	logrus.Infof("Document %d processed successfully", body.DocumentID)
	return nil
}
