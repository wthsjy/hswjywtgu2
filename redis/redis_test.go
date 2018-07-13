package redis

import "testing"

func TestRegisterRedis(t *testing.T) {
	RegisterRedis("local", "127.0.0.1:6379:pwd")
}

func TestClient(t *testing.T) {
	Client("local")
}

func TestCloseAll(t *testing.T) {
	CloseAll()
}
