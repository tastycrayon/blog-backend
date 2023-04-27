package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/tastycrayon/blog-backend/util"
)

// store provides all db functions
type Store struct {
	Queries *Queries
	db      *sql.DB
}

// NewStore creates new store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
// 	tx, err := store.db.BeginTx(ctx, nil)
// 	if err != nil {
// 		return err
// 	}
// 	q := New(tx)
// 	err = fn(q)
// 	if err != nil {
// 		if rbErr := tx.Rollback(); rbErr != nil {
// 			return fmt.Errorf("tx Err: %v, Rb Err: %v", err, rbErr)
// 		}
// 		return err
// 	}
// 	return tx.Commit()
// }

func InitDB(config util.Config) (*sql.DB, error) {
	var args string = "?parseTime=true&multiStatements=true&tls=preferred"
	var dsn string = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBDockerPort,
		config.DBName,
	)
	db, err := sql.Open("mysql", dsn+args)

	if err != nil {
		log.Printf("❗failed to connect: %v", err)
		return nil, err
	}
	// connection setting
	// Maximum Idle Connections
	db.SetMaxIdleConns(5)
	// Maximum Open Connections
	db.SetMaxOpenConns(10)
	// Idle Connection Timeout
	db.SetConnMaxIdleTime(1 * time.Second)
	// Connection Lifetime
	db.SetConnMaxLifetime(30 * time.Second)

	// migrate start
	runMigrationOnDB(config, db)
	// migrate start
	if err := db.Ping(); err != nil {
		log.Printf("❗ failed to ping: %v", err)
		return nil, err
	}

	var version string
	err = db.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	log.Printf("🤟 Successfully connected to MySQL %v!\n", version)
	// defer db.Close() // will close in query

	return db, nil
}

func runMigrationOnDB(config util.Config, db *sql.DB) {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatalf("❗ failed to connect driver: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		config.MigrationURL, // Ex: file:///db/migrate
		config.DBDriver,     // Ex: mysql
		driver,
	)
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("❗ failed to migrate: %v", err)
	}
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("❗ failed to migrate: %v", err)
	}
	if err == migrate.ErrNoChange {
		log.Println("🎉 Migrated Successfully! (No Change)")
		return
	}
	log.Println("🎉 Migrated Successfully! (No Change)")
}
