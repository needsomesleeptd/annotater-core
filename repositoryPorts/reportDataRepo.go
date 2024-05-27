package repository

import (
	"github.com/google/uuid"
	"github.com/needsomesleeptd/annotater-core/models"
)

type IReportDataRepository interface {
	AddDocument(doc *models.ErrorReport) error
	DeleteReportByID(id uuid.UUID) error
	GetDocumentByID(id uuid.UUID) (*models.ErrorReport, error)
}
