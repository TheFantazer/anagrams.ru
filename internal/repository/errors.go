package repository

import "errors"

var (
	ErrNotFound            = errors.New("record not found")
	ErrDuplicateResult     = errors.New("result already submitted for this session and player")
	ErrForeignKeyViolation = errors.New("foreign key constraint violation")
)
