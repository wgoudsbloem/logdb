package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestPut(t *testing.T) {
	k := "testPut"
	v := 99
	i64, err := put(k, v)
	defer os.Remove(k)
	if err != nil {
		t.Error(err)
	}
	if i64 < 0 {
		t.Errorf("\nwant:\t%+v\ngot:\t%+v", 0, i64)
	}
	b, err := ioutil.ReadFile(k)
	if err != nil {
		t.Error(err)
	}
	got := string(b)
	want := fmt.Sprintf("%v\n", v)
	if want != got {
		t.Errorf("\nwant:\t'%+v'\ngot:\t'%+v'", want, got)
	}
}

func TestCliPut(t *testing.T) {
	k := "testCliPut"
	v := 99
	in := fmt.Sprintf(`put={"key":"%s", "value":%d}`, k, v)
	t.Log(in)
	var r, w bytes.Buffer
	r.WriteString(in)
	err := cmd(&r, &w)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(k)
	b, _, err := bufio.NewReader(&w).ReadLine()
	if err != nil {
		t.Error(err)
	}
	want := `ok {"key":"position","value":0}`
	got := string(b)
	if err != nil && err != io.EOF {
		t.Error(err)
	}
	if want != got {
		t.Errorf("\nwant:\t'%+v'\ngot:\t'%+v'", want, got)
	}
	// test offset position
	var r2, w2 bytes.Buffer
	r2.WriteString(in)
	err = cmd(&r2, &w2)
	if err != nil {
		t.Error(err)
	}
	b2, _, err := bufio.NewReader(&w2).ReadLine()
	if err != nil {
		t.Error(err)
	}
	want2 := `ok {"key":"position","value":3}`
	got2 := string(b2)
	if err != nil && err != io.EOF {
		t.Error(err)
	}
	if want2 != got2 {
		t.Errorf("\nwant:\t'%+v'\ngot:\t'%+v'", want2, got2)
	}
}

func TestSum(t *testing.T) {
	var r bytes.Buffer
	r.WriteString("1\n")
	r.WriteString("2\n")
	r.WriteString("3\n")
	got, err := sum(&r)
	if err != nil {
		t.Error(err)
	}
	want := 6
	if want != got {
		t.Errorf("\nwant:\t'%+v'\ngot:\t'%+v'", want, got)
	}
}

func TestGetSum(t *testing.T) {
	// start add values
	k := "testGetSum"
	v1 := 40
	_, err := put(k, v1)
	if err != nil {
		t.Error(err)
	}
	v2 := 60
	_, err = put(k, v2)
	// remove test file
	defer os.Remove(k)
	if err != nil {
		t.Error(err)
	}
	// end add values
	got, err := get(k, sum)
	if err != nil {
		t.Error(err)
	}
	want := 100
	if want != got {
		t.Errorf("\nwant:\t'%+v'\ngot:\t'%+v'", want, got)
	}
}

func TestCliGetSum(t *testing.T) {
	k := "testCliGetSum"
	// start add values
	v1 := 110
	in1 := fmt.Sprintf(`put={"key":"%s", "value":%d}`, k, v1)
	t.Log(in1)
	var r1, w1 bytes.Buffer
	r1.WriteString(in1)
	cmd(&r1, &w1)
	v2 := 90
	in2 := fmt.Sprintf(`put={"key":"%s", "value":%d}`, k, v2)
	t.Log(in2)
	var r2, w2 bytes.Buffer
	r2.WriteString(in2)
	cmd(&r2, &w2)
	// end add values
	in := fmt.Sprintf(`get={"key":"%s", "value":"sum"}`, k)
	t.Log(in)
	var r3, w3 bytes.Buffer
	r3.WriteString(in)
	err := cmd(&r3, &w3)
	// remove test file
	defer os.Remove(k)
	if err != nil {
		t.Error(err)
	}
	b, _, err := bufio.NewReader(&w3).ReadLine()
	if err != nil {
		t.Error(err)
	}
	want := `ok {"key":"sum","value":200}`
	got := string(b)
	if err != nil && err != io.EOF {
		t.Error(err)
	}
	if want != got {
		t.Errorf("\nwant:\t'%+v'\ngot:\t'%+v'", want, got)
	}
}

func TestLast(t *testing.T) {
	var r bytes.Buffer
	r.WriteString("4\n")
	r.WriteString("5\n")
	r.WriteString("6\n")
	got, err := last(&r)
	if err != nil {
		t.Error(err)
	}
	want := 6
	if want != got {
		t.Errorf("\nwant:\t'%+v'\ngot:\t'%+v'", want, got)
	}
}

func TestGetLast(t *testing.T) {
	// start add values
	k := "testGetLast"
	v1 := 22
	_, err := put(k, v1)
	if err != nil {
		t.Error(err)
	}
	// remove test file
	defer os.Remove(k)
	v2 := 11
	_, err = put(k, v2)
	if err != nil {
		t.Error(err)
	}
	// end add values
	got, err := get(k, last)
	if err != nil {
		t.Error(err)
	}
	want := v2
	if want != got {
		t.Errorf("\nwant:\t'%+v'\ngot:\t'%+v'", want, got)
	}
}

func TestCliGetLast(t *testing.T) {
	k := "testCLiGetLast"
	// start add values
	v1 := 222
	in1 := fmt.Sprintf(`put={"key":"%s", "value":%d}`, k, v1)
	var r1, w1 bytes.Buffer
	r1.WriteString(in1)
	cmd(&r1, &w1)
	v2 := 111
	in2 := fmt.Sprintf(`put={"key":"%s", "value":%d}`, k, v2)
	var r2, w2 bytes.Buffer
	r2.WriteString(in2)
	err := cmd(&r2, &w2)
	// remove test file
	defer os.Remove(k)
	if err != nil {
		t.Error(err)
	}
	// end add values
	in := fmt.Sprintf(`get={"key":"%s", "value":"last"}`, k)
	var r3, w3 bytes.Buffer
	r3.WriteString(in)
	cmd(&r3, &w3)
	b, _, err := bufio.NewReader(&w3).ReadLine()
	if err != nil {
		t.Error(err)
	}
	want := `ok {"key":"last","value":111}`
	got := string(b)
	if err != nil && err != io.EOF {
		t.Error(err)
	}
	if want != got {
		t.Errorf("\nwant:\t'%+v'\ngot:\t'%+v'", want, got)
	}
}

func TestValidate(t *testing.T) {
	in := []byte(`put={"key":"somekey", "value":99}`)
	want := true
	got := validate(in)
	if want != got {
		t.Errorf("\nwant:\t'%+v'\ngot:\t'%+v'", want, got)
	}
	in2 := []byte(`put={"key":"somekey", "value":"99"`)
	want2 := false
	got2 := validate(in2)
	if want2 != got2 {
		t.Errorf("\nwant:\t'%+v'\ngot:\t'%+v'", want2, got2)
	}
}
