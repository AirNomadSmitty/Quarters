package models

const (
	POSITION_STATE_OPEN = 0
	POSITION_STATE_WIN  = 1
	POSITION_STATE_LOSE = 2
	BET_STATUS_OPEN     = 0
	BET_STATUS_CLOSED   = 1
)

type Position struct {
	PositionID     int64
	UserID         int64
	Description    string
	OddsMultiplier float64
}

type BetClose struct {
	BetCloseID        *int64
	BetID             *int64
	CloseDate         *string
	WinningPositionID *int64
}
type Bet struct {
	BetID          int64
	MyPosition     *Position
	OtherPositions []*Position
	BaseValue      float64
	Created        string
	Note           string
	Close          *BetClose
}
