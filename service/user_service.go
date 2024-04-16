package service

import (
	"github.com/hashicorp/go-memdb"
	"oddschecker/model"
)

type UserService struct {
	db *memdb.MemDB
}

func NewUserService(db *memdb.MemDB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) getUserID(userID string) (string, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	if raw, err := txn.First("user", "id", userID); err != nil {
		return "", err
	} else if raw != nil {
		return userID, nil
	} else {
		return "", nil
	}
}

func (s *UserService) saveUserID(userID string) (string, error) {
	txn := s.db.Txn(true)
	defer txn.Abort()

	newUser := &model.User{ID: userID}
	if err := txn.Insert("user", newUser); err != nil {
		return "", err
	}

	txn.Commit()
	return userID, nil
}

func (s *UserService) GetOrCreateUserID(userID string) (string, error) {
	existingUserID, err := s.getUserID(userID)
	if err != nil {
		return "", err
	} else {
		if existingUserID != "" {
			return existingUserID, nil
		} else {
			return s.saveUserID(userID)
		}
	}
}
