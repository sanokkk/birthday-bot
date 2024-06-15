package local_storage

import (
	"rutube/internal/models"
	"rutube/internal/storage"
	"slices"
	"time"
)

type LocalStorage struct {
	Users map[int64]*models.User
}

func New() *LocalStorage {
	return &LocalStorage{
		Users: make(map[int64]*models.User),
	}
}

func (s *LocalStorage) GetUserInfo(id int64) (*models.User, error) {
	if user, ok := s.Users[id]; !ok {
		return nil, storage.NotFoundErr
	} else {
		return user, nil
	}
}

func (s *LocalStorage) AddUserInfo(user *models.User) (*models.User, error) {
	if _, ok := s.Users[user.Id]; ok {
		return nil, storage.DuplicateErr
	}

	s.Users[user.Id] = user

	return user, nil
}

func (s *LocalStorage) GetUsersByDate(time time.Time) []*models.User {
	result := make([]*models.User, 0)

	_, month, day := time.Date()
	for _, user := range s.Users {
		if user.Birthday == nil {
			continue
		}

		_, userMonth, userDay := user.Birthday.Date()
		if userMonth == month && userDay == day {
			result = append(result, user)
		}
	}

	return result
}

func (s *LocalStorage) UpdateUser(id int64, user *models.User) {
	usr, err := s.GetUserInfo(id)
	if err != nil {
		return
	}

	if user.Birthday != nil {
		usr.Birthday = user.Birthday
	}

	if user.ChatIds != nil && !slices.Equal(user.ChatIds, usr.ChatIds) {
		usr.ChatIds = user.ChatIds
	}

	s.Users[id] = usr
}
