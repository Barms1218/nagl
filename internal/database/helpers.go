package database

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func UUIDToPgtype(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

func PgTypeToUUID(id pgtype.UUID) uuid.UUID {
	var gUUID uuid.UUID
	if id.Valid {
		gUUID = id.Bytes
	}

	return gUUID
}

func IntToPgtype(n int32) pgtype.Int4 {
	return pgtype.Int4{Int32: n, Valid: true}
}

func StringToPgtype(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: true}
}
