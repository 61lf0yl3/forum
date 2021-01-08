package database

import (
	"fmt"
	"strings"

	"github.com/astgot/forum/internal/model"
)

// InsertThreadInfo ..
func (d *Database) InsertThreadInfo(threadName string, postID int64) error {
	stmnt, err := d.db.Prepare("INSERT INTO Threads (Name) VALUES (?)")
	if err != nil {
		fmt.Println("insert Threads error")
		return err
	}
	defer stmnt.Close()
	res, err := stmnt.Exec(threadName)
	if err != nil {
		fmt.Println(err.Error(), "---> exec Threads error")
		// if thread exists in DB, we will get his threadID, and insert to PostMapping Table
		if err.Error() == "UNIQUE constraint failed: Threads.Name" {
			ID := d.GetThreadID(threadName)
			d.InsertPostMapInfo(postID, ID)
		}
		return err
	}
	threadID, _ := res.LastInsertId()
	d.InsertPostMapInfo(postID, threadID)

	return nil
}

// GetAllThreads ...
func (d *Database) GetAllThreads() []*model.Thread {
	var threads []*model.Thread
	res, err := d.db.Query("SELECT * FROM Threads")
	if err != nil {
		fmt.Println("thread query error")
		return nil
	}
	defer res.Close()
	for res.Next() {
		thread := model.NewThread()
		if err := res.Scan(&thread.ID, &thread.Name); err != nil {
			fmt.Println("thread scan error")
			return nil
		}
		threads = append(threads, thread)

	}
	return threads
}

// GetThreadID ...
func (d *Database) GetThreadID(name string) int64 {
	var id int64
	if err := d.db.QueryRow("SELECT ID FROM Threads WHERE Name = ?", name).Scan(&id); err != nil {
		fmt.Println("threadID retrieve error")
	}

	return id
}

// GetThreadByID ...
func (d *Database) GetThreadByID(id int64) (*model.Thread, error) {
	thread := model.NewThread()
	if err := d.db.QueryRow("SELECT * FROM Threads WHERE ID = ?", id).Scan(&thread.ID, &thread.Name); err != nil {
		fmt.Println("error on func GetThreadByID()")
		return nil, err
	}
	return thread, nil
}

// SearchThread ...
func (d *Database) SearchThread(search string) ([]*model.Post, error) {
	if err := d.Open(); err != nil {
		return nil, err
	}
	threads := d.GetAllThreads()
	threadIDs := []int64{}
	// Collect all threads IDs by comparing thread Names with search query
	for _, thread := range threads {
		if ok := strings.Contains(thread.Name, search); ok {
			threadIDs = append(threadIDs, thread.ID)
		}
	}
	if len(threadIDs) == 0 {
		return nil, nil
	}
	var result []*model.Post
	for _, threadID := range threadIDs {
		posts, err := d.FindPostsByThreadID(threadID)
		if err != nil {
			fmt.Println("SearchThread", err.Error())
			return nil, err
		}
		result = append(result, posts...)
	}
	return result, nil
}
