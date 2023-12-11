package xgor

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TestEntity struct {
	gorm.Model
	Name string
	Age  int
	Ulid string
}

func TestGormBaseRepository_SHouldAddEntity(t *testing.T) {
	// Connect to an in-memory SQLite database
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Migrate the database schema
	err = db.AutoMigrate(&TestEntity{})
	assert.NoError(t, err)

	// HandleCreate a new repository for the TestEntity
	repo := New[TestEntity](db, errors.New("entity not found"))

	// Test adding an entity
	entity := &TestEntity{Name: "test"}
	err = repo.Add(entity)
	assert.Equal(t, "test", entity.Name)
	assert.Equal(t, uint(1), entity.ID)
	assert.NoError(t, err)
	_ = repo.Delete(entity)
}

func TestGormBaseRepository_ShouldGetByID(t *testing.T) {
	// Connect to an in-memory SQLite database
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Migrate the database schema
	err = db.AutoMigrate(&TestEntity{})
	assert.NoError(t, err)

	// HandleCreate a new repository for the TestEntity
	repo := New[TestEntity](db, errors.New("entity not found"))

	// Test getting an entity by ID
	entity := &TestEntity{Name: "test"}
	err = repo.Add(entity)
	assert.NoError(t, err)

	result, err := repo.GetByID(entity.ID)
	assert.Equal(t, "test", result.Name)
	assert.NoError(t, err)
	_ = repo.Delete(entity)

	//// Test getting a non-existent entity by ID
	_, err = repo.GetByID(0)
	assert.Error(t, err)
	//assert.Equal(t, "No entity with id 0 found", err.(*EntityNotFoundError).Message)
}

func TestGormBaseRepository_ShouldUpdate(t *testing.T) {
	// Connect to an in-memory SQLite database
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Migrate the database schema
	err = db.AutoMigrate(&TestEntity{})
	assert.NoError(t, err)

	// HandleCreate a new repository for the TestEntity
	repo := New[TestEntity](db, errors.New("entity not found"))

	// Test updating an entity
	entity := &TestEntity{Name: "test"}
	err = repo.Add(entity)
	assert.NoError(t, err)

	entity.Name = "updated"
	err = repo.Update(entity)
	assert.NoError(t, err)
	_ = repo.Delete(entity)
}

func TestGormBaseRepository_ShouldDelete(t *testing.T) {
	// Connect to an in-memory SQLite database
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Migrate the database schema
	err = db.AutoMigrate(&TestEntity{})
	assert.NoError(t, err)

	// HandleCreate a new repository for the TestEntity
	repo := New[TestEntity](db, errors.New("entity not found"))

	// Test deleting an entity
	entity := &TestEntity{Name: "test"}
	err = repo.Add(entity)
	assert.NoError(t, err)

	err = repo.Delete(entity)
	assert.NoError(t, err)
}

func TestGormBaseRepository_ShouldGetAll(t *testing.T) {
	// Connect to an in-memory SQLite database
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Migrate the database schema
	err = db.AutoMigrate(&TestEntity{})
	assert.NoError(t, err)

	// HandleCreate a new repository for the TestEntity
	repo := New[TestEntity](db, errors.New("entity not found"))

	// Test getting all entities
	entity1 := &TestEntity{Name: "test1"}
	err = repo.Add(entity1)
	assert.NoError(t, err)

	entity2 := &TestEntity{Name: "test2"}
	err = repo.Add(entity2)
	assert.NoError(t, err)
	//entities := make([]*TestEntity, 0)
	entities, err := repo.GetAll(nil, nil, nil, nil)
	result := *entities.Items
	assert.NoError(t, err)
	assert.Equal(t, int64(2), entities.ResultCount)
	assert.Equal(t, int64(2), entities.TotalCount)
	assert.Equal(t, result[0].Name, "test1")
	assert.Equal(t, result[1].Name, "test2")

	//
	// Test pagination
	entities, err = repo.GetAll(nil, nil, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), entities.TotalCount)
	assert.Equal(t, int64(2), entities.ResultCount)

	limit := 1
	offset := 0
	entities, err = repo.GetAll(&limit, &offset, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), entities.TotalCount)
	assert.Equal(t, int64(1), entities.ResultCount)

	// Test ordering
	entity3 := &TestEntity{Name: "test3"}
	err = repo.Add(entity3)
	assert.NoError(t, err)

	orderBy := "name DESC"
	entities, err = repo.GetAll(nil, nil, &orderBy, nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), entities.TotalCount)
	assert.Equal(t, int64(3), entities.ResultCount)
	result = *entities.Items
	assert.Equal(t, result[0].Name, "test3")
	assert.Equal(t, result[1].Name, "test2")
	assert.Equal(t, result[2].Name, "test1")

	orderBy = "name ASC"
	entities, err = repo.GetAll(nil, nil, &orderBy, nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), entities.TotalCount)
	assert.Equal(t, int64(3), entities.ResultCount)
	result = *entities.Items
	assert.Equal(t, result[0].Name, "test1")
	assert.Equal(t, result[1].Name, "test2")
	assert.Equal(t, result[2].Name, "test3")

	// Test filtering
	filters := map[string]interface{}{"name": "test2"}
	entities, err = repo.GetAll(nil, nil, nil, filters)
	result = *entities.Items
	assert.NoError(t, err)
	assert.Equal(t, int64(1), entities.TotalCount)
	assert.Equal(t, int64(1), entities.ResultCount)
	assert.Equal(t, result[0].Name, "test2")
}

