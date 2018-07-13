package mysql

import (
	"testing"
	"time"
)

func TestRegisterMySQL(t *testing.T) {
	RegisterMySQL("debug", &Config{
		DSN:          "[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]",
		Active:       128,
		Idle:         128,
		IdleTimeout:  time.Minute * 5,
		QueryTimeout: time.Second * 10,
		ExecTimeout:  time.Second * 10,
		TranTimeout:  time.Second * 10,
	})
}
func TestMySQLClient(t *testing.T) {
	MySQLClient("debug")
}
func TestCloseMySQL(t *testing.T) {
	CloseMySQL()
}
