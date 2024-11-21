package mysql

import (
	"testing"
)

func TestBackendConn_Connect(t *testing.T) {
	c := NewBackendConn()
	c.Connect("127.0.0.1", 3306, "root", "taotaoJJ1986@", "test")
}
