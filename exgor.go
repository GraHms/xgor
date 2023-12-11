package xgor

import (
	"errors"
	"gorm.io/gorm"
	"strings"
)

type filterType map[string]any
type ListItems[T any] struct {
	Items       *[]T
	TotalCount  int64
	ResultCount int64
}

type Repository[T any] interface {
	Add(item *T) error
	Update(item *T) error
	Delete(item *T) error
	GetByID(id interface{}) (*T, error)
	GetWithCustomFilters(filters filterType) (*T, error)
	GetAll(limit *int, offset *int, orderBy *string, filters filterType) (ListItems[T], error)
	PerformTransaction(fn func(tx *gorm.DB) error) error
	DeleteRelationship(item *T, relationship string) error
}

type EntityNotFoundError struct {
	Message string
}

func (e *EntityNotFoundError) Error() string {
	return e.Message
}

type BaseRepository[T any] struct {
	db            *gorm.DB
	notFoundError error
	relationships []string
}

func New[T any](db *gorm.DB, notFoundError error) *BaseRepository[T] {
	return &BaseRepository[T]{
		db:            db,
		notFoundError: notFoundError,
		relationships: []string{},
	}
}

func NewWithRelationships[T any](db *gorm.DB, notFoundError error, relationships ...string) *BaseRepository[T] {
	return &BaseRepository[T]{
		db:            db,
		notFoundError: notFoundError,
		relationships: relationships,
	}
}

func (r *BaseRepository[T]) Add(item *T) error {
	return r.db.Create(item).Error
}

func (r *BaseRepository[T]) Update(item *T) error {
	return r.db.Save(item).Error
}

func (r *BaseRepository[T]) Delete(item *T) error {
	return r.db.Delete(item).Error
}

func (r *BaseRepository[T]) GetWithCustomFilters(filters filterType) (*T, error) {
	limit := 1
	entity, err := r.GetAll(&limit, nil, nil, filters)
	if err != nil || entity.TotalCount == 0 {
		return nil, &EntityNotFoundError{Message: "No entity found"}
	}

	return &(*entity.Items)[0], nil
}

func (r *BaseRepository[T]) GetByID(id interface{}) (*T, error) {
	entity := new(T)

	result := r.db.First(entity, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &EntityNotFoundError{Message: "No entity with id " + string(rune(id.(int))) + " found"}
		}
	}
	return entity, nil
}

func (r *BaseRepository[T]) GetAll(limit *int, offset *int, orderBy *string, filters filterType) (ListItems[T], error) {
	var total int64

	countQuery := r.preload().Model(new(T)).Scopes(r.applyFilters(filters))
	_ = countQuery.Count(&total)

	query := r.preload().Model(new(T)).Scopes(r.paginate(limit, offset), r.applyFilters(filters), r.orderBy(orderBy))

	var results []T
	result := query.Find(&results)
	if result.Error != nil {
		return ListItems[T]{}, result.Error
	}

	if len(results) == 0 {
		return ListItems[T]{}, &EntityNotFoundError{Message: "No entities found"}
	}

	return ListItems[T]{Items: &results, TotalCount: total, ResultCount: int64(len(results))}, nil
}

func (r *BaseRepository[T]) PerformTransaction(fn func(tx *gorm.DB) error) error {

	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *BaseRepository[T]) preload() *gorm.DB {
	db := r.db
	if len(r.relationships) > 0 {
		for _, relationship := range r.relationships {
			db = db.Preload(relationship)
		}
	}
	return db
}

func (r *BaseRepository[T]) DeleteRelationship(item *T, relationship string) error {
	return r.db.Model(item).Association(relationship).Clear()
}

func (r *BaseRepository[T]) applyFilters(filters filterType) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for c, v := range filters {
			col, val := c, v
			relationshipPath := strings.Split(col, ".")
			for i := 0; i < len(relationshipPath)-1; i++ {
				db = db.Joins(relationshipPath[i])
			}
			finalFieldWithCondition := relationshipPath[len(relationshipPath)-1]
			col = finalFieldWithCondition
			switch {
			case strings.HasSuffix(col, "__eq"):
				db = db.Where(strings.TrimSuffix(col, "__eq")+" = ?", val)
			case strings.HasSuffix(col, "__gt"):
				db = db.Where(strings.TrimSuffix(col, "__gt")+" > ?", val)
			case strings.HasSuffix(col, "__lt"):
				db = db.Where(strings.TrimSuffix(col, "__lt")+" < ?", val)
			case strings.HasSuffix(col, "__gte"):
				db = db.Where(strings.TrimSuffix(col, "__gte")+" >= ?", val)
			case strings.HasSuffix(col, "__lte"):
				db = db.Where(strings.TrimSuffix(col, "__lte")+" <= ?", val)
			case strings.HasSuffix(col, "__in"):
				db = db.Where(strings.TrimSuffix(col, "__in")+" IN ?", val)
			case strings.HasSuffix(col, "__not"):
				db = db.Where(strings.TrimSuffix(col, "__not")+" <> ?", val)
			case strings.HasSuffix(col, "__not_in"):
				db = db.Where(strings.TrimSuffix(col, "__not_in")+" NOT IN ?", val)
			case strings.HasSuffix(col, "__like"):
				db = db.Where(strings.TrimSuffix(col, "__like")+" LIKE ?", val)
			default:
				db = db.Where(col+" = ?", val)
			}
		}
		return db
	}
}

func (r *BaseRepository[T]) paginate(limit *int, offset *int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if limit != nil {
			db = db.Limit(*limit)
		}
		if offset != nil {
			db = db.Offset(*offset)
		}
		return db
	}
}

func (r *BaseRepository[T]) orderBy(orderBy *string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if orderBy != nil {
			db = db.Order(*orderBy)
		}
		return db
	}
}
