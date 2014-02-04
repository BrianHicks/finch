/*
Package persist wraps a LevelDB database, adding prefix iteration for indexing
and logged operations (for future syncing)

To get a database, use `NewFile`:

    store, err := NewFile("testDB")

Then you can store and retrieve documents:

    err := store.Put([]byte("key"), []byte("value"))
    doc, err := store.Get([]byte("key"))

Batch operations (logged to "_log" prefix):

    batch := NewLoggedBatch()
    batch.Put([]byte("key"), []byte("value"))
    batch.Delete([]byte("otherKey"))

    err = store.Write(batch.Batch)

And range operations:

    firstDoc, err := db.Prefix([]byte("somePrefix")).First() // Also Last and All
*/
package persist
