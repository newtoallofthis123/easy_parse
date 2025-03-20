package db

import "github.com/Masterminds/squirrel"

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

func (s *Store) CreateUser(user User) (User, error) {
	_, err := s.pq.Insert("users").Columns("id", "email").Values(user.Id, user.Email).RunWith(s.db).Exec()
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *Store) GetUser(id string) (User, error) {
	var user User
	err := s.pq.Select("id", "email").From("users").Where(squirrel.Eq{"id": id}).RunWith(s.db).Scan(&user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *Store) DeleteUser(id string) error {
	_, err := s.pq.Delete("users").Where(squirrel.Eq{"id": id}).RunWith(s.db).Exec()
	if err != nil {
		return err
	}

	return nil
}
