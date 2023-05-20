package gormboot

import (
	"testing"
)

func TestGetInstance(t *testing.T) {
	factory := &DatabaseInstanceFactory{
		instances: make(map[DatabaseType]DatabaseInstance),
	}

	mysqlInstance := &MySQLInstance{
		Host:     "localhost",
		Port:     3306,
		Username: "root",
		Password: "password",
		Database: "mydatabase",
	}

	postgresInstance := &PostgresInstance{
		Host:     "localhost",
		Port:     5432,
		Username: "postgres",
		Password: "password",
		Database: "mydatabase",
	}

	sqliteInstance := &SQLiteInstance{
		DatabasePath: "/path/to/mydatabase.db",
	}

	factory.RegisterInstance(MySQL, mysqlInstance)
	factory.RegisterInstance(Postgres, postgresInstance)
	factory.RegisterInstance(SQLite, sqliteInstance)

	_, err := factory.GetInstance("invalid")
	if err == nil {
		t.Errorf("Expected error for invalid database type")
	}

	instance, err := factory.GetInstance(MySQL)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	if _, ok := instance.(*MySQLInstance); !ok {
		t.Errorf("Expected MySQLInstance, got %T", instance)
	}

	instance, err = factory.GetInstance(Postgres)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	if _, ok := instance.(*PostgresInstance); !ok {
		t.Errorf("Expected PostgresInstance, got %T", instance)
	}

	instance, err = factory.GetInstance(SQLite)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	if _, ok := instance.(*SQLiteInstance); !ok {
		t.Errorf("Expected SQLiteInstance, got %T", instance)
	}
}
