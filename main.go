package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/torderonex/crudgenerator/internal/generator"
	"github.com/torderonex/crudgenerator/internal/parser"
)

func getModuleName() (string, error) {
	cmd := exec.Command("go", "list", "-m")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	moduleName := strings.TrimSpace(string(output))
	return moduleName, nil
}

func createCrudFile(s parser.StructWithFields, outFilePath string) error {
	module, _ := getModuleName()

	//generator init

	gen := generator.CodeGenerator{Struct: s, Module: module}

	if !strings.HasSuffix(outFilePath, ".go") {
		outFilePath += ".go"
	}

	file, err := os.Create(outFilePath)
	if err != nil {
		return fmt.Errorf("Error creating file: %v", err)
	}
	defer file.Close()

	file.Write([]byte(gen.GenerateCRUD()))

	return nil
}

func main() {

	in := flag.String("i", "", "path to input file")
	out := flag.String("o", "", "path to output directory")
	flag.Parse()
	if *in == "" || *out == "" {
		log.Fatal("Wrong command line args")
	}

	dir, err := os.Getwd()
	err = os.Chdir(dir)

	if !strings.HasPrefix(*in, dir) {
		*in = dir + *in
	}
	if !strings.HasPrefix(*out, dir) {
		*out = dir + *out
	}

	_, err = exec.Command("cmd", "/c", "go get github.com/jmoiron/sqlx").Output()
	if err != nil {
		log.Fatal("Ошибка при скачивании пакета github.com/jmoiron/sqlx:", err)
	}

	structs, err := parser.ParseGoFile(*in)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range structs {
		o := path.Join(*out, strings.ToLower(v.Name))
		createCrudFile(v, o)
	}
}
