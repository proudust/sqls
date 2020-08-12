package database

import (
	"context"
	"database/sql"
)

type MockDBRepository struct {
	MockDatabase                      func(context.Context) (string, error)
	MockDatabases                     func(context.Context) ([]string, error)
	MockDatabaseTables                func(context.Context) (map[string][]string, error)
	MockTables                        func(context.Context) ([]string, error)
	MockDescribeTable                 func(context.Context, string) ([]*ColumnDesc, error)
	MockDescribeDatabaseTable         func(context.Context) ([]*ColumnDesc, error)
	MockDescribeDatabaseTableBySchema func(context.Context, string) ([]*ColumnDesc, error)
	MockExec                          func(context.Context, string) (sql.Result, error)
	MockQuery                         func(context.Context, string) (*sql.Rows, error)
}

func NewMockDBRepository(conn *sql.DB) DBRepository {
	return &MockDBRepository{
		MockDatabase:       func(ctx context.Context) (string, error) { return "world", nil },
		MockDatabases:      func(ctx context.Context) ([]string, error) { return dummyDatabases, nil },
		MockDatabaseTables: func(ctx context.Context) (map[string][]string, error) { return dummyDatabaseTables, nil },
		MockTables:         func(ctx context.Context) ([]string, error) { return dummyTables, nil },
		MockDescribeTable: func(ctx context.Context, tableName string) ([]*ColumnDesc, error) {
			switch tableName {
			case "city":
				return dummyCityColumns, nil
			case "country":
				return dummyCountryColumns, nil
			case "countrylanguage":
				return dummyCountryLanguageColumns, nil
			}
			return nil, nil
		},
		MockDescribeDatabaseTable: func(ctx context.Context) ([]*ColumnDesc, error) {
			res := []*ColumnDesc{}
			res = append(res, dummyCityColumns...)
			res = append(res, dummyCountryColumns...)
			res = append(res, dummyCountryLanguageColumns...)
			return res, nil

		},
		MockDescribeDatabaseTableBySchema: func(ctx context.Context, schemaName string) ([]*ColumnDesc, error) {
			res := []*ColumnDesc{}
			res = append(res, dummyCityColumns...)
			return res, nil

		},
		MockExec: func(ctx context.Context, query string) (sql.Result, error) {
			return &MockResult{
				MockLastInsertID: func() (int64, error) { return 11, nil },
				MockRowsAffected: func() (int64, error) { return 22, nil },
			}, nil
		},
		MockQuery: func(ctx context.Context, query string) (*sql.Rows, error) {
			return &sql.Rows{}, nil
		},
	}
}

func (m *MockDBRepository) CurrentDatabase(ctx context.Context) (string, error) {
	return m.MockDatabase(ctx)
}

func (m *MockDBRepository) Databases(ctx context.Context) ([]string, error) {
	return m.MockDatabases(ctx)
}

func (m *MockDBRepository) CurrentSchema(ctx context.Context) (string, error) {
	return m.MockDatabase(ctx)
}

func (m *MockDBRepository) Schemas(ctx context.Context) ([]string, error) {
	return m.MockDatabases(ctx)
}

func (m *MockDBRepository) SchemaTables(ctx context.Context) (map[string][]string, error) {
	return m.MockDatabaseTables(ctx)
}

func (m *MockDBRepository) Tables(ctx context.Context) ([]string, error) {
	return m.MockTables(ctx)
}

func (m *MockDBRepository) DescribeDatabaseTable(ctx context.Context) ([]*ColumnDesc, error) {
	return m.MockDescribeDatabaseTable(ctx)
}

func (m *MockDBRepository) DescribeDatabaseTableBySchema(ctx context.Context, schemaName string) ([]*ColumnDesc, error) {
	return m.MockDescribeDatabaseTableBySchema(ctx, schemaName)
}

func (m *MockDBRepository) Exec(ctx context.Context, query string) (sql.Result, error) {
	return m.MockExec(ctx, query)
}

