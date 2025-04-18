package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func ConnectDatabase() (*gorm.DB, error) {
	HOST := os.Getenv("DB_HOST")
	USER := os.Getenv("DB_USERNAME")
	PASSWORD := os.Getenv("DB_PASSWORD")
	NAME := os.Getenv("DB_NAME")
	PORT := os.Getenv("DB_PORT")
	SSLMODE := os.Getenv("DB_SSLMODE")
	if HOST == "" || USER == "" || PASSWORD == "" || NAME == "" || PORT == "" || SSLMODE == "" {
		log.Println("DATABASE HOST , USER, PASSWORD, NAME, PORT, SSLMODE missing in configuration.")
		return nil, fmt.Errorf("DATABASE HOST , USER, PASSWORD, NAME, PORT, SSLMODE missing in configuration")
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", HOST, USER, PASSWORD, NAME, PORT, SSLMODE)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		SkipDefaultTransaction: true, // skip wrapping operations in transaction

		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "app_", // prefix for all tables (e.g. app_users)
			SingularTable: false,  // keep plural table names
			NameReplacer:  nil,    // can use strings.NewReplacer() to customize further
		},
		Logger: logger.Default.LogMode(logger.Info), // log SQL at Info level

		NowFunc: func() time.Time {
			return time.Now().UTC() // use UTC time instead of local
		},
		DryRun: false, // if true, SQL statements are generated but not run

		PrepareStmt: true, // cache prepared statements for reuse

		DisableNestedTransaction: false, // allow nested transactions (savepoints)

		AllowGlobalUpdate: false, // block updates without WHERE clause

		DisableAutomaticPing:                     false, // ping DB after connecting to check it's up
		DisableForeignKeyConstraintWhenMigrating: true,  // don't auto-create FK constraints in DB

	})

	if err != nil {
		return nil, err
	}
	log.Println("Connected to Database")

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDb.SetMaxOpenConns(20)
	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetConnMaxLifetime(20 * time.Minute)
	sqlDb.SetConnMaxIdleTime(5 * time.Minute)
	log.Println("✅ DB connection pool configured successfully!")

	DB = db

	return DB, nil

}
