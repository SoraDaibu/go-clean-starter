package common

import (
	"database/sql"
	"fmt"
)

// Helper functions for type conversion between database and domain
func SqlNullInt32ToUint(nullInt32 sql.NullInt32) (uint, error) {
	if !nullInt32.Valid {
		return 0, fmt.Errorf("type_id cannot be null")
	}
	if nullInt32.Int32 < 0 {
		return 0, fmt.Errorf("type_id cannot be negative: %d", nullInt32.Int32)
	}
	return uint(nullInt32.Int32), nil
}

func UintToSqlNullInt32(arg uint) sql.NullInt32 {
	return sql.NullInt32{Int32: int32(arg), Valid: true}
}
