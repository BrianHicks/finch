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

	db, err := NewTaskDB(filestore)
	assert.Nil(t, err)
	defer os.RemoveAll(name)

	assert.Nil(t, err)

	db.Close()

	assert.Nil(t, db.db)
}

type TaskDBSuite struct {
	suite.Suite
	db *TaskDB
}

func (suite *TaskDBSuite) SetupTest() {
	db, err := NewTaskDB(storage.NewMemStorage())
	assert.Nil(suite.T(), err)
	suite.db = db
}

func (suite *TaskDBSuite) TearDownTest() {
	suite.db.Close()
}

func (suite *TaskDBSuite) TestTasksIndexing() {
	t := NewTask("test", time.Now())

	err := suite.db.PutTasks(t)
	assert.Nil(suite.T(), err)

	task, err := suite.db.GetTask(t.Key())
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), t, task)
}

func (suite *TaskDBSuite) TestPendingIndexing() {
	nope := NewTask("test", time.Now())
	nope.Attrs[TagPending] = false

	yep := NewTask("test", time.Now())
	nope.Attrs[TagPending] = true

	suite.db.PutTasks(nope, yep)

	pending, err := suite.db.GetPendingTasks()
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), []*Task{yep}, pending)
}

func (suite *TaskDBSuite) TestSelectedIndexing() {
	nope := NewTask("test", time.Now())
	nope.Attrs[TagSelected] = false

	yep := NewTask("selected", time.Now())
	yep.Attrs[TagSelected] = true

	suite.db.PutTasks(nope, yep)

	selected, err := suite.db.GetSelectedTasks()
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), []*Task{yep}, selected)
}

func (suite *TaskDBSuite) TestGetNextSelected() {
	// try it empty first
	_, err := suite.db.GetNextSelected()
	assert.Equal(suite.T(), ErrNoNextTask, err)

	t := NewTask("test", time.Now())
	t.Attrs[TagSelected] = true

	suite.db.PutTasks(t)

	next, err := suite.db.GetNextSelected()
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), t, next)
}

func TestTaskDBSuite(t *testing.T) {
	suite.Run(t, new(TaskDBSuite))
}
