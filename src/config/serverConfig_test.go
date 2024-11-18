package config

import (
	"testing"
)

func genTestConfig() *ServerConfig {
	ins := new(Instance)
	ins.host = "127.0.0.1"
	ins.port = 3306
	ins.name = "test-instance"
	ins.role = MASTER
	ins.status = ALIVE

	n := new(Node)
	n.index = 0
	n.start = 0
	n.end = 1000000
	n.name = "test-node"
	n.readStrategy = 1
	n.writeStrategy = 2

	n.master = *ins
	n.slaves = []Instance{*ins}

	c := new(Cluster)
	c.clusterName = "test-cluster"
	c.hashAlg = "region"
	c.nodes = []Node{*n}

	config := new(ServerConfig)
	config.listen = "0.0.0.0"
	config.port = 9999
	config.logFile = "/tmp/server.log"
	config.cluster = *c
	return config
}

func TestServerConfig_DumpConfig(t *testing.T) {
	s := genTestConfig()
	DumpConfig(s)
}

func TestServerConfig_LoadConfig(t *testing.T) {

}

func Test_genTestConfig(t *testing.T) {

}
