package migrations

import "github.com/BurntSushi/migration"

func ForgotNotNull(tx migration.LimitedTx) error {
	_, err := tx.Exec(`
					ALTER TABLE resource_billing
					CHANGE COLUMN account_number account_number VARCHAR(255) NOT NULL
	`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
					ALTER TABLE resource_billing
					CHANGE COLUMN account_name account_name VARCHAR(255)
	`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
					ALTER TABLE resource_billing
					CHANGE COLUMN service_type service_type VARCHAR(255) NOT NULL
	`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
					ALTER TABLE resource_billing
					CHANGE COLUMN region region VARCHAR(255)
	`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
					ALTER TABLE resource_billing
					CHANGE COLUMN resource resource VARCHAR(255) NOT NULL
	`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
					ALTER TABLE resource_billing
					CHANGE COLUMN unit_of_measure unit_of_measure VARCHAR(255)
	`)
	if err != nil {
		return err
	}

	return nil
}
