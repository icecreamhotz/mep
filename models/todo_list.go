package models

import (
	"context"
	"log"
	"mime/multipart"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/gofrs/uuid"
)

type TodoListReporer interface {
	Create(todoList TodoLists) error
	UpdateDoneStatus(todoListID uuid.UUID) (TodoLists, error)
	Update(todoList *TodoLists) error
	GetAll() ([]TodoLists, int, error)
	GetById(todoListID uuid.UUID) (TodoLists, error)
	DeleteById(todoListID uuid.UUID) error
}

type TodoListRepository struct {
	DB *pg.DB
}

type TodoLists struct {
	tableName struct{}              `pg:"todo_lists"`
	ID        uuid.UUID             `pg:"type:uuid,pk,column_name:id" json:"id"`
	Name      string                `pg:"column_name:name" form:"name" json:"name" binding:"required"`
	Done      *bool                 `pg:"column_name:done,use_zero" form:"done" json:"done" binding:"bool"`
	ImageFile *multipart.FileHeader `form:"image" pg:"-"`
	Image     string                `pg:"column_name:image" json:"image"`
	CreatedAt time.Time             `pg:"column_name:created_at,null" json:"created_at"`
	UpdatedAt time.Time             `pg:"column_name:updated_at,null" json:"updated_at"`
	DeletedAt time.Time             `pg:"column_name:deleted_at,soft_delete" json:"deleted_at"`
}

var _ pg.BeforeInsertHook = (*TodoLists)(nil)

func (t *TodoLists) BeforeInsert(ctx context.Context) (context.Context, error) {
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now()
	}
	return ctx, nil
}

var _ pg.BeforeUpdateHook = (*TodoLists)(nil)

func (t *TodoLists) BeforeUpdate(ctx context.Context) (context.Context, error) {
	t.UpdatedAt = time.Now()
	return ctx, nil
}

func NewTodolistRepository(db *pg.DB) TodoListReporer {
	return &TodoListRepository{
		DB: db,
	}
}

func (repo *TodoListRepository) GetById(id uuid.UUID) (TodoLists, error) {
	todoList := TodoLists{
		ID: id,
	}
	err := repo.DB.Select(&todoList)
	return todoList, err
}

func (repo *TodoListRepository) Create(todoList TodoLists) error {
	err := repo.DB.Insert(&todoList)
	return err
}

func (repo *TodoListRepository) UpdateDoneStatus(todoListID uuid.UUID) (TodoLists, error) {
	var done bool
	err := repo.DB.Model((*TodoLists)(nil)).
		Column("done").
		Where("id = ?", todoListID).
		Select(&done)
	if err != nil {
		return TodoLists{}, err
	}
	done = !done
	todoList := TodoLists{
		ID:   todoListID,
		Done: &done,
	}

	_, err = repo.DB.Model(&todoList).Column("done").WherePK().Returning("*").Update()

	return todoList, err
}

func (repo *TodoListRepository) GetAll() ([]TodoLists, int, error) {
	var todoLists []TodoLists
	err := repo.DB.Model(&todoLists).Select()
	if err != nil {
		log.Fatal(err)
	}
	return todoLists, len(todoLists), nil
}

func (repo *TodoListRepository) Update(todoList *TodoLists) error {
	_, err := repo.DB.Model(todoList).WherePK().Returning("*").Update()
	return err
}

func (repo *TodoListRepository) DeleteById(todoListID uuid.UUID) error {
	todoList := TodoLists{
		ID: todoListID,
	}
	err := repo.DB.Delete(&todoList)
	return err
}
