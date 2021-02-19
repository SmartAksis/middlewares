package relational

import (
	"fmt"
	"os"
)

func PostgresSequenceNextVal(sequence string) int64 {
	var nextId int64
	stmt:=fmt.Sprintf("SELECT nextval('%s.\"%s\"')", os.Getenv("SCHEMA_DEFAULT"), sequence)
	postgresDb.Raw(stmt).Scan(&nextId)
	return nextId
}

func PostgresSequenceCurrVal(sequence string) int64 {
	var nextId int64
	stmt:=fmt.Sprintf("SELECT currval('%s.\"%s\"')", os.Getenv("SCHEMA_DEFAULT"), sequence)
	postgresDb.Raw(stmt).Scan(&nextId)
	return nextId
}