func (m *MockDBRepository) Query(ctx context.Context, query string) (*sql.Rows, error) {
	return m.MockQuery(ctx, query)
}

var dummyDatabases = []string{
	"information_schema",
	"mysql",
	"performance_schema",
	"sys",
	"world",
}
var dummyDatabaseTables = map[string][]string{
	"world": []string{
		"city",
		"country",
		"countrylanguage",
	},
}
var dummyTables = []string{
	"city",
	"country",
	"countrylanguage",
}
var dummyCityColumns = []*ColumnDesc{
	{
		Schema: "world",
		Table:  "city",
		Name:   "ID",
		Type:   "int(11)",
		Null:   "NO",
		Key:    "PRI",
		Default: sql.NullString{
			String: "<null>",
			Valid:  false,
		},
		Extra: "auto_increment",
	},
	{
		Schema: "world",
		Table:  "city",
		Name:   "Name",
		Type:   "char(35)",
		Null:   "NO",
		Key:    "",
		Default: sql.NullString{
			String: "",
			Valid:  false,
		},
		Extra: "",
	},
	{
		Schema: "world",
		Table:  "city",
		Name:   "CountryCode",
		Type:   "char(3)",
		Null:   "NO",
		Key:    "MUL",
		Default: sql.NullString{
			String: "",
			Valid:  false,
		},
		Extra: "",
	},
	{
		Schema: "world",
		Table:  "city",
		Name:   "District",
		Type:   "char(20)",
		Null:   "NO",
		Key:    "",
		Default: sql.NullString{
			String: "",
			Valid:  false,
		},
		Extra: "",
	},
	{
		Schema: "world",
		Table:  "city",
		Name:   "Population",
		Type:   "int(11)",
		Null:   "NO",
		Key:    "",
		Default: sql.NullString{
			String: "",
			Valid:  false,
		},
		Extra: "",
	},
}
var dummyCountryColumns = []*ColumnDesc{
	{
		Schema: "world",
		Table:  "country",
		Name:   "Code",
		Type:   "char(3)",
		Null:   "NO",
		Key:    "PRI",
		Default: sql.NullString{
			String: "<null>",
			Valid:  false,
		},
		Extra: "auto_increment",
	},
	{
		Schema: "world",
		Table:  "country",
		Name:   "Name",
		Type:   "char(52)",
		Null:   "NO",
		Key:    "",
		Default: sql.NullString{
			String: "",
			Valid:  false,
		},
		Extra: "",
	},
	{
		Schema: "world",
		Table:  "country",
		Name:   "CountryCode",
		Type:   "char(3)",
		Null:   "NO",
		Key:    "",
		Default: sql.NullString{
			String: "",
			Valid:  false,
		},
		Extra: "",
	},
	{
		Schema: "world",
		Table:  "country",
		Name:   "Continent",
		Type:   "enum('Asia','Europe','North America','Africa','Oceania','Antarctica','South America')",
		Null:   "NO",
		Key:    "",
		Default: sql.NullString{
			String: "Asia",
			Valid:  false,
		},
		Extra: "",
	},
	{
		Schema: "world",
		Table:  "country",
		Name:   "Region",
		Type:   "char(26)",
		Null:   "NO",
		Key:    "",
		Default: sql.NullString{
			String: "",
			Valid:  false,
		},
		Extra: "",
	},
	{
		Schema: "world",
		Table:  "country",
		Name:   "SurfaceArea",
		Type:   "decimal(10,2)",
		Null:   "NO",
		Key:    "",
		Default: sql.NullString{
			String: "0.00",
			Valid:  false,
		},
		Extra: "auto_increment",
	},
	{
		Schema: "world",
		Table:  "country",
		Name:   "IndepYear",
		Type:   "smallint(6)",
		Null:   "YES",
		Key:    "",
		Default: sql.NullString{
			String: "0",
			Valid:  false,
		},
		Extra: "",
	},
	{
		Schema: "world",
		Table:  "country",
		Name:   "LifeExpectancy",
		Type:   "decimal(3,1)",
		Null:   "YES",
		Key:    "",
		Default: sql.NullString{
			String: "<null>",
			Valid:  false,
		},
		Extra: "",
	},
	{
		Schema: "world",
		Table:  "country",
		Name:   "GNP",
		Type:   "decimal(10,2)",
		Null:   "YES",
		Key:    "",
		Default: sql.NullString{
			String: "<null>",
			Valid:  false,
		},
		Extra: "",
	},
	{
		Schema: "world",
		Table:  "country",
		Name:   "GNPOld",
		Type:   "decimal(10,2)",
		Null:   "YES",
		Key:    "",
		Default: sql.NullString{
			String: "<null>",
			Valid:  false,
		},
		Extra: "",
	},
	{
		Schema: "world",
		Table:  "country",
		Name:   "LocalName",
		Type:   "char(45)",
		Null:   "NO",
		Key:    "",
		Default: sql.NullString{
			String: "",
			Valid:  false,
		},
		Extra: "",
	},
	{
		Schema: "world",
		Table:  "country",
		Name:   "GovernmentForm",
		Type:   "char(45)",
		Null:   "NO",
		Key:    "",
		Default: sql.NullString{
			String: "",
			Valid:  false,
		},
		Extra: "",
	},
	{
		Schema: "world",
		Table:  "country",
		Name:   "HeadOfState",
		Type:   "char(60)",
		Null:   "YES",
		Key:    "",
		Default: sql.NullString{
			String: "<null>",
			Valid:  false,
		},
		Extra: "",
	},
	{
		Schema: "world",
		Table:  "country",
		Name:   "Capital",
		Type:   "int(11)",
		Null:   "YES",
		Key:    "",
		Default: sql.NullString{
			String: "<null>",
			Valid:  false,
		},
		Extra: "",
	},
	{
		Schema: "world",
		Table:  "country",
		Name:   "Code2",
		Type:   "char(2)",
		Null:   "NO",
		Key:    "",
		Default: sql.NullString{
			String: "",
			Valid:  false,
		},
		Extra: "",
	},
}
var dummyCountryLanguageColumns = []*ColumnDesc{
	{
		Schema: "world",
		Table:  "countrylanguage",
		Name:   "CountryCode",
		Type:   "char(3)",
		Null:   "NO",
		Key:    "PRI",
		Default: sql.NullString{
			String: "",
			Valid:  false,
		},
		Extra: "",
	},
	{
		Schema: "world",
		Table:  "countrylanguage",
		Name:   "Language",
		Type:   "char(30)",
		Null:   "NO",
		Key:    "PRI",
		Default: sql.NullString{
			String: "",
			Valid:  false,
		},
		Extra: "",
	},
	{
		Schema: "world",
		Table:  "countrylanguage",
		Name:   "IsOfficial",
		Type:   "enum('T','F')",
		Null:   "NO",
		Key:    "F",
		Default: sql.NullString{
			String: "",
			Valid:  false,
		},
		Extra: "",
	},
	{
		Schema: "world",
		Table:  "countrylanguage",
		Name:   "Percentage",
		Type:   "decimal(4,1)",
		Null:   "NO",
		Key:    "",
		Default: sql.NullString{
			String: "0.0",
			Valid:  false,
		},
		Extra: "",
	},
}

type MockResult struct {
	MockLastInsertID func() (int64, error)
	MockRowsAffected func() (int64, error)
}

func (m *MockResult) LastInsertId() (int64, error) {
	return m.MockLastInsertID()
}
func (m *MockResult) RowsAffected() (int64, error) {
	return m.MockRowsAffected()
}

func init() {
	RegisterOpen("mock", func(connCfg *DBConfig) (*DBConnection, error) { return &DBConnection{}, nil })
	RegisterFactory("mock", NewMockDBRepository)
}
