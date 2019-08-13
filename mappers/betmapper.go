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

func (mapper *BetMapper) GetFromBetIdsWithUserId(betIds []int64, userId int64) (map[int64]models.Bet, error) {
	rows, err := mapper.db.Query(`
	SELECT b.*, p.*, u2p.*, bc.* FROM bets b
	INNER JOIN positions p ON p.bet_id = b.bet_id
	INNER JOIN users_to_positions u2p ON u2p.position_id = p.position_id
	INNER JOIN bet_closes bc ON bc.bet_id = b.bet_id
	WHERE b.bet_id IN (?)`, StringFromIntSlice(betIds))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	bets := make(map[int64]models.Bet)
	var bet models.Bet
	var position *models.Position
	var winningPositionId int64

	for rows.Next() {
		position = &models.Position{}
		err = rows.Scan(&bet.BetID, &bet.Created, &bet.Closed, &position.PositionID, &position.Description, &position.OddsMultiplier, &position.UserID, &winningPositionId)
		if err != nil {
			return nil, err
		}

		if winningPositionId == 0 {
			position.State = 0
		} else if position.PositionID == winningPositionId {
			position.State = 1
		} else {
			position.State = 2
		}
		if foundBet, ok := bets[bet.BetID]; ok {
			bet = foundBet
		}

		if position.UserID == userId {
			bet.MyPosition = position
		} else {
			bet.OtherPositions = append(bet.OtherPositions, position)
		}

		prettyPrint(bet)

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

func (mapper *BetMapper) GetFromUserId(userId int64) {

}
