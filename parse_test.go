package flow

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	var p Parser = ParserImpl{}
	p.Parse(nil, nil)
}

func TestParseBpmn(t *testing.T) {
	file, err := os.Open("test_data/basic.bpmn") // For read access.
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	var p Parser = ParserImpl{}
	v, err := p.Parse(nil, data)
	if err != nil {
		fmt.Println(err.Error())
	}
	buf, _ := json.Marshal(v)
	fmt.Println(string(buf))
}

func TestParseBpmn2(t *testing.T) {
	file, err := os.Open("test_data/route.bpmn") // For read access.
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	var p Parser = ParserImpl{}
	v, err := p.Parse(nil, data)
	if err != nil {
		fmt.Println(err.Error())
	}
	buf, _ := json.Marshal(v)
	fmt.Println(string(buf))
}
