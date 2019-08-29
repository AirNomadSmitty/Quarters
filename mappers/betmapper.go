package mappers

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/airnomadsmitty/quarters/models"
)

type BetMapper struct {
	db *sql.DB
}

func NewBetMapper(db *sql.DB) *BetMapper {
	return &BetMapper{db}
}

func (mapper *BetMapper) getFromBetIdsWithUserId(betIds []int64, userId int64) (map[int64]models.Bet, error) {
	// Don't mind using a placeholder since this is internal and will always be safe
	sql := fmt.Sprintf(`
	SELECT b.bet_id, b.created, p.position_id, p.description, p.odds_multiplier, u2p.user_id, bc.bet_close_id, bc.bet_id, bc.close_date, bc.winning_position_id FROM bets b
	INNER JOIN positions p ON p.bet_id = b.bet_id
	INNER JOIN users_to_positions u2p ON u2p.position_id = p.position_id
	LEFT JOIN bet_closes bc ON bc.bet_id = b.bet_id
	WHERE b.bet_id IN (%s)`, StringFromIntSlice(betIds))

	rows, err := mapper.db.Query(sql)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	bets := make(map[int64]models.Bet)
	var (
		bet      models.Bet
		position *models.Position
		close    *models.BetClose
	)

	for rows.Next() {
		position = &models.Position{}
		close = &models.BetClose{}
		err = rows.Scan(&bet.BetID, &bet.Created, &position.PositionID, &position.Description, &position.OddsMultiplier, &position.UserID, &close.BetCloseID, &close.BetID, &close.CloseDate, &close.WinningPositionID)
		if err != nil {
			return nil, err
		}
		if foundBet, ok := bets[bet.BetID]; ok {
			bet = foundBet
		}
		if position.UserID == userId {
			bet.MyPosition = position
		} else {
			bet.OtherPositions = append(bet.OtherPositions, position)
		}
		if close.BetID != nil {
			bet.Close = close
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

func (mapper *BetMapper) GetFromBetIdWithUserId(betID int64, userID int64) (*models.Bet, error) {
	bets, err := mapper.getFromBetIdsWithUserId([]int64{betID}, userID)
	if err != nil {
		return nil, err
	}

	bet := bets[betID]
	return &bet, nil
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

func (mapper *BetMapper) GetClosedFromUserId(userID int64) (map[int64]models.Bet, error) {
	rows, err := mapper.db.Query(`
	SELECT b.bet_id FROM bets b
	INNER JOIN positions p ON p.bet_id = b.bet_id
	INNER JOIN users_to_positions u2p ON u2p.position_id = p.position_id
	INNER JOIN bet_closes bc ON b.bet_id = bc.bet_id
	WHERE u2p.user_id = ?`, userID)

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

	return mapper.getFromBetIdsWithUserId(betIDs, userID)
}

func (mapper *BetMapper) GetOpenFromUserId(userID int64) (map[int64]models.Bet, error) {
	rows, err := mapper.db.Query(`
	SELECT b.bet_id FROM bets b
	INNER JOIN positions p ON p.bet_id = b.bet_id
	INNER JOIN users_to_positions u2p ON u2p.position_id = p.position_id
	LEFT JOIN bet_closes bc ON b.bet_id = bc.bet_id
	WHERE u2p.user_id = ? AND bc.bet_id is NULL`, userID)

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

	return mapper.getFromBetIdsWithUserId(betIDs, userID)
}

func StringFromIntSlice(intSlice []int64) string {
	stringArray := make([]string, len(intSlice))
	for i, v := range intSlice {
		stringArray[i] = strconv.FormatInt(v, 10)
	}

	return strings.Join(stringArray, ",")
}

func (mapper *BetMapper) Create(baseValue float64, positionDescription string, oddsMultiplier float64, userID int64, note string) (*models.Bet, error) {
	now := time.Now().UTC().Format("2006-01-02 15:04:05")

	betResult, err := mapper.db.Exec("INSERT INTO bets (base_value, note) VALUES (?, ?)", baseValue, note)
	if err != nil {
		return nil, err
	}

	betID, err := betResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	positionResult, err := mapper.db.Exec("INSERT INTO positions (description, odds_multiplier, bet_id) VALUES (?,?,?)", positionDescription, oddsMultiplier, betID)
	if err != nil {
		return nil, err
	}

	positionID, err := positionResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	_, err = mapper.db.Exec("INSERT INTO users_to_positions (user_id, position_id) VALUES (?,?)", userID, positionID)
	if err != nil {
		return nil, err
	}

	var otherPositions []*models.Position
	position := &models.Position{positionID, userID, positionDescription, oddsMultiplier}
	// TODO: timestamp stuff
	bet := &models.Bet{betID, position, otherPositions, baseValue, now, note, nil}

	return bet, nil
}
