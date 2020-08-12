package database

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"strconv"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh"
	"golang.org/x/xerrors"
)

func init() {
	RegisterOpen("mysql", mysqlOpen)
	RegisterFactory("mysql", NewMySQLDBRepository)
}

func mysqlOpen(dbConnCfg *DBConfig) (*DBConnection, error) {
	var (
		conn    *sql.DB
		sshConn *ssh.Client
	)
	cfg, err := genMysqlConfig(dbConnCfg)
	if err != nil {
		return nil, err
	}

	if dbConnCfg.SSHCfg != nil {
		dbConn, dbSSHConn, err := openMySQLViaSSH(cfg.FormatDSN(), dbConnCfg.SSHCfg)
		if err != nil {
			return nil, err
		}
		conn = dbConn
		sshConn = dbSSHConn
	} else {
		dbConn, err := sql.Open("mysql", cfg.FormatDSN())
		if err != nil {
			return nil, err
		}
		conn = dbConn
	}
	if err := conn.Ping(); err != nil {
		return nil, err
	}

	conn.SetMaxIdleConns(DefaultMaxIdleConns)
	conn.SetMaxOpenConns(DefaultMaxOpenConns)

	return &DBConnection{
		Conn:    conn,
		SSHConn: sshConn,
	}, nil
}

type MySQLViaSSHDialer struct {
	client *ssh.Client
}

func (d *MySQLViaSSHDialer) Dial(ctx context.Context, addr string) (net.Conn, error) {
	return d.client.Dial("tcp", addr)
}

func openMySQLViaSSH(dsn string, sshCfg *SSHConfig) (*sql.DB, *ssh.Client, error) {
	sshConfig, err := sshCfg.ClientConfig()
	if err != nil {
		return nil, nil, err
	}
	sshConn, err := ssh.Dial("tcp", sshCfg.Endpoint(), sshConfig)
	if err != nil {
		return nil, nil, xerrors.Errorf("cannot ssh dial, %+v", err)
	}
	mysql.RegisterDialContext("mysql+tcp", (&MySQLViaSSHDialer{sshConn}).Dial)
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, nil, xerrors.Errorf("cannot connect database, %+v", err)
	}
	return conn, sshConn, nil
}

func genMysqlConfig(connCfg *DBConfig) (*mysql.Config, error) {
	cfg := mysql.NewConfig()

	if connCfg.DataSourceName != "" {
		return mysql.ParseDSN(connCfg.DataSourceName)
	}

	cfg.User = connCfg.User
	cfg.Passwd = connCfg.Passwd
	cfg.DBName = connCfg.DBName

	switch connCfg.Proto {
	case ProtoTCP, ProtoUDP:
		host, port := connCfg.Host, connCfg.Port
		if host == "" {
			host = "127.0.0.1"
		}
		if port == 0 {
			port = 3306
		}
		cfg.Addr = host + ":" + strconv.Itoa(port)
		cfg.Net = string(connCfg.Proto)
	case ProtoUnix:
		if connCfg.Path != "" {
			cfg.Addr = "/tmp/mysql.sock"
			break
		}
		cfg.Addr = connCfg.Path
		cfg.Net = string(connCfg.Proto)
	default:
		return nil, fmt.Errorf("default addr for network %s unknown", connCfg.Proto)
	}

	cfg.Params = connCfg.Params

	return cfg, nil
}

type MySQLDBRepository struct {
	Conn *sql.DB
}

func NewMySQLDBRepository(conn *sql.DB) DBRepository {
	return &MySQLDBRepository{Conn: conn}
}

func (db *MySQLDBRepository) CurrentDatabase(ctx context.Context) (string, error) {
	row := db.Conn.QueryRowContext(ctx, "SELECT DATABASE()")
	var database string
	if err := row.Scan(&database); err != nil {
		return "", err
	}
	return database, nil
}

func (db *MySQLDBRepository) Databases(ctx context.Context) ([]string, error) {
	rows, err := db.Conn.QueryContext(ctx, "select SCHEMA_NAME from information_schema.SCHEMATA")
	if err != nil {
		return nil, err
	}
	databases := []string{}
	for rows.Next() {
		var database string
		if err := rows.Scan(&database); err != nil {
			return nil, err
		}
		databases = append(databases, database)
	}
	return databases, nil
}

