package mappers

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/airnomadsmitty/quarters/models"
)

type BetMapper struct {
	db *sql.DB
}

func NewBetMapper(db *sql.DB) *BetMapper {
	return &BetMapper{db}
}

func (mapper *BetMapper) getFromBetIdsWithUserId(betIds []int64, userId int64) (map[int64]models.Bet, error) {
	rows, err := mapper.db.Query(`
	SELECT b.bet_id, b.created, b.closed, p.position_id, p.description, p.odds_multiplier, u2p.user_id, bc.winning_position_id FROM bets b
	INNER JOIN positions p ON p.bet_id = b.bet_id
	INNER JOIN users_to_positions u2p ON u2p.position_id = p.position_id
	LEFT JOIN bet_closes bc ON bc.bet_id = b.bet_id
	WHERE b.bet_id IN (?)`, StringFromIntSlice(betIds))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	bets := make(map[int64]models.Bet)
	var bet models.Bet
	var position *models.Position
	var winningPositionId *int64

	for rows.Next() {
		position = &models.Position{}
		err = rows.Scan(&bet.BetID, &bet.Created, &bet.Closed, &position.PositionID, &position.Description, &position.OddsMultiplier, &position.UserID, &winningPositionId)
		if err != nil {
			return nil, err
		}
		if winningPositionId == nil {
			position.State = models.POSITION_STATE_OPEN
		} else if position.PositionID == *winningPositionId {
			position.State = models.POSITION_STATE_WIN
		} else {
			position.State = models.POSITION_STATE_LOSE
		}
		if foundBet, ok := bets[bet.BetID]; ok {
			bet = foundBet
		}
		if position.UserID == userId {
			bet.MyPosition = position
		} else {
			bet.OtherPositions = append(bet.OtherPositions, position)
		}

		bets[bet.BetID] = bet
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return bets, nil
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func StringFromIntSlice(intSlice []int64) string {
	stringArray := make([]string, len(intSlice))
	for i, v := range intSlice {
		stringArray[i] = strconv.FormatInt(v, 10)
	}

	return strings.Join(stringArray, ",")
}

func (mapper *BetMapper) GetFromUserId(userId int64) (map[int64]models.Bet, error) {
	rows, err := mapper.db.Query(`
	SELECT b.bet_id FROM bets b
	INNER JOIN positions p ON p.bet_id = b.bet_id
	INNER JOIN users_to_positions u2p ON u2p.position_id = p.position_id
	WHERE u2p.user_id = ?`, userId)

	defer rows.Close()
	var betIDs []int64
	var betID int64

	for rows.Next() {
		err = rows.Scan(&betID)
		if err != nil {
			return nil, err
		}

		betIDs = append(betIDs, betID)
	}

	return mapper.getFromBetIdsWithUserId(betIDs, userId)
}
