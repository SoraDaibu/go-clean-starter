package common

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// Helper functions for type conversion between database and domain

// UUID conversions
func UUIDToPgtype(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: id,
		Valid: true,
	}
}

func PgtypeToUUID(id pgtype.UUID) (uuid.UUID, error) {
	if !id.Valid {
		return uuid.Nil, fmt.Errorf("UUID is not valid")
	}
	return id.Bytes, nil
}

// Int32 pointer conversions for nullable fields
func Int32PtrToUint(ptr *int32) (uint, error) {
	if ptr == nil {
		return 0, fmt.Errorf("type_id cannot be null")
	}
	if *ptr < 0 {
		return 0, fmt.Errorf("type_id cannot be negative: %d", *ptr)
	}
	return uint(*ptr), nil
}

func UintToInt32Ptr(arg uint) *int32 {
	val := int32(arg)
	return &val
}
