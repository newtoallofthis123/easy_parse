package db

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/newtoallofthis123/ranhash"
)

type Request struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateRequestRequest struct {
	UserId string `json:"user_id"`
}

func (s *Store) CreateRequest(req CreateRequestRequest, status string) (Request, error) {
	id := ranhash.Generate(16, ranhash.HashPool)

	_, err := s.pq.Insert("requests").Columns("id", "user_id", "status", "created_at").Values(id, req.UserId, status, time.Now()).RunWith(s.db).Exec()
	if err != nil {
		return Request{}, err
	}

	return Request{
		Id:        id,
		UserId:    req.UserId,
		Status:    status,
		CreatedAt: time.Now(),
	}, nil
}

func (s *Store) GetRequest(id string) (Request, error) {
	var req Request
	err := s.pq.Select("*").From("requests").Where(squirrel.Eq{"id": id}).RunWith(s.db).Scan(&req)
	if err != nil {
		return Request{}, err
	}

	return req, nil
}

func (s *Store) GetRequestByUserId(userId string) ([]Request, error) {
	var reqs []Request
	err := s.pq.Select("*").From("requests").Where(squirrel.Eq{"user_id": userId}).RunWith(s.db).Scan(&reqs)
	if err != nil {
		return nil, err
	}

	return reqs, nil
}

func (s *Store) DeleteRequest(id string) error {
	_, err := s.pq.Delete("requests").Where(squirrel.Eq{"id": id}).RunWith(s.db).Exec()
	if err != nil {
		return err
	}

	return nil
}
