package mappers

import (
	"database/sql"

	"github.com/airnomadsmitty/quarters/models"
)

type BetMapper struct {
	db *sql.DB
}

func NewBetMapper(db *sql.DB) *BetMapper {
	return &BetMapper{db}
}

func (mapper *BetMapper) GetFromBetIdsWithUserId(betIds []int64, userId int64) ([]models.Bet, error) {
	rows, err := mapper.db.Query(`
	SELECT b.*, p.*, u2p.*, bc.* FROM bets b
	INNER JOIN positions p ON p.bet_id = b.bet_id
	INNER JOIN users_to_positions u2p ON utp.position_id = p.position_id
	INNER JOIN bet_closes bc ON bc.bet_id = b.bet_id
	WHERE b.bet_id IN ()`, 1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bets []models.Bet

	for rows.Next() {
		bet := &models.Bet{}
		// TODO: does this work
		err = rows.Scan()
		if err != nil {
			return nil, err
		}
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return bets, nil
}

func (mapper *BetMapper) GetFromUserId(userId int64) {

}
