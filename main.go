package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	sql2 "github.com/minitauros/swagen/sql"
	"github.com/minitauros/swagen/swagger"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DB struct {
		DSN string
	}
	Service   swagger.ServiceInfo
	Resources map[string]swagger.Resource // Table name => resource
}

func main() {
	configFilePath := flag.String("conf", "", "config file path")

	flag.Parse()

	if *configFilePath == "" {
		log.Fatal("No config file path given. Use the -conf flag.")
	}

	configBytes, err := ioutil.ReadFile(*configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	conf := Config{}
	err = yaml.Unmarshal(configBytes, &conf)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("mysql", conf.DB.DSN)
	if err != nil {
		log.Fatal(err)
	}

	generator := swagger.Generator{
		TableService: sql2.NewTableService(db),
		Resources:    conf.Resources,
		ServiceInfo:  conf.Service,
	}

	swag, err := generator.Generate()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(swag)
}
