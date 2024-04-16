package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/go-memdb"
	"github.com/repeale/fp-go"
	"net/http"
	"oddschecker/model"
	"oddschecker/payloads"
	"regexp"
	"strconv"
	"strings"
)

func isStringZero(s string) bool {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i == 0
}

func (s *OddsService) validateOdds(odds string) bool {
	if s.oddsRegex.MatchString(odds) {
		if odds == "SP" {
			return true
		}

		oddsParts := strings.Split(odds, "/")
		if isStringZero(oddsParts[0]) || isStringZero(oddsParts[1]) {
			return false
		}

		return true
	}

	return false
}

func (s *OddsService) saveOdds(userID string, betID int64, odds string) error {
	txn := s.db.Txn(true)
	defer txn.Abort()

	newOdds := &model.Odds{
		ID:     uuid.New().String(), // Do not believe in odds repeating UUID in near future :)
		Odds:   odds,
		BetID:  betID,
		UserID: userID,
	}

	if err := txn.Insert("odds", newOdds); err != nil {
		return err
	}

	txn.Commit()
	return nil
}

func (s *OddsService) getAllOddsByBetID(betID int64) ([]model.Odds, error) {
	txn := s.db.Txn(false)
	defer txn.Abort()

	it, err := txn.Get("odds", "betId", betID)
	if err != nil {
		return nil, err
	}

	var results []model.Odds
	for obj := it.Next(); obj != nil; obj = it.Next() {
		odds, ok := obj.(*model.Odds)
		if !ok {
			continue
		}
		results = append(results, *odds)
	}

	return results, nil
}

func (s *OddsService) toOddsPayload(odds model.Odds) payloads.OddsPayload {
	return payloads.OddsPayload{
		BetID:  odds.BetID,
		UserID: odds.UserID,
		Odds:   odds.Odds,
	}
}

type OddsService struct {
	db          *memdb.MemDB
	userService *UserService
	betService  *BetService
	oddsRegex   *regexp.Regexp
}

func NewOddsService(db *memdb.MemDB, userService *UserService, betService *BetService) *OddsService {
	oddsPatterns := []string{
		`^\d+/\d+$`,
		`^SP$`,
	}

	combinedPattern := strings.Join(oddsPatterns, "|")
	oddsRegex, err := regexp.Compile(combinedPattern)

	if err != nil {
		panic(err)
		return nil
	}

	return &OddsService{
		db:          db,
		userService: userService,
		betService:  betService,
		oddsRegex:   oddsRegex,
	}
}

func (s *OddsService) InsertOdds(payload payloads.OddsPayload) (int, string) {
	isOddsValid := s.validateOdds(payload.Odds)

	if !isOddsValid {
		return http.StatusBadRequest, "Invalid format of Odds!"
	}

	userID, err := s.userService.GetOrCreateUserID(payload.UserID)
	if err != nil {
		panic(err)
	}

	betID, err := s.betService.GetOrCreateBetID(payload.BetID)
	if err != nil {
		panic(err)
	}

	err = s.saveOdds(userID, betID, payload.Odds)
	if err != nil {
		panic(err)
	}

	return http.StatusCreated, "Odds have been created for bet !"
}

func (s *OddsService) GetOddsByBetID(betID int64) ([]payloads.OddsPayload, error) {
	actualBetID, err := s.betService.GetBetID(betID)
	if err != nil {
		panic(err)
	}
	if actualBetID == -1 {
		return []payloads.OddsPayload{}, errors.New(fmt.Sprintf("Bet not found for given ID %d !", betID))
	}

	odds, err := s.getAllOddsByBetID(actualBetID)
	if err != nil {
		panic(err)
	}

	return fp.Map(s.toOddsPayload)(odds), nil
}
