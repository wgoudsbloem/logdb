package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	cmdPut = "put"
	cmdGet = "get"
	cmdOk  = "ok"
)

const (
	sumFunc  = "sum"
	lastFunc = "last"
)

const (
	helpText = "\nadd an entry:\t\t\tput={\"key\":\"<entry name>\", \"value\":<integer value>}\n" +
		"get an sum of the values:\tget{\"key\":\"<entry name>\", \"value\":\"sum\"}\n" +
		"get the last entry value:\tget{\"key\":\"<entry name>\", \"value\":\"last\"}\n\n"
)

func main() {
	fmt.Printf("welcome to logDB\n" + helpText)
	cmd(os.Stdin, os.Stdout)
}

type keyvalue struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// cmd is the parser for the command line instructions
func cmd(r io.Reader, w io.Writer) (err error) {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		if sc.Text() == "help" {
			fmt.Fprintln(w, helpText)
			continue
		}
		if !validate(sc.Bytes()) {
			w.Write([]byte("wrong format type help for instructions\n"))
			continue
		}
		s := strings.Split(sc.Text(), "=")
		var p keyvalue
		err = json.Unmarshal([]byte(s[1]), &p)
		if err != nil {
			break
		}
		var b []byte
		switch s[0] {
		case cmdPut:
			str, ok := p.Value.(string)
			if ok {
				p.Value, err = strconv.Atoi(str)
				if err != nil {
					break
				}
			}
			ret := keyvalue{Key: "position"}
			ret.Value, err = put(p.Key, p.Value)
			if err != nil {
				break
			}
			b, err = json.Marshal(ret)
			if err != nil {
				break
			}
			break
		case cmdGet:
			ret := keyvalue{Key: p.Value.(string)}
			switch p.Value.(string) {
			case sumFunc:
				ret.Value, err = get(p.Key, sum)
			case lastFunc:
				ret.Value, err = get(p.Key, last)
			}
			b, err = json.Marshal(ret)
		}
		if len(b) > 0 {
			_, err = w.Write([]byte(cmdOk + " "))
			_, err = w.Write(b)
			_, err = w.Write([]byte("\n"))
		}
	}
	if err != nil {
		fmt.Println(err)
	}
	return
}

// validates if a command follows the accepted pattern
func validate(b []byte) bool {
	reg := `[a-z]*={"key":"[A-Za-z]*",.?"value":["0-9a-z\-]*}`
	m, err := regexp.Match(reg, b)
	if err != nil {
		m = false
	}
	return m
}

// put enters the given val to the store with the key name
func put(key string, val interface{}) (pos int64, err error) {
	var f *os.File
	f, err = os.OpenFile(key, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	defer f.Close()
	if pos, err = f.Seek(0, os.SEEK_END); err != nil {
		return
	}
	s := fmt.Sprintf("%v", val)
	_, err = f.Write([]byte(s))
	_, err = f.Write([]byte("\n"))
	return
}

type act func(r io.Reader) (int, error)

// get retrieves the value computed by the given function act
func get(key string, fn act) (i int, err error) {
	var f *os.File
	f, err = os.OpenFile(key, os.O_RDONLY, 0666)
	if err != nil {
		return
	}
	i, err = fn(f)
	return
}

// sum is a convenience function to add up values
func sum(r io.Reader) (i int, err error) {
	rd := bufio.NewReader(r)
	for {
		var s string
		s, err = rd.ReadString('\n')
		if err != nil && err == io.EOF {
			err = nil
			break
		}
		var d int
		d, err = strconv.Atoi(s[:len(s)-1])
		if err != nil {
			break
		}
		i += d
	}
	return
}

// last is a convenience function to seek the last entered value
func last(r io.Reader) (i int, err error) {
	rd := bufio.NewReader(r)
	for {
		var s string
		s, err = rd.ReadString('\n')
		if err != nil && err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			break
		}
		i, err = strconv.Atoi(s[:len(s)-1])
	}
	return
}
