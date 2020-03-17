package models

import (
	"time"

	"github.com/go-pg/pg/v9"
)

type UserReporer interface {
	GetById(id string) (BackofficeUser, error)
	FindByUsername(username string) (BackofficeUser, error)
	FindByEmail(email string) (BackofficeUser, error)
	Create(backofficeUser BackofficeUser) error
}

type UserRepository struct {
	DB *pg.DB
}

type BackofficeUser struct {
	tableName       struct{}  `pg:"backoffice_users"`
	ID              string    `pg:"type:uuid,pk,column_name:id"`
	Username        string    `pg:"column_name:username,unique" json:"username" binding:"required"`
	Password        string    `pg:"column_name:username" json:"password" binding:"required,min=8"`
	ConfirmPassword string    `pg:"-" json:"confirm_password" binding:"required,min=8,eqfield=Password"`
	Email           string    `pg:"column_name:email,unique" json:"email" binding:"required,email"`
	Name            string    `pg:"column_name:name" json:"name" binding:"required"`
	Lastname        string    `pg:"column_name:lastname" json:"lastname" binding:"required"`
	Role            string    `pg:"column_name:role" json:"role" binding:"required,oneof=ADMIN EMPLOYEE"`
	CreatedAt       time.Time `pg:"column_name:created_at,null"`
	UpdatedAt       time.Time `pg:"column_name:updated_at,null"`
	DeletedAt       time.Time `pg:"column_name:deleted_at,soft_delete"`
}

func NewUserRepository(db *pg.DB) UserReporer {
	return &UserRepository{
		DB: db,
	}
}

func (repo *UserRepository) GetById(id string) (BackofficeUser, error) {
	backofficeUser := BackofficeUser{
		ID: id,
	}
	err := repo.DB.Select(&backofficeUser)
	return backofficeUser, err
}

func (repo *UserRepository) Create(backofficeUser BackofficeUser) error {
	err := repo.DB.Insert(&backofficeUser)
	return err
}

func (repo *UserRepository) FindByUsername(username string) (BackofficeUser, error) {
	var backofficeUser BackofficeUser
	err := repo.DB.Model(&backofficeUser).Where("username = ?", username).Select()
	return backofficeUser, err
}

func (repo *UserRepository) FindByEmail(email string) (BackofficeUser, error) {
	var backofficeUser BackofficeUser
	err := repo.DB.Model(&backofficeUser).Where("email = ?", email).Select()
	return backofficeUser, err
}
