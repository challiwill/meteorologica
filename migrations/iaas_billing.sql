CREATE TABLE cf_82db49c0_7f6f_4844_a8af_97788db5e92e.iaas_billing
(
 AccountNumber VARCHAR(15) NOT NULL, 
AccountName VARCHAR(30),
Day TINYINT(2) NOT NULL,
Month CHAR(9) NOT NULL,
Year SMALLINT(4) NOT NULL,
ServiceType VARCHAR(30) NOT NULL,
UsageQuantity DOUBLE NOT NULL,
Cost DECIMAL(15,2) NOT NULL,
Region VARCHAR(10),
UnitOfMeasure VARCHAR(10),
IAAS VARCHAR(10) NOT NULL,
UNIQUE KEY(AccountNumber, Day, Month, Year, ServiceType, Region, IAAS)
);
