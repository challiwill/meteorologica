package migrations

import (
	"database/sql"
	"time"

	"github.com/BurntSushi/migration"
	"github.com/Sirupsen/logrus"
	"github.com/challiwill/meteorologica/db"
	"github.com/challiwill/meteorologica/errare"
)

func LockDBAndMigrate(log *logrus.Logger, sqlDriver, username, password, address, name string) (*db.Client, error) {
	log.Debug("Entering db.Ping")
	defer log.Debug("Returning db.Ping")

	if username == "" && password != "" {
		return nil, errare.NewCreationError("database client", "cannot have a database password with a username\n Please set the DB_PASSWORD environment variable")
	}

	sqlDataSource := username + ":" + password + "@" + "tcp(" + address + ")/" + name
	lockName := "meteorologica-mysql-migration-lock"

	dbLockConn, err := sql.Open(sqlDriver, sqlDataSource)
	if err != nil {
		return nil, err
	}
	defer dbLockConn.Close()

	for {
		var res int
		err := dbLockConn.QueryRow(`SELECT GET_LOCK(?, 10)`, lockName).Scan(&res)
		if err != nil || res != 1 {
			errStr := "timed out"
			if err != nil {
				errStr = err.Error()
			}
			log.Warnf("failed to get lock: %s. Trying again...", errStr)
			time.Sleep(5 * time.Second)
			continue
		}

		defer func() {
			_, err = dbLockConn.Exec(`SELECT RELEASE_LOCK(?)`, lockName)
			if err != nil {
				log.Error("failed to release lock: ", err)
			}
			log.Debug("migration lock released")
		}()
		log.Debug("migration lock acquired")

		dbClient, err := migration.OpenWith(sqlDriver, sqlDataSource, Migrations, mariadbGetVersion, mariadbSetVersion)
		if err != nil {
			log.Fatal("failed to run migrations: ", err)
		}

		return db.NewClientWith(log, dbClient), nil
	}
}

func mariadbGetVersion(tx migration.LimitedTx) (int, error) {
	v, err := getVersion(tx)
	if err != nil {
		if err := createVersionTable(tx); err != nil {
			return 0, err
		}
		return getVersion(tx)
	}
	return v, nil
}

func mariadbSetVersion(tx migration.LimitedTx, version int) error {
	if err := setVersion(tx, version); err != nil {
		if err := createVersionTable(tx); err != nil {
			return err
		}
		return setVersion(tx, version)
	}
	return nil
}

func getVersion(tx migration.LimitedTx) (int, error) {
	var version int
	r := tx.QueryRow("SELECT version FROM migration_version")
	if err := r.Scan(&version); err != nil {
		return 0, err
	}
	return version, nil
}

func setVersion(tx migration.LimitedTx, version int) error {
	_, err := tx.Exec("UPDATE migration_version SET version = ?", version)
	return err
}

func createVersionTable(tx migration.LimitedTx) error {
	_, err := tx.Exec(`
		CREATE TABLE migration_version (
			version INTEGER
		);`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO migration_version (version) VALUES (0);`)
	return err
}
