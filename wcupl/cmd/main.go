package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"log"

	"github.com/aymerick/raymond"
)

const layoutISO = "2006-01-02"

var (
	cuplargs []string
	file     string
	headers  []string
	pld      []string
	simu     []string
	mode     string
	hasSimu  bool
	cuplHome string
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

func parseParameters() {
	ch, ok := os.LookupEnv("CUPL_HOME")
	if !ok {
		ch = os.Args[0]
		ch = strings.TrimSuffix(ch, "wcupl.exe")
	}

	cuplHome = ch

	args := os.Args[1:]
	for _, arg := range args {
		log.Println("arg: ", arg)
		if strings.HasSuffix(arg, ".wpld") {
			file = arg
		} else {
			cuplargs = append(cuplargs, arg)
		}
	}
}

func main() {
	log.Println("start")

	parseParameters()

	log.Println("cupl args:", cuplargs)
	log.Println("file     :", file)

	workdir := filepath.Join(filepath.Dir(file), "wcupl")
	err := os.Mkdir(workdir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Fatalf("error creating temp dir: %v\n", err)
	}

	destfile := filepath.Join(workdir, filepath.Base(file))

	file, err = renderFile(file, destfile)
	if err != nil {
		log.Fatalf("error rendering file: %v\n", err)
	}

	err, headers, pld, simu := scanfile(file)
	if err != nil {
		log.Fatalf("error scanning file: %v\n", err)
	}

	hasSimu, err := validate(headers, pld, simu)
	if err != nil {
		log.Fatalf("error validating file: %v\n", err)
	}

	log.Println("starting file generation")
	pldFile := strings.TrimSuffix(file, filepath.Ext(file)) + ".pld"
	err = generateFile(pldFile, headers, pld)
	if err != nil {
		log.Fatalf("error creating file: %v\n", err)
	}
	if hasSimu {
		siFile := strings.TrimSuffix(file, filepath.Ext(file)) + ".si"
		err = generateFile(siFile, headers, simu)
		if err != nil {
			log.Fatalf("error creating file: %v\n", err)
		}
	}

	log.Println("starting cupl from ", cuplHome)
	err = startcupl(cuplargs, pldFile)
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

func renderFile(file, destfile string) (string, error) {
	filename := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
	ctx := map[string]string{
		"filename": filename,
		"date":     time.Now().Format(layoutISO),
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return destfile, fmt.Errorf("error loading source file: %v\n", err)
	}

	result, err := raymond.Render(string(b), ctx)
	if err != nil {
		return destfile, fmt.Errorf("error render file: %v\n", err)
	}

	err = ioutil.WriteFile(destfile, []byte(result), 0644)
	if err != nil {
		return destfile, fmt.Errorf("error writing file: %v\n", err)
	}
	return destfile, nil
}

func scanfile(file string) (err error, headers, pld, simu []string) {

	headers = make([]string, 0)
	pld = make([]string, 0)
	simu = make([]string, 0)

	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("error opening file: %v\n", err), nil, nil, nil
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(strings.ToLower(line), "header:") {
			mode = "hdr"
			continue
		}
		if strings.HasPrefix(strings.ToLower(line), "pld:") {
			mode = "pld"
			continue
		}
		if strings.HasPrefix(strings.ToLower(line), "simulator:") {
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

	return
}

func validate(headers, pld, simu []string) (hasSimu bool, err error) {
	if len(headers) == 0 {
		err = errors.New("no headers found")
	}
	if len(pld) == 0 {
		err = errors.New("no pld found")
	}
	hasSimu = true
	if len(simu) == 0 {
		log.Print("no simulation found")
		hasSimu = false
	}
	return
}

func generateFile(file string, headers, part []string) error {
	content := append(headers, part...)
	return WriteLines(file, content)
}

func startcupl(cuplargs []string, pldFile string) error {
	cuplargs = append(cuplargs, pldFile)
	cupl := filepath.Join(cuplHome, "cupl.exe")
	cmnd := exec.Command(cupl, cuplargs...)
	var outb, errb bytes.Buffer
	cmnd.Stdout = &outb
	cmnd.Stderr = &errb
	err := cmnd.Run()
	log.Println(outb.String())
	log.Println(errb.String())

	return err
}
