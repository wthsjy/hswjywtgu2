package mysql

import (
	"time"
	"context"
	"sync"
	"errors"
)

type safeMap struct {
	sync.RWMutex
	confs   map[string]interface{}
	clients map[string]interface{}
}

var sqlPool = &safeMap{
	confs:   make(map[string]interface{}),
	clients: make(map[string]interface{}),
}

// RegisterMySQL dsn format -> [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
func RegisterMySQL(name string, c *Config) error {
	// set default
	if c.Idle == 0 {
		c.Idle = 32
	}

	if c.Active == 0 {
		c.Active = 128
	}

	if c.IdleTimeout == 0 {
		c.IdleTimeout = time.Minute * 5
	}

	if c.TranTimeout == 0 {
		c.TranTimeout = time.Second * 5
	}

	if c.ExecTimeout == 0 {
		c.ExecTimeout = time.Second * 5
	}

	if c.QueryTimeout == 0 {
		c.QueryTimeout = time.Second * 5
	}

	sqlPool.Lock()
	defer sqlPool.Unlock()

	sqlPool.confs[name] = c

	db, err := NewMySQL(c)
	if err != nil {
		return err
	}

	sqlPool.clients[name] = db

	return nil
}

// MySQLClient get mysql client
func MySQLClient(name string) (*DB, error) {
	sqlPool.RLock()

	if v, ok := sqlPool.clients[name]; ok {
		sqlPool.RUnlock()
		return v.(*DB), nil
	}

	sqlPool.RUnlock()

	c, ok := sqlPool.confs[name]
	if !ok {
		return nil, errors.New("db  need init")
	}

	sqlPool.Lock()

	if v, ok := sqlPool.clients[name]; ok {
		sqlPool.Unlock()
		return v.(*DB), nil
	}

	db, err := NewMySQL(c.(*Config))
	if err != nil {
		return nil, err
	}

	sqlPool.clients[name] = db

	sqlPool.Unlock()

	return db, nil
}

// CloseMySQL close all my sql conn
func CloseMySQL() error {
	sqlPool.Lock()

	for _, v := range sqlPool.clients {
		v.(*DB).Close()
	}

	sqlPool.Unlock()

	return nil
}

// HealthCheckMySQL TODO
func HealthCheckMySQL() {
	sqlPool.RLock()

	if len(sqlPool.confs) != len(sqlPool.clients) {
		// unhealth
	}

	for k, v := range sqlPool.clients {
		err := v.(*DB).Ping(context.Background())
		if err != nil {
			_ = k
			return
		}
	}
}
