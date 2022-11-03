package boot

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"casper-dao-middleware/pkg/config"
	"casper-dao-middleware/pkg/initializer"
	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
)

const (
	defaultInitInterval = 5
	defaultInitTimeout  = 300
)

// InitMySQL initialize connection to MySQL
func InitMySQL(ctx context.Context, dbConfig config.DBConfig) (*sqlx.DB, error) {
	expInitializer := initializer.NewExponential(defaultInitTimeout, defaultInitInterval)

	rawConn, err := expInitializer.Run(ctx, initMySQLConnection, dbConfig.DatabaseURI)
	if err != nil {
		return nil, err
	}

	conn, ok := rawConn.(*sqlx.DB)
	if !ok {
		return nil, errors.New("invalid connection format")
	}

	conn.SetMaxIdleConns(dbConfig.MaxIdleConnections)
	conn.SetMaxOpenConns(dbConfig.MaxOpenConnections)

	return conn, nil
}

func CloseMySQL(db *sqlx.DB) {
	if err := db.Close(); err != nil {
		log.Printf("Failed to properly close connection to MySQL: %s", err)
	}
	log.Println("MySQL connection properly closed")
}

func initMySQLConnection(connectionURI string) (interface{}, error) {
	connect, err := sqlx.Open("mysql", connectionURI)
	if err != nil {
		log.Printf("Failed to open mysql connection: %s", err)
		return nil, err
	}

	if err := connect.Ping(); err != nil {
		log.Printf("Could not ping mysql: %s", err)
		return nil, err
	}

	log.Println("Successfully connected to MySQL")
	return connect, nil
}

func SetUpTestDB() *sqlx.DB {
	return setUpTestMySQL(os.Getenv("TEST_DATABASE_URI"))
}

func setUpTestMySQL(uri string) *sqlx.DB {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	db, err := InitMySQL(ctx, config.DBConfig{
		DatabaseURI:        uri,
		MaxOpenConnections: 5,
		MaxIdleConnections: 5,
	})
	if err != nil {
		panic(err)
	}

	return db
}
