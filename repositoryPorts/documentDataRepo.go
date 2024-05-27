package repository

import (
	"github.com/google/uuid"
	"github.com/needsomesleeptd/annotater-core/models"
)

type IDocumentDataRepository interface {
	AddDocument(doc *models.DocumentData) error
	DeleteDocumentByID(id uuid.UUID) error
	GetDocumentByID(id uuid.UUID) (*models.DocumentData, error)
}
