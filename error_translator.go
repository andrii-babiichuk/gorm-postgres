package postgres

import (
	"encoding/json"

	"gorm.io/gorm"
)

// The error codes to map PostgreSQL errors to gorm errors, here is the PostgreSQL error codes reference https://www.postgresql.org/docs/current/errcodes-appendix.html.
var errCodes = map[string]error{
	"23505": gorm.ErrDuplicatedKey,
	"23503": gorm.ErrForeignKeyViolated,
	"42703": gorm.ErrInvalidField,
	"23514": gorm.ErrCheckConstraintViolated,
}

type ErrMessage struct {
	Code     string
	Severity string
	Message  string
}

// Translate it will translate the error to native gorm errors.
func (dialector Dialector) Translate(err error) error {
	parsedErr, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		return err
	}

	var errMsg ErrMessage
	unmarshalErr := json.Unmarshal(parsedErr, &errMsg)
	if unmarshalErr != nil {
		return err
	}

	if translatedErr, found := errCodes[errMsg.Code]; found {
		return translatedErr
	}
	return err
}