func (db *MySQLDBRepository) CurrentSchema(ctx context.Context) (string, error) {
	return db.CurrentDatabase(ctx)
}

func (db *MySQLDBRepository) Schemas(ctx context.Context) ([]string, error) {
	return db.Databases(ctx)
}

func (db *MySQLDBRepository) SchemaTables(ctx context.Context) (map[string][]string, error) {
	rows, err := db.Conn.QueryContext(
		ctx,
		`
	SELECT 
		TABLE_SCHEMA,
		TABLE_NAME
	FROM
		information_schema.TABLES
	ORDER BY
		TABLE_SCHEMA,
		TABLE_NAME
	`)
	if err != nil {
		return nil, err
	}
	databaseTables := map[string][]string{}
	for rows.Next() {
		var schema, table string
		if err := rows.Scan(&schema, &table); err != nil {
			return nil, err
		}

		if arr, ok := databaseTables[schema]; ok {
			databaseTables[schema] = append(arr, table)
		} else {
			databaseTables[schema] = []string{table}
		}
	}
	return databaseTables, nil
}

func (db *MySQLDBRepository) Tables(ctx context.Context) ([]string, error) {
	rows, err := db.Conn.QueryContext(ctx, "SHOW TABLES")
	if err != nil {
		return nil, err
	}
	tables := []string{}
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}
	return tables, nil
}

func (db *MySQLDBRepository) DescribeDatabaseTable(ctx context.Context) ([]*ColumnDesc, error) {
	rows, err := db.Conn.QueryContext(
		ctx,
		`
SELECT
	TABLE_SCHEMA,
	TABLE_NAME,
	COLUMN_NAME,
	COLUMN_TYPE,
	IS_NULLABLE,
	COLUMN_KEY,
	COLUMN_DEFAULT,
	EXTRA
FROM information_schema.COLUMNS
`)
	if err != nil {
		return nil, err
	}
	tableInfos := []*ColumnDesc{}
	for rows.Next() {
		var tableInfo ColumnDesc
		err := rows.Scan(
			&tableInfo.Schema,
			&tableInfo.Table,
			&tableInfo.Name,
			&tableInfo.Type,
			&tableInfo.Null,
			&tableInfo.Key,
			&tableInfo.Default,
			&tableInfo.Extra,
		)
		if err != nil {
			return nil, err
		}
		tableInfos = append(tableInfos, &tableInfo)
	}
	return tableInfos, nil
}

func (db *MySQLDBRepository) DescribeDatabaseTableBySchema(ctx context.Context, schemaName string) ([]*ColumnDesc, error) {
	rows, err := db.Conn.QueryContext(
		ctx,
		`
SELECT
	TABLE_SCHEMA,
	TABLE_NAME,
	COLUMN_NAME,
	COLUMN_TYPE,
	IS_NULLABLE,
	COLUMN_KEY,
	COLUMN_DEFAULT,
	EXTRA
FROM information_schema.COLUMNS
WHERE information_schema.COLUMNS.TABLE_SCHEMA = ?
`, schemaName)
	if err != nil {
		return nil, err
	}
	tableInfos := []*ColumnDesc{}
	for rows.Next() {
		var tableInfo ColumnDesc
		err := rows.Scan(
			&tableInfo.Schema,
			&tableInfo.Table,
			&tableInfo.Name,
			&tableInfo.Type,
			&tableInfo.Null,
			&tableInfo.Key,
			&tableInfo.Default,
			&tableInfo.Extra,
		)
		if err != nil {
			return nil, err
		}
		tableInfos = append(tableInfos, &tableInfo)
	}
	return tableInfos, nil
}

func (db *MySQLDBRepository) Exec(ctx context.Context, query string) (sql.Result, error) {
	return db.Conn.ExecContext(ctx, query)
}

func (db *MySQLDBRepository) Query(ctx context.Context, query string) (*sql.Rows, error) {
	return db.Conn.QueryContext(ctx, query)
}
