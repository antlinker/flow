package flow

import (
	"context"
	"fmt"

	"github.com/beevik/etree"
)

const ionicURL = "http://119.27.161.206:8080/ionic.js"

// NewIonicRenderer 渲染到ionic
func NewIonicRenderer() Renderer {
	return &ionicRenderer{}
}

type ionicRenderer struct {
}

func (p *ionicRenderer) Render(ctx context.Context, form *NodeFormResult) ([]byte, error) {
	// content := page.CreateElement("ion-content")
	doc := etree.NewDocument()
	doc.WriteSettings.CanonicalText = true
	// doc.WriteSettings.
	list := p.createFormTemplate(&doc.Element)
	for _, field := range form.Fields {
		createFunc, _ := p.getCreateFunctionByFieldType(field)

		if createFunc != nil {
			createFunc(list, field)
			fmt.Println("domfffffffffffffff", list)
		}
	}
	// var buf = make([]byte, 0)
	doc.Indent(2)

	buf, err := doc.WriteToBytes()
	return buf, err
}

func (p *ionicRenderer) createFormTemplate(doc *etree.Element) *etree.Element {

	body := doc.CreateElement("body")
	viewport := body.CreateElement("meta")
	viewport.CreateAttr("name", "viewport")
	viewport.CreateAttr("content", "width=device-width,minimum-scale=1.0,maximum-scale=1.0,user-scalable=no")

	ionicScript := body.CreateElement("script")
	ionicScript.SetText("")
	ionicScript.CreateAttr("src", ionicURL)

	app := body.CreateElement("ion-app")
	page := app.CreateElement("ion-page")
	content := page.CreateElement("ion-content")
	// form := content.CreateElement("form")
	// form.CreateAttr("action", "/api/form")
	// form.CreateAttr("method", "post")
	list := content.CreateElement("ion-list")
	return list
}

func (p *ionicRenderer) getCreateFunctionByFieldType(field *FormFieldResult) (func(*etree.Element, *FormFieldResult), error) {
	switch field.Type {
	case "string":
		return p.createStringInput, nil
	case "long":
		return p.createLongInput, nil
	case "date":
		return p.createDateInput, nil
	case "enum":
		return p.createSelectInput, nil
	case "boolean":
		return p.createRatioInput, nil
	}
	return nil, nil
}

func (p *ionicRenderer) createStringInput(form *etree.Element, field *FormFieldResult) {
	item := form.CreateElement("ion-item")
	label := item.CreateElement("ion-label")
	label.SetText(field.Label)
	input := item.CreateElement("ion-input")
	input.CreateAttr("id", field.ID)
	// input.SetText(field.DefaultValue)
}

func (p *ionicRenderer) createLongInput(form *etree.Element, field *FormFieldResult) {
	item := form.CreateElement("ion-item")
	label := item.CreateElement("ion-label")
	label.SetText(field.Label)
	input := item.CreateElement("ion-input")
	input.CreateAttr("id", field.ID)
	// input.SetText(field.DefaultValue)
}

func (p *ionicRenderer) createDateInput(form *etree.Element, field *FormFieldResult) {
	item := form.CreateElement("ion-item")
	label := item.CreateElement("ion-label")
	label.SetText(field.Label)
	input := item.CreateElement("ion-datetime")
	input.CreateAttr("id", field.ID)
	// input.SetText(field.DefaultValue)
}

func (p *ionicRenderer) createSelectInput(form *etree.Element, field *FormFieldResult) {
	item := form.CreateElement("ion-item")
	label := item.CreateElement("ion-label")
	label.SetText(field.Label)
	input := item.CreateElement("ion-select")
	input.CreateAttr("id", field.ID)

	for _, option := range field.Values {
		optionElement := input.CreateElement("ion-select-option")
		optionElement.SetText(option.Name)
		// optionElement.CreateAttr("id", option)
		optionElement.CreateAttr("id", option.ID)
	}
}

func (p *ionicRenderer) createRatioInput(form *etree.Element, field *FormFieldResult) {
	item := form.CreateElement("ion-item")
	label := item.CreateElement("ion-label")
	label.SetText(field.Label)
	input := item.CreateElement("ion-radio")
	input.CreateAttr("id", field.ID)
	// input.SetText(field.DefaultValue)
}
