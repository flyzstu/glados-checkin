package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Users []User `yaml:"user"`
}

type User struct {
	Email  string `yaml:"email"`
	Cookie string `yaml:"cookie"`
}

// 将yaml映射对象
func loadConf(confName string) *Config {
	b, err := ioutil.ReadFile(confName)
	if err != nil {
		panic(err)
	}
	conf := new(Config)
	err = yaml.Unmarshal(b, conf)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("conf: %v\n", conf)
	// fmt.Printf("conf.Users: %v\n", conf.Users)
	return conf

}
