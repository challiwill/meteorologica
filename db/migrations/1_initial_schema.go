package migrations

import "github.com/BurntSushi/migration"

func InitialSchema(tx migration.LimitedTx) error {
	_, err := tx.Exec(`
					CREATE TABLE resource_billing (
						id VARCHAR(11) PRIMARY KEY,
						account_number VARCHAR(15) NOT NULL,
						account_name VARCHAR(30),
						day TINYINT(2) NOT NULL,
						month TINYINT(2) NOT NULL,
						year SMALLINT(4) NOT NULL,
						service_type VARCHAR(300) NOT NULL,
						region VARCHAR(10),
						resource VARCHAR(10) NOT NULL,
						usage_quantity DOUBLE NOT NULL,
						unit_of_measure VARCHAR(10),
						cost DOUBLE NOT NULL
					)
	`)
	if err != nil {
		return err
	}

	return nil
}
