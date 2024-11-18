package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const (
	START = iota
	ALIVE
	ISOLATED
	DEAD
)

const (
	MASTER = iota
	SLAVE
)

type Instance struct {
	name           string
	host           string
	port           int
	startTime      int
	status         int
	role           int
	addTime        int
	lastModifyTime int
}

type Node struct {
	name           string
	master         Instance
	slaves         []Instance
	index          int
	start          int64
	end            int64
	addTime        int
	lastModifyTime int
	readStrategy   int
	writeStrategy  int
}

type Cluster struct {
	clusterName    string
	hashAlg        string
	nodes          []Node
	addTime        int
	lastModifyTime int
}

type ServerConfig struct {
	listen  string
	port    int
	cluster Cluster
	logFile string
}

func LoadConfig(path string) (*ServerConfig, error) {
	var data []byte
	var err error
	if data, err = ioutil.ReadFile(path); err != nil {
		return nil, err
	}

	var config *ServerConfig
	if err = yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}
	return config, nil
}

func DumpConfig(aa *ServerConfig) {
	data, _ := yaml.Marshal(aa)
	fmt.Println(fmt.Sprintf("the data len is %d\n", len(data)))
	ioutil.WriteFile("/tmp/test.yaml", data, 0777)
}
