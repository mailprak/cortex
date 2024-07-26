package main

import (
	"io/ioutil"
	"log"
	"fmt"
	"gopkg.in/yaml.v2"
	"bytes"
	"os"
	"text/template"
)

type Fixes map[int]string

type Neuron struct {
	Name   string
	Path   string
	Fixes  Fixes
	Config interface{}
}

type Plan struct {
	ExitOnFirstError bool
	Parallel         []string
}

type AppConfig struct {
	Name       string
	Definition []Neuron
	Plan       Plan
}

func ReadConfig(filename string) (*AppConfig, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config AppConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func main() {
	config, err := ReadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// Generate the code based on the config
	GenerateCode(config)
}

func GenerateCode(config *AppConfig) {
	tmpl, err := template.ParseFiles("template.go.tmpl")
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, config)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}

	err = os.WriteFile("generated_code.go", buf.Bytes(), 0644)
	if err != nil {
		log.Fatalf("Error writing generated code: %v", err)
	}

	fmt.Println("Code generated successfully in generated_code.go")
}