package finch

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

type TaskDB struct {
	db *leveldb.DB
	wo *opt.WriteOptions
	ro *opt.ReadOptions
}

func NewTaskDB(store storage.Storage) (*TaskDB, error) {
	tdb := new(TaskDB)

	// Open the Database with the provided Storage
	options := &opt.Options{}
	db, err := leveldb.Open(store, options)
	if err != nil {
		return tdb, err
	}
	tdb.db = db

	// Set default read and write options
	tdb.wo = &opt.WriteOptions{
		Sync: true,
	}
	tdb.ro = &opt.ReadOptions{
		DontFillCache: false,
	}

	return tdb, nil
}

func (tdb *TaskDB) Close() {
	tdb.db.Close()
	tdb.db = nil
}

// PutTasks inserts tasks into the database and overwrites those which
// already exist
func (tdb *TaskDB) PutTasks(tasks ...*Task) error {
	batch := new(leveldb.Batch)

	for i := 0; i < len(tasks); i++ {
		task := tasks[i]

		szd, err := task.Serialize()
		if err != nil {
			return err
		}
		key := task.Key(TasksIndex)

		// write the Task to the main storage
		batch.Put(key.Serialize(), szd)
	}

	if err := tdb.db.Write(batch, tdb.wo); err != nil {
		return err
	}

	return nil
}

func (tdb *TaskDB) GetTask(key *Key) (*Task, error) {
	szd, err := tdb.db.Get(key.Serialize(), tdb.ro)
	if err != nil {
		return new(Task), err
	}

	task, err := DeserializeTask(szd)
	return task, err
}
