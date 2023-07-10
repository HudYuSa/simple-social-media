package connection

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/HudYuSa/comments/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// migrate_up:
// migrate -path database/migrations/ -database "postgres://cyicsnej:pkzAVWH--U1AcE4IMQb3rfWvyE-gYK22@arjuna.db.elephantsql.com/cyicsnej?sslmode=disable" -verbose up

// migrate_down:
// migrate -path database/migrations/ -database "postgres://cyicsnej:pkzAVWH--U1AcE4IMQb3rfWvyE-gYK22@arjuna.db.elephantsql.com/cyicsnej?sslmode=disable" -verbose down

// migrateUp := exec.Command("make", "migrate_up")

// migrateErr := migrateUp.Run()
// if migrateErr != nil {
// 	panic(migrateErr)
// }

func ConnectDB(config *config.Config) {
	var err error
	// dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", config.DBHost, config.DBUserName, config.DBUserPassword, config.DBName, config.DBPort)
	dsn := config.DSN

	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// log gorm query
	newLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold:             time.Second, // Slow SQL threshold
		LogLevel:                  logger.Info, // Log level
		IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
		Colorful:                  true,        // Disable color
	})

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Fatal("Failed to connect to the Database")
	}

	// migrate
	m, err := migrate.New("file://"+cwd+"/database/migrations", dsn)
	if err != nil {
		log.Fatal(err)
	}
	err = m.Up()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("? Connected Successfully to the Database")
}
