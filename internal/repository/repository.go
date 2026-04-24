package repository

import (
	"context"
	"medvisitlog/cmd/common"
	"medvisitlog/internal/entities/patient"
)

type PatientRepo interface {
	Create(ctx context.Context, patient patient.Patient) error
	GetByID(ctx context.Context, id common.ID) (patient.Patient, error)
}
