package repositories

import (
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/kmrhemant916/iam/entities"
	"github.com/kmrhemant916/iam/global"
	"gorm.io/gorm"
)

type SignupRepository interface {
	CreateRootAccount(user *entities.User, organization *entities.Organization) error
}

type signupRepository struct {
	db *gorm.DB
}

func NewSignupRepository(db *gorm.DB) *signupRepository {
	return &signupRepository {
		db: db,
	}
}

func (u *signupRepository) CreateRootAccount (user *entities.User, organization *entities.Organization) error {
	var mysqlErr *mysql.MySQLError
	tx := u.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	if err := tx.Create(organization).Error; err != nil {
		tx.Rollback()
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return global.ErrOrgExists
		}
		return fmt.Errorf("failed to create organization: %w", err)
	}
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return global.ErrUserExists
		}
		return fmt.Errorf("failed to create user: %w", err)
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
