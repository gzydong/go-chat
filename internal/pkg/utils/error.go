package utils

import (
	"errors"

	"gorm.io/gorm"
)

func IsSqlNoRows(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
