package flow

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestRenderIonic(t *testing.T) {
	var testString = `
	{
		"ID": "key_test",
		"Fields": [{
			"ID": "id_1",
			"Type": "string",
			"Label": "label_1",
			"DefaultValue": "default_1",
			"Values": null,
			"Validations": [{
				"Name": "constraint_1",
				"Config": "constraint_1_config"
			}, {
				"Name": "constraint_2",
				"Config": "constraint_2_config"
			}],
			"Properties": [{
				"ID": "Property_1",
				"Value": "Property_1_value"
			}, {
				"ID": "Property_2",
				"Value": "Property_2_value"
			}]
		}, {
			"ID": "FormField_2",
			"Type": "long",
			"Label": "FormField_2_label",
			"DefaultValue": "",
			"Values": null,
			"Validations": null,
			"Properties": null
		}, {
			"ID": "FormField_3",
			"Type": "date",
			"Label": "",
			"DefaultValue": "",
			"Values": null,
			"Validations": null,
			"Properties": null
		}, {
			"ID": "FormField_4",
			"Type": "enum",
			"Label": "FormField_4_label",
			"DefaultValue": "FormField_4_label_default",
			"Values": [{
				"ID": "Value_1",
				"Name": "Value_1_value"
			}, {
				"ID": "Value_2",
				"Name": "Value_2_value"
			}],
			"Validations": null,
			"Properties": null
		}, {
			"ID": "FormField_5",
			"Type": "boolean",
			"Label": "",
			"DefaultValue": "",
			"Values": null,
			"Validations": null,
			"Properties": null
		}]
	}
	`
	var form = &NodeFormResult{}
	err := json.Unmarshal([]byte(testString), form)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(form.ID)
	render := NewIonicRenderer()
	result, err := render.Render(nil, form)
	if err != nil {
		t.Error(err.Error())
	}
	// fmt.Println()
	ioutil.WriteFile("../tmp", result, 0644)
}
