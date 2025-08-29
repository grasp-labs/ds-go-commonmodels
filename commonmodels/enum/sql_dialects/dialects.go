package sqldialects

type DatabaseDialect string

const (
	DialectPostgres  DatabaseDialect = "postgres"
	DialectMySQL     DatabaseDialect = "mysql"
	DialectSQLServer DatabaseDialect = "sqlserver" // aka mssql
	DialectSQLite    DatabaseDialect = "sqlite"
	DialectSnowflake DatabaseDialect = "snowflake"
	DialectBigQuery  DatabaseDialect = "bigquery"
	DialectRedshift  DatabaseDialect = "redshift"
	DialectOracle    DatabaseDialect = "oracle"
	DialectDuckDB    DatabaseDialect = "duckdb"
	DialectTrino     DatabaseDialect = "trino"
)
