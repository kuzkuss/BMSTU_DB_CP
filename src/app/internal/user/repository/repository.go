package repository

import (
	"github.com/kuzkuss/VK_DB_Project/app/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RepositoryI interface {
	SelectUserByNickName(nickname string) (*models.User, error)
	SelectUserByEmail(email string) (*models.User, error)
	SelectUsersByNickNameOrEmail(nickname string, email string) ([]*models.User, error)
	CreateUser(user *models.User) (error)
	UpdateUser(user *models.User) (error)
}

type dataBase struct {
	db *gorm.DB
}

func New(db *gorm.DB) RepositoryI {
	return &dataBase{
		db: db,
	}
}

func (dbUser *dataBase) CreateUser(user *models.User) (error) {
	tx := dbUser.db.Create(user)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table users)")
	}

	return nil
}

func (dbUser *dataBase) SelectUserByNickName(nickname string) (*models.User, error) {
	user := models.User{}

	tx := dbUser.db.Where("nickname = ?", nickname).Take(&user)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	} else if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table users)")
	}

	return &user, nil
}

func (dbUser *dataBase) SelectUsersByNickNameOrEmail(nickname string, email string) ([]*models.User, error) {
	users := make([]*models.User, 0, 10)

	tx := dbUser.db.Where("email = ? OR nickname = ?", email, nickname).Find(&users)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table users)")
	}

	return users, nil
}

func (dbUser *dataBase) SelectUserByEmail(email string) (*models.User, error) {
	user := models.User{}

	tx := dbUser.db.Where("email = ?", email).Take(&user)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	} else if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table users)")
	}

	return &user, nil
}

func (dbUser *dataBase) UpdateUser(user *models.User) (error) {
	tx := dbUser.db.Model(user).Clauses(clause.Returning{}).Updates(models.User{About:user.About, Email: user.Email, FullName: user.FullName})
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table users)")
	}

	return nil
}

