package db

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/newtoallofthis123/ranhash"
)

type Token struct {
	Id     string    `json:"id"`
	UserId string    `json:"user_id"`
	Expiry time.Time `json:"expiry"`
}

type CreateTokenRequest struct {
	UserId string `json:"user_id"`
}

func (s *Store) CreateToken(token CreateTokenRequest) (Token, error) {
	id := ranhash.Generate(16, ranhash.HashPool)

	// Set the expiry to 1 month from now
	expiry := time.Now().Add(time.Hour * 24 * 7 * 4)

	_, err := s.pq.Insert("tokens").Columns("id", "user_id", "expiry").Values(id, token.UserId, expiry).RunWith(s.db).Exec()
	if err != nil {
		return Token{}, err
	}
	return Token{
		Id:     id,
		UserId: token.UserId,
		Expiry: expiry,
	}, nil
}

func (s *Store) GetToken(id string) (Token, error) {
	var token Token
	err := s.pq.Select("id", "user_id", "expiry").From("tokens").Where(squirrel.Eq{"id": id}).RunWith(s.db).Scan(&token.Id, &token.UserId, &token.Expiry)
	if err != nil {
		return Token{}, err
	}
	return token, nil
}

func (s *Store) DeleteToken(id string) error {
	_, err := s.pq.Delete("tokens").Where(squirrel.Eq{"id": id}).RunWith(s.db).Exec()
	return err
}
