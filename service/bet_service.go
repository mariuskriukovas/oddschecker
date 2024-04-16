package service

import (
	"github.com/hashicorp/go-memdb"
	"oddschecker/model"
)

type BetService struct {
	db *memdb.MemDB
}

func NewBetService(db *memdb.MemDB) *BetService {
	return &BetService{
		db: db,
	}
}

func (s *BetService) GetBetID(betID int64) (int64, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	if raw, err := txn.First("bet", "id", betID); err != nil {
		return -1, err
	} else if raw != nil {
		return betID, nil
	} else {
		return -1, nil
	}
}

func (s *BetService) SaveBetID(betID int64) (int64, error) {
	txn := s.db.Txn(true)
	defer txn.Abort()

	newBet := &model.Bet{ID: betID}
	if err := txn.Insert("bet", newBet); err != nil {
		txn.Abort()
		return 0, err
	}
	txn.Commit()

	return betID, nil
}

func (s *BetService) GetOrCreateBetID(betID int64) (int64, error) {
	existingBetID, err := s.GetBetID(betID)
	if err != nil {
		return -1, err
	} else {
		if existingBetID != -1 {
			return existingBetID, nil
		} else {
			return s.SaveBetID(betID)
		}
	}
}
