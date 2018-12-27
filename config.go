package main

type Mysql struct {
	Dsn string `yaml:"dsn"`
}

type Address struct {
	Ip   string `yaml:"ip"`
	Port int    `yaml:"port"`
}

// Config struct corresponding to config.yml structure
type Config struct {
	Debug bool    `yaml:"debug"`
	Mysql Mysql   `yaml:"mysql"`
	Http  Address `yaml:"http"`
}
