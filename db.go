package finch

import "github.com/jmhodges/levigo"

type TaskDB struct {
	db *levigo.DB
	wo *levigo.WriteOptions
	ro *levigo.ReadOptions
}

func NewTaskDB(path string) (*TaskDB, error) {
	tdb := new(TaskDB)

	opts := levigo.NewOptions()
	opts.SetCache(levigo.NewLRUCache(3 << 30))
	opts.SetCreateIfMissing(true)
	db, err := levigo.Open(path, opts)
	if err != nil {
		return tdb, err
	}
	tdb.db = db

	tdb.wo = levigo.NewWriteOptions()
	tdb.wo.SetSync(true)

	tdb.ro = levigo.NewReadOptions()
	tdb.ro.SetFillCache(true)

	return tdb, nil
}

func (tdb *TaskDB) Close() {
	tdb.db.Close()
	tdb.db = nil

	tdb.wo.Close()
	tdb.wo = nil

	tdb.ro.Close()
	tdb.ro = nil
}

// PutTasks inserts tasks into the database and overwrites those which
// already exist
func (tdb *TaskDB) PutTasks(tasks ...*Task) error {
	batch := levigo.NewWriteBatch()
	defer batch.Close()

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

	if err := tdb.db.Write(tdb.wo, batch); err != nil {
		return err
	}

	return nil
}

func (tdb *TaskDB) GetTask(key *Key) (*Task, error) {
	szd, err := tdb.db.Get(tdb.ro, key.Serialize())
	if err != nil {
		return new(Task), err
	}

	task, err := DeserializeTask(szd)
	return task, err
}
