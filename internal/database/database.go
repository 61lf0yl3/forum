package database

import (
	"database/sql"
)

// Database ...
type Database struct {
	config *Config
	db     *sql.DB
}

// NewDB ...
func NewDB(config *Config) *Database {
	return &Database{
		config: config,
	}
}

// InitDB ...
func (d *Database) InitDB() error {
	if err := d.Open(); err != nil {
		return err
	}
	d.BuildSchema()
	d.Close()
	return nil
}

// Open ...
func (d *Database) Open() error {
	store, err := sql.Open("sqlite3", d.config.DatabaseAddress)
	if err != nil {
		return err
	}
	// Complete checking of connection with DB
	if err := store.Ping(); err != nil {
		return err
	}
	d.db = store // fill "db" field with completely configured DB
	return nil
}

// Close DB
func (d *Database) Close() {
	d.db.Close()
}
