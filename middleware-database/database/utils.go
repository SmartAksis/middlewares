package database

func SequenceNextVal(sequence string) int64 {
	var nextId int64
	database.Raw("SELECT nextval('dominus.\"seq_accounts\"')").Scan(&nextId)
	return nextId
}

func SequenceCurrVal(sequence string) int64 {
	var nextId int64
	database.Raw("SELECT currval('dominus.\"seq_accounts\"')").Scan(&nextId)
	return nextId
}
