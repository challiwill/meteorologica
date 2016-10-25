package migrations

import "github.com/BurntSushi/migration"

func LengthenIDsAgain(tx migration.LimitedTx) error {
	_, err := tx.Exec(`
					ALTER TABLE resource_billing
					CHANGE COLUMN id id VARCHAR(30)
	`)
	if err != nil {
		return err
	}

	return nil
}
