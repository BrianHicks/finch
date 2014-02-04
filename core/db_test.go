package core

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/syndtr/goleveldb/leveldb/storage"
	"os"

	"testing"
	"time"
)

// TestTaskStoreDiskLifecycle makes sure that we can read and write to disk. The
// rest of the TaskStore tests should use the in-memory versions of the database.
func TestTaskStoreDiskLifecycle(t *testing.T) {
	name := "_taskdb_lifecycle"
	filestore, err := storage.OpenFile(name)
	assert.Nil(t, err)
	defer filestore.Close()

	_, err = NewTaskStore(filestore)
	assert.Nil(t, err)
	defer os.RemoveAll(name)

	assert.Nil(t, err)
}

type TaskStoreSuite struct {
	suite.Suite
	Store *TaskStore
}

func (suite *TaskStoreSuite) SetupTest() {
	Store, err := NewTaskStore(storage.NewMemStorage())
	assert.Nil(suite.T(), err)
	suite.Store = Store
}

func (suite *TaskStoreSuite) TearDownTest() {
	suite.Store.Store.Close()
}

func (suite *TaskStoreSuite) TestTasksIndexing() {
	t := NewTask("test", time.Now())

	err := suite.Store.PutTasks(t)
	assert.Nil(suite.T(), err)

	task, err := suite.Store.GetTask(t.Key())
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), t, task)
}

func (suite *TaskStoreSuite) TestMove() {
	task := NewTask("test", time.Now())
	key := task.Key()

	err := suite.Store.PutTasks(task)
	if !assert.Nil(suite.T(), err) {
		// the rest of this test wouldn't make sense now, abort!
		return
	}

	task.ID = "some-other-value"
	suite.Store.MoveTask(key, task)

	_, err = suite.Store.GetTask(key)
	assert.Equal(suite.T(), ErrNoTask, err)

	present, err := suite.Store.GetTask(task.Key())
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), task, present)
}

func (suite *TaskStoreSuite) TestPendingIndexing() {
	nope := NewTask("test", time.Now())
	nope.Attrs[TagPending] = false

	yep := NewTask("test", time.Now())
	nope.Attrs[TagPending] = true

	suite.Store.PutTasks(nope, yep)

	pending, err := suite.Store.GetPendingTasks()
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), []*Task{yep}, pending)
}

func (suite *TaskStoreSuite) TestSelectedIndexing() {
	nope := NewTask("test", time.Now())
	nope.Attrs[TagSelected] = false

	yep := NewTask("selected", time.Now())
	yep.Attrs[TagSelected] = true

	suite.Store.PutTasks(nope, yep)

	selected, err := suite.Store.GetSelectedTasks()
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), []*Task{yep}, selected)
}

func (suite *TaskStoreSuite) TestGetNextSelected() {
	// try it empty first
	_, err := suite.Store.GetNextSelected()
	assert.Equal(suite.T(), ErrNoTask, err)

	t := NewTask("test", time.Now())
	t.Attrs[TagSelected] = true

	suite.Store.PutTasks(t)

	next, err := suite.Store.GetNextSelected()
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), t, next)
}

func TestTaskStoreSuite(t *testing.T) {
	suite.Run(t, new(TaskStoreSuite))
}
