package model

import (
	"github.com/hashicorp/go-memdb"
)

type User struct {
	ID string
}

type Bet struct {
	ID int64
}

type Odds struct {
	ID     string
	BetID  int64
	Odds   string
	UserID string
}

func CreateSchema() *memdb.DBSchema {
	return &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"user": {
				Name: "user",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID"},
					},
				},
			},
			"bet": {
				Name: "bet",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.IntFieldIndex{Field: "ID"},
					},
				},
			},
			"odds": {
				Name: "odds",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID"},
					},
					"userId": {
						Name:    "userId",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "UserID"},
					},
					"betId": {
						Name:    "betId",
						Unique:  false,
						Indexer: &memdb.IntFieldIndex{Field: "BetID"},
					},
					"odds": {
						Name:    "odds",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Odds"},
					},
				},
			},
		},
	}
}
