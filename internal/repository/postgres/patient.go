package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"medvisitlog/internal/entities/patient"
)

//go:embed requests/insert_person.sql
var insertPersonSQL string

//go:embed requests/insert_patient.sql
var insertPatientSQL string

type PostgresPatientRepo struct {
	db *sql.DB
}

func (sgpr PostgresPatientRepo) Create(ctx context.Context, patient patient.Patient) error {

	var personID int
	err := sgpr.db.QueryRow(insertPersonSQL,
		patient.Person.GetFirstName(),
		patient.Person.GetLastName(),
		patient.Person.GetDateOfBirth(),
		patient.Person.GetAge(),
	).Scan(&personID)
	if err != nil {
		return err
	}
	err = sgpr.db.QueryRow(insertPatientSQL,
		personID,
		patient.GetCondition(),
		patient.GetDiagnosis(),
	).Err()
	if err != nil {
		return err
	}
	return nil
}
