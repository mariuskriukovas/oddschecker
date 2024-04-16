package payloads

type OddsPayload struct {
	BetID  int64  `json:"betId"`
	UserID string `json:"userId"`
	Odds   string `json:"odds"`
}
