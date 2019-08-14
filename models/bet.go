package models

const (
	POSITION_STATE_OPEN = 0
	POSITION_STATE_WIN  = 1
	POSITION_STATE_LOSE = 2
)

type Position struct {
	PositionID     int64
	UserID         int64
	Description    string
	OddsMultiplier float64
	State          int
}
type Bet struct {
	BetID          int64
	MyPosition     *Position
	OtherPositions []*Position
	BaseValue      float64
	Created        string
	Notes          string
	Closed         bool
	CloseDate      string
}
