package flow

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

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

	p := NewXMLParser()
	v, err := p.Parse(context.Background(), data)
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

	p := NewXMLParser()
	v, err := p.Parse(nil, data)
	if err != nil {
		fmt.Println(err.Error())
	}
	buf, _ := json.Marshal(v)
	fmt.Println(string(buf))
}

func TestParseBpmnForm(t *testing.T) {
	file, err := os.Open("test_data/form_test.bpmn") // For read access.
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

	p := NewXMLParser()
	v, err := p.Parse(nil, data)
	if err != nil {
		fmt.Println(err.Error())
	}
	buf, _ := json.Marshal(v)
	fmt.Println(string(buf))
}
