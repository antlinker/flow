package flow

import (
	"context"

	"github.com/beevik/etree"
)

// NewIonicRenderer 渲染到ionic
func NewIonicRenderer() Renderer {
	return &ionicRenderer{}
}

type ionicRenderer struct {
}

func (p *ionicRenderer) Render(context.Context, *NodeFormResult) ([]byte, error) {
	// content := page.CreateElement("ion-content")

	return nil, nil
}

func (p *ionicRenderer) getCreateFunctionByFieldType(field *FormFieldResult) (func(*etree.Element, *FormFieldResult), error) {
	switch field.Type {
	case "string":
		return p.createInput, nil
		// case ID:
		// 	return createInput, nil
		// case Image:
		// 	return createImageInput, nil
		// case Select:
		// 	return createSelect, nil
		// case Date:
		// 	return createDateInput, nil
		// case Dynamic:
		// 	return createDynamic, nil
		// case Textarea:
		// 	return createTextarea, nil
		// }
	}
	return nil, nil
}

func (p *ionicRenderer) createInput(form *etree.Element, field *FormFieldResult) {
	// item := form.CreateElement("ion-item")
	// label := item.CreateElement("ion-label")
	// label.SetText(field.FieldName)
	// //p.SetText(field)
	// input := item.CreateElement("ion-input")
	// // input.CreateAttr("id", field.FieldName)
	// input.CreateAttr("name", field.FieldName)
}

// func createSelect(form *etree.Element, field *Field) {

// 	item := form.CreateElement("ion-item")
// 	label := item.CreateElement("ion-label")
// 	label.SetText(field.FieldName)

// 	selectElement := item.CreateElement("ion-select")
// 	// selectElement.CreateAttr("id", field.FieldName)
// 	selectElement.CreateAttr("name", field.FieldName)
// 	for _, option := range field.FieldValues {
// 		optionElement := selectElement.CreateElement("ion-select-option")
// 		optionElement.SetText(option)
// 		// optionElement.CreateAttr("id", option)
// 		optionElement.CreateAttr("value", option)
// 	}
// }
// func createTextarea(form *etree.Element, field *Field) {
// 	item := form.CreateElement("ion-item")
// 	textarea := item.CreateElement("ion-textarea")
// 	// textarea.CreateAttr("id", field.FieldName)
// 	textarea.CreateAttr("name", field.FieldName)
// 	textarea.SetText(field.FieldName)
// 	form.CreateElement("br")
// }
// func createImageInput(form *etree.Element, field *Field) {
// 	item := form.CreateElement("ion-item")
// 	label := item.CreateElement("ion-label")
// 	label.SetText(field.FieldName)
// 	input := item.CreateElement("ion-input")
// 	// input.CreateAttr("id", field.FieldName)
// 	input.CreateAttr("name", field.FieldName)
// 	input.CreateAttr("type", "file")
// }

// func createDateInput(form *etree.Element, field *Field) {
// 	item := form.CreateElement("ion-item")
// 	label := item.CreateElement("ion-label")
// 	label.SetText(field.FieldName)
// 	input := item.CreateElement("ion-datetime")
// 	// input.CreateAttr("id", field.FieldName)
// 	input.CreateAttr("name", field.FieldName)
// 	input.CreateAttr("display-format", "MM DD YY")
// 	//input.CreateAttr("type", "date")
// }
