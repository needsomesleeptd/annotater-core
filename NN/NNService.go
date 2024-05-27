package nn

import "github.com/needsomesleeptd/annotater-core/models"

type INeuralNetwork interface {
	Predict(document models.DocumentData) ([]models.Markup, error)
}
