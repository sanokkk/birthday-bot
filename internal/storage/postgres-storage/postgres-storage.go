package postgres_storage

import (
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"rutube/internal/config"
	"rutube/internal/models"
	"rutube/internal/storage"
	"slices"
	"sync"
	"time"
)

type PostgresStorage struct {
	*sync.RWMutex
	db *gorm.DB
}

func New(cfg *config.Config) (*PostgresStorage, error) {
	db, err := gorm.Open(postgres.Open(cfg.ConnectionString))

	if err != nil {
		return nil, storage.OpenDbErr
	}

	return &PostgresStorage{
		db:      db,
		RWMutex: &sync.RWMutex{},
	}, nil
}

func (ps *PostgresStorage) GetUserInfo(id int64) (*models.User, error) {

	user := models.User{Id: id}
	db := ps.db.First(&user)

	if db.Error != nil || db.RowsAffected == 0 {
		if errors.Is(ps.db.Error, gorm.ErrRecordNotFound) {
			return nil, storage.NotFoundErr
		}

		return nil, storage.UnSpecifiedErr
	}

	return &user, nil
}

func (ps *PostgresStorage) GetUsersByDate(time time.Time) []models.User {
	ps.RWMutex.Lock()
	defer ps.RWMutex.Unlock()

	var users []models.User
	_, month, day := time.Date()

	db := ps.db.Table("users").Where(
		"extract('month' from users.birthday) = ? and extract('day' from users.birthday) = $2",
		int(month),
		day).Find(&users)

	if db.Error != nil {
		return []models.User{}
	}

	return users
}

func (ps *PostgresStorage) UpdateUser(id int64, user *models.User) {
	ps.RWMutex.Lock()
	defer ps.RWMutex.Unlock()

	dbUser, err := ps.GetUserInfo(id)
	if err != nil {
		if errors.Is(err, storage.NotFoundErr) {
			return
		}
	}

	if user.Birthday != nil && user.Birthday != dbUser.Birthday {
		dbUser.Birthday = user.Birthday
	}

	if user.ChatIds != nil && len(user.ChatIds) > 0 && !slices.Equal(dbUser.ChatIds, user.ChatIds) {
		dbUser.ChatIds = user.ChatIds
	}

	ps.db.Save(dbUser)
}

func (ps *PostgresStorage) AddUserInfo(user *models.User) (*models.User, error) {
	ps.RWMutex.Lock()
	defer ps.RWMutex.Unlock()

	db := ps.db.Create(user)

	if db.Error != nil {
		fmt.Println("error while creating user in db")
		if errors.Is(ps.db.Error, gorm.ErrDuplicatedKey) {
			return nil, storage.DuplicateErr
		}

		return nil, storage.UnSpecifiedErr
	}

	return user, nil
}

func (ps *PostgresStorage) Migrate() error {
	return ps.db.AutoMigrate(&models.User{})
}
