package finch

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"

	"testing"
	"time"
)

func TestTaskDBLifecycle(t *testing.T) {
	name := "_taskdb_lifecycle"
	db, err := NewTaskDB(name)
	defer os.RemoveAll(name)

	assert.Nil(t, err)

	db.Close()

	assert.Nil(t, db.db)
	assert.Nil(t, db.wo)
	assert.Nil(t, db.ro)
}

type TaskDBSuite struct {
	suite.Suite
	db    *TaskDB
	fname string
}

func (suite *TaskDBSuite) SetupTest() {
	suite.fname = "_taskdb_suite"
	db, err := NewTaskDB(suite.fname)
	assert.Nil(suite.T(), err)
	suite.db = db
}

func (suite *TaskDBSuite) TearDownTest() {
	suite.db.Close()
	os.RemoveAll(suite.fname)
}

func (suite *TaskDBSuite) TestTasksIndexing() {
	t := NewTask("test", time.Now())

	err := suite.db.PutTasks(t)
	assert.Nil(suite.T(), err)

	task, err := suite.db.GetTask(t.Key(TasksIndex))
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), t, task)
}

func TestTaskDBSuite(t *testing.T) {
	suite.Run(t, new(TaskDBSuite))
}
