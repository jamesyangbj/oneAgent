package config

import (
	"fmt"
	"testing"
)

func TestServerConfig_DumpConfig(t *testing.T) {
	//s := GenTestConfig()
	//DumpConfig("/tmp/test.yaml", s)
}

func TestServerConfig_LoadConfig(t *testing.T) {
	s, _ := LoadConfig("/tmp/test.yaml")
	fmt.Println(s.Listen)
	//if s.Listen != "127.0.0.1" {
	//	t.Fail()
	//}

}

func Test_genTestConfig(t *testing.T) {

}
