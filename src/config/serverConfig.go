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
	Name           string `yaml:"name"`
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	StartTime      int    `yaml:"startTime"`
	Status         int    `yaml:"status"`
	Role           int    `yaml:"role"`
	AddTime        int    `yaml:"addTime"`
	LastModifyTime int    `yaml:"lastModifyTime"`
}

type Node struct {
	Name           string     `yaml:"name"`
	Index          int        `yaml:"index"`
	Start          int64      `yaml:"start"`
	End            int64      `yaml:"end"`
	AddTime        int        `yaml:"addTime"`
	LastModifyTime int        `yaml:"lastModifyTime"`
	ReadStrategy   int        `yaml:"readStrategy"`
	WriteStrategy  int        `yaml:"writeStrategy"`
	Master         Instance   `yaml:"master"`
	Slaves         []Instance `yaml:"slaves"`
}

type Cluster struct {
	ClusterName    string `yaml:"clusterName"`
	HashAlg        string `yaml:"hashAlg"`
	Nodes          []Node `yaml:"nodes"`
	AddTime        int    `yaml:"addTime"`
	LastModifyTime int    `yaml:"lastModifyTime"`
}

type ServerConfig struct {
	Listen  string `yaml:"listen"`
	Port    int    `yaml:"port"`
	Cl      Cluster
	LogFile string
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

func DumpConfig(path string, aa *ServerConfig) {
	data, _ := yaml.Marshal(aa)
	fmt.Println(fmt.Sprintf("the data len is %d\n", len(data)))
	ioutil.WriteFile(path, data, 0777)
}

func GenTestConfig() *ServerConfig {
	ins := new(Instance)
	ins.Host = "127.0.0.1"
	ins.Port = 3306
	ins.Name = "test-instance"
	ins.Role = MASTER
	ins.Status = ALIVE

	n := new(Node)
	n.Index = 0
	n.Start = 0
	n.End = 1000000
	n.Name = "test-node"
	n.ReadStrategy = 1
	n.WriteStrategy = 2

	n.Master = *ins
	n.Slaves = []Instance{*ins}

	c := new(Cluster)
	c.ClusterName = "test-Cl"
	c.HashAlg = "region"
	c.Nodes = []Node{*n}

	config := new(ServerConfig)
	config.Listen = "0.0.0.0"
	config.Port = 9999
	config.LogFile = "/tmp/server.log"
	config.Cl = *c
	return config
}
