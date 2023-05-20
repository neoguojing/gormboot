package gormboot

import (
	"fmt"
	"sync"
)

type DatabaseType string

const (
	MySQL    DatabaseType = "mysql"
	Postgres DatabaseType = "postgres"
	SQLite   DatabaseType = "sqlite"
)

type DatabaseInstance interface {
	Connect() error
	Disconnect() error
}

type MySQLInstance struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

func (m *MySQLInstance) Connect() error {
	fmt.Printf("Connecting to MySQL database %s:%d\n", m.Host, m.Port)
	// Connect to MySQL database
	return nil
}

func (m *MySQLInstance) Disconnect() error {
	fmt.Printf("Disconnecting from MySQL database %s:%d\n", m.Host, m.Port)
	// Disconnect from MySQL database
	return nil
}

type PostgresInstance struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

func (p *PostgresInstance) Connect() error {
	fmt.Printf("Connecting to Postgres database %s:%d\n", p.Host, p.Port)
	// Connect to Postgres database
	return nil
}

func (p *PostgresInstance) Disconnect() error {
	fmt.Printf("Disconnecting from Postgres database %s:%d\n", p.Host, p.Port)
	// Disconnect from Postgres database
	return nil
}

type SQLiteInstance struct {
	DatabasePath string
}

func (s *SQLiteInstance) Connect() error {
	fmt.Printf("Connecting to SQLite database %s\n", s.DatabasePath)
	// Connect to SQLite database
	return nil
}

func (s *SQLiteInstance) Disconnect() error {
	fmt.Printf("Disconnecting from SQLite database %s\n", s.DatabasePath)
	// Disconnect from SQLite database
	return nil
}

type DatabaseInstanceFactory struct {
	instances map[DatabaseType]DatabaseInstance
	mutex     sync.Mutex
}

func (f *DatabaseInstanceFactory) GetInstance(dbType DatabaseType) (DatabaseInstance, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	instance, ok := f.instances[dbType]
	if !ok {
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}

	return instance, nil
}

func (f *DatabaseInstanceFactory) RegisterInstance(dbType DatabaseType, instance DatabaseInstance) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	_, ok := f.instances[dbType]
	if ok {
		return fmt.Errorf("database type %s already registered", dbType)
	}

	f.instances[dbType] = instance

	return nil
}
