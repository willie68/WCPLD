package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"log"
)

var (
	cuplargs []string
	file     string
	headers  []string
	pld      []string
	simu     []string
	mode     string
	hasSimu  bool
)

func init() {
	log.Println("init wcupl")
	cuplargs = make([]string, 0)
	headers = make([]string, 0)
	pld = make([]string, 0)
	simu = make([]string, 0)
	mode = "nn"
	hasSimu = false
}

func main() {
	log.Println("start")

	cuplHome, ok := os.LookupEnv("CUPL_HOME")
	if !ok {
		cuplHome = os.Args[0]
		cuplHome = strings.TrimSuffix(cuplHome, "wcupl.exe")
	}
	args := os.Args[1:]
	for _, arg := range args {
		log.Println("arg: ", arg)
		if strings.HasSuffix(arg, ".wpld") {
			file = arg
		} else {
			cuplargs = append(cuplargs, arg)
		}
	}

	log.Println("cupl args:", cuplargs)
	log.Println("file     :", file)
	workdir := filepath.Join(filepath.Dir(file), "wcupl")
	err := os.Mkdir(workdir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Fatalf("error creating temp dir: %v\n", err)
	}
	destfile := filepath.Join(workdir, filepath.Base(file))
	file, err = copy(file, destfile)
	if err != nil {
		log.Fatalf("error copy source file: %v\n", err)
	}

	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("error opening file: %v\n", err)
	}
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(strings.ToLower(line), "header") {
			mode = "hdr"
			continue
		}
		if strings.HasPrefix(strings.ToLower(line), "pld") {
			mode = "pld"
			continue
		}
		if strings.HasPrefix(strings.ToLower(line), "simulator") {
			mode = "sim"
			continue
		}
		switch mode {
		case "hdr":
			headers = append(headers, line)
		case "pld":
			pld = append(pld, line)
		case "sim":
			simu = append(simu, line)
		}
	}

	if len(headers) == 0 {
		log.Fatal("no headers found")
	}
	if len(pld) == 0 {
		log.Fatal("no pld found")
	}
	hasSimu = true
	if len(simu) == 0 {
		log.Print("no simulation found")
		hasSimu = false
	}
	log.Println("starting file genration")
	pldFile := strings.TrimSuffix(file, filepath.Ext(file)) + ".pld"

	pld = append(headers, pld...)
	err = WriteLines(pldFile, pld)
	if err != nil {
		log.Fatalf("error creating file: %v\n", err)
	}
	if hasSimu {
		siFile := strings.TrimSuffix(file, filepath.Ext(file)) + ".si"
		simu = append(headers, simu...)
		err = WriteLines(siFile, simu)
		if err != nil {
			log.Fatalf("error creating file: %v\n", err)
		}
	}

	cuplargs = append(cuplargs, pldFile)
	log.Println("starting cupl from ", cuplHome)
	cupl := filepath.Join(cuplHome, "cupl.exe")
	cmnd := exec.Command(cupl, cuplargs...)
	var outb, errb bytes.Buffer
	cmnd.Stdout = &outb
	cmnd.Stderr = &errb
	err = cmnd.Run()
	log.Println(outb.String())
	log.Println(errb.String())
	if err != nil {
		log.Fatalf("error starting cupl: %v\n", err)
	}
	log.Println("finnished")
}

func WriteLines(path string, lines []string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func copy(src, dest string) (string, error) {
	bytesRead, err := ioutil.ReadFile(src)

	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(dest, bytesRead, 0644)

	if err != nil {
		return "", err
	}
	return dest, nil
}
