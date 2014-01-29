package finch

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/syndtr/goleveldb/leveldb/storage"
	"os"

	"testing"
	"time"
)

// TestTaskDBDiskLifecycle makes sure that we can read and write to disk. The
// rest of the TaskDB tests should use the in-memory versions of the database.
func TestTaskDBDiskLifecycle(t *testing.T) {
	name := "_taskdb_lifecycle"
	filestore, err := storage.OpenFile(name)
	assert.Nil(t, err)
	defer filestore.Close()

	DB, err := NewTaskDB(filestore)
	assert.Nil(t, err)
	defer os.RemoveAll(name)

	assert.Nil(t, err)

	DB.Close()

	assert.Nil(t, DB.DB)
}

type TaskDBSuite struct {
	suite.Suite
	DB *TaskDB
}

func (suite *TaskDBSuite) SetupTest() {
	DB, err := NewTaskDB(storage.NewMemStorage())
	assert.Nil(suite.T(), err)
	suite.DB = DB
}

func (suite *TaskDBSuite) TearDownTest() {
	suite.DB.Close()
}

func (suite *TaskDBSuite) TestTasksIndexing() {
	t := NewTask("test", time.Now())

	err := suite.DB.PutTasks(t)
	assert.Nil(suite.T(), err)

	task, err := suite.DB.GetTask(t.Key())
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), t, task)
}

func (suite *TaskDBSuite) TestMove() {
	task := NewTask("test", time.Now())
	key := task.Key()

	err := suite.DB.PutTasks(task)
	if !assert.Nil(suite.T(), err) {
		// the rest of this test wouldn't make sense now, abort!
		return
	}

	task.ID = "some-other-value"
	suite.DB.MoveTask(key, task)

	_, err = suite.DB.GetTask(key)
	assert.Equal(suite.T(), ErrNoTask, err)

	present, err := suite.DB.GetTask(task.Key())
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), task, present)
}

func (suite *TaskDBSuite) TestPendingIndexing() {
	nope := NewTask("test", time.Now())
	nope.Attrs[TagPending] = false

	yep := NewTask("test", time.Now())
	nope.Attrs[TagPending] = true

	suite.DB.PutTasks(nope, yep)

	pending, err := suite.DB.GetPendingTasks()
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), []*Task{yep}, pending)
}

func (suite *TaskDBSuite) TestSelectedIndexing() {
	nope := NewTask("test", time.Now())
	nope.Attrs[TagSelected] = false

	yep := NewTask("selected", time.Now())
	yep.Attrs[TagSelected] = true

	suite.DB.PutTasks(nope, yep)

	selected, err := suite.DB.GetSelectedTasks()
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), []*Task{yep}, selected)
}

func (suite *TaskDBSuite) TestGetNextSelected() {
	// try it empty first
	_, err := suite.DB.GetNextSelected()
	assert.Equal(suite.T(), ErrNoTask, err)

	t := NewTask("test", time.Now())
	t.Attrs[TagSelected] = true

	suite.DB.PutTasks(t)

	next, err := suite.DB.GetNextSelected()
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), t, next)
}

func TestTaskDBSuite(t *testing.T) {
	suite.Run(t, new(TaskDBSuite))
}