func TestApplyFilters(t *testing.T) {

	t.Run("Filter by Greater Than", func(t *testing.T) {
		db, err := Open(sqlite.Open("file::memory:"), &gorm.Config{})
		assert.NoError(t, err)

		err = db.AutoMigrate(&TestEntity{})
		assert.NoError(t, err)

		repo := New[TestEntity](db, errors.New("entity not found"))
		entity1 := &TestEntity{Name: "ismael", Age: 10}
		err = repo.Add(entity1)
		assert.NoError(t, err)

		entity2 := &TestEntity{Name: "grahms", Age: 20}
		err = repo.Add(entity2)
		assert.NoError(t, err)

		filters := map[string]interface{}{"age__gt": "10"}
		entities, err := repo.GetAll(nil, nil, nil, filters)
		result := *entities.Items
		assert.NoError(t, err)
		assert.Equal(t, int64(1), entities.TotalCount)
		assert.Equal(t, int64(1), entities.ResultCount)
		assert.Equal(t, result[0].Name, "grahms")
	})

	t.Run("Filter by Greater Than or Equal", func(t *testing.T) {
		db, err := Open(sqlite.Open("file::memory:"), &gorm.Config{})
		assert.NoError(t, err)

		err = db.AutoMigrate(&TestEntity{})
		assert.NoError(t, err)

		repo := New[TestEntity](db, errors.New("entity not found"))
		// Test filtering by range
		entity1 := &TestEntity{Name: "ismael", Age: 10}
		err = repo.Add(entity1)
		assert.NoError(t, err)

		entity2 := &TestEntity{Name: "grahms", Age: 20}
		err = repo.Add(entity2)
		assert.NoError(t, err)

		filters := map[string]interface{}{"age__gte": "10"}
		entities, err := repo.GetAll(nil, nil, nil, filters)
		result := *entities.Items
		assert.NoError(t, err)
		assert.Equal(t, int64(2), entities.TotalCount)
		assert.Equal(t, int64(2), entities.ResultCount)
		assert.Equal(t, result[0].Name, "ismael")
		assert.Equal(t, result[1].Name, "grahms")
	})

	// Continue with other filter conditions...
}

func TestShouldReturnErrorIfEntityNotFound(t *testing.T) {
	// Connect to an in-memory SQLite database
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Migrate the database schema
	err = db.AutoMigrate(&TestEntity{})
	assert.NoError(t, err)

	// HandleCreate a new repository for the TestEntity
	repo := New[TestEntity](db, errors.New("entity not found"))

	// Test getting all entities
	entity1 := &TestEntity{Name: "ismael"}
	err = repo.Add(entity1)
	assert.NoError(t, err)

	entity2 := &TestEntity{Name: "grahms"}
	err = repo.Add(entity2)
	assert.NoError(t, err)

	// Test getting an entity by id

	id := 3
	entity, err := repo.GetByID(id)
	assert.Error(t, err)
	assert.Nil(t, entity)
	errorMessage := &EntityNotFoundError{Message: "No entity with id " + string(rune(id)) + " found"}

	assert.EqualError(t, err, errorMessage.Error())

}

func TestShouldReturnErrorIfListIsEmpty(t *testing.T) {
	// Connect to an in-memory SQLite database
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Migrate the database schema
	err = db.AutoMigrate(&TestEntity{})
	assert.NoError(t, err)

	// HandleCreate a new repository for the TestEntity
	repo := New[TestEntity](db, errors.New("entity not found"))

	// Test getting all entities
	entities, err := repo.GetAll(nil, nil, nil, nil)

	assert.Equal(t, int64(0), entities.TotalCount)
	assert.Equal(t, int64(0), entities.ResultCount)
}

func TestShouldGetByULID(t *testing.T) {
	// Connect to an in-memory SQLite database
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Migrate the database schema
	err = db.AutoMigrate(&TestEntity{})
	assert.NoError(t, err)

	// HandleCreate a new repository for the TestEntity
	repo := New[TestEntity](db, errors.New("entity not found"))

	entity := &TestEntity{Name: "ismael", Age: 10, Ulid: "123"}
	err = repo.Add(entity)
	assert.NoError(t, err)
	result, err := repo.GetWithCustomFilters(FilterType{"ulid": "123"})
	assert.NoError(t, err)
	assert.Equal(t, result.Name, "ismael")

	//	not found
	_, err = repo.GetWithCustomFilters(FilterType{"ulid": "nothing"})
	assert.NotNil(t, err)

}

func TestGormBaseRepository_ShouldHandleTransaction(t *testing.T) {
	// Connect to an in-memory SQLite database
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Migrate the schema
	err = db.AutoMigrate(&TestEntity{})
	assert.NoError(t, err)

	// HandleCreate a new repository for the TestEntity
	repo := New[TestEntity](db, errors.New("entity not found"))

	err = repo.PerformTransaction(func(tx *gorm.DB) error {
		repoWithTx := New[TestEntity](tx, errors.New("entity not found"))

		entity := &TestEntity{Name: "InTransaction"}
		err := repoWithTx.Add(entity)
		assert.NoError(t, err)

		// Let's simulate an error to see if it rolls back correctly
		return errors.New("simulated error")
	})

	assert.Error(t, err)

	// Ensure the entity was not added due to rollback
	result, _ := repo.GetByID(1)
	assert.Nil(t, result)
}
