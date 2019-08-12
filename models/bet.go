package models

type Position struct {
	PositionID     int64
	UserID         int64
	Description    string
	OddsMultiplier float64
}
type Bet struct {
	MyPosition     Position
	OtherPositions []Position
	BaseValue      float64
	Created        string
	Notes          string
	Closed         bool
	CloseDate      string
}