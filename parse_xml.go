package flow

import (
	"context"
	"strings"

	"ant-flow/util"

	"github.com/beevik/etree"
)

// NewXMLParser xml解析器
func NewXMLParser() Parser {
	return &xmlParser{}
}

type xmlParser struct {
}

func (p *xmlParser) Parse(ctx context.Context, content []byte) (*ParseResult, error) {
	result := &ParseResult{}
	var err error

	doc := etree.NewDocument()
	if err = doc.ReadFromBytes(content); err != nil {
		panic(err)
	}

	root := doc.SelectElement("definitions")
	//root.Tag
	process := root.SelectElement("process")

	if id := process.SelectAttr("id"); id != nil {
		result.FlowID = id.Value
	}
	if name := process.SelectAttr("name"); name != nil {
		result.FlowName = name.Value
	}
	if version := process.SelectAttr("versionTag"); version != nil {
		result.FlowVersion, err = util.StringToInt(version.Value)
		if err != nil {
			return nil, err
		}
	}

	// 定义一个用于辅助的map，由节点id映射到noderesult
	nodeMap := make(map[string]*NodeResult)
	// 遍历找到所有的节点，因为是解析一个树，所以先解析节点，再解析sequenceFlow部分
	// 解析sequenceFlow部分时，nodeMap里面应该已经有对应的nodeId了
	for _, element := range process.ChildElements() {
		if element.Tag == "documentation" ||
			element.Tag == "extensionElements" ||
			element.Tag == "sequenceFlow" {
			continue
		}
		node, _ := p.ParseNode(element)
		var nodeResult NodeResult
		nodeResult.NodeID = node.Code
		nodeResult.NodeName = node.Name
		nodeResult.NodeType, err = GetNodeTypeByName(node.Type)
		if err != nil {
			return nil, err
		}
		nodeResult.CandidateExpressions = node.CandidateUsers
		// yupengfei 2018-01-17 增加了form的解析
		nodeResult.FormResult = node.FormResult
		nodeResult.Properties = node.Properties
		nodeMap[nodeResult.NodeID] = &nodeResult
		// 如果节点是一个路由的话，需要特殊处理
	}

	for _, element := range process.ChildElements() {
		if element.Tag == "sequenceFlow" {
			sequenceFlow, _ := p.ParsesequenceFlow(element)
			var routerResult RouterResult
			routerResult.Expression = sequenceFlow.Expression
			routerResult.Explain = sequenceFlow.Explain
			routerResult.TargetNodeID = sequenceFlow.TargetRef
			if nodeResult, exist := nodeMap[sequenceFlow.SourceRef]; exist {
				nodeResult.Routers = append(nodeResult.Routers, &routerResult)
			}
		}
	}

	for _, nodeResult := range nodeMap {
		result.Nodes = append(result.Nodes, nodeResult)
	}
	return result, nil
}

func (p *xmlParser) ParseNode(element *etree.Element) (*nodeInfo, error) {
	var node nodeInfo

	node.Type = element.Tag
	if node.Type == "endEvent" {
		for _, e := range element.ChildElements() {
			if e.Tag == "terminateEventDefinition" {
				node.Type = "terminateEvent"
			}
		}
	}
	if name := element.SelectAttr("name"); name != nil {
		node.Name = name.Value
	}
	if id := element.SelectAttr("id"); id != nil {
		node.Code = id.Value
	}
	if candidateUsers := element.SelectAttr("candidateUsers"); candidateUsers != nil {
		candidateUserList := strings.Split(candidateUsers.Value, ";")
		node.CandidateUsers = candidateUserList
	}

	if extensionElements := element.SelectElement("extensionElements"); extensionElements != nil {
		if formData := extensionElements.SelectElement("formData"); formData != nil {
			form, err := p.ParseFormData(formData)
			if err != nil {
				return nil, err
			}
			if form != nil {
				if formKey := element.SelectAttr("formKey"); formKey != nil {
					form.ID = formKey.Value
				}
			}
			node.FormResult = form
		}

		if propertyData := extensionElements.SelectElement("properties"); propertyData != nil {
			// 解析节点属性
			for _, p := range propertyData.SelectElements("property") {
				var item PropertyResult
				if name := p.SelectAttr("name"); name != nil {
					item.Name = name.Value
				}
				if value := p.SelectAttr("value"); value != nil {
					item.Value = value.Value
				}

				if item.Name != "" {
					node.Properties = append(node.Properties, &item)
				}
			}
		}
	}

	return &node, nil
}

func (p *xmlParser) ParsesequenceFlow(element *etree.Element) (*sequenceFlow, error) {
	hasExpression := false
	var seq sequenceFlow
	seq.XMLName = element.Tag
	seq.Code = element.SelectAttr("id").Value
	seq.SourceRef = element.SelectAttr("sourceRef").Value
	seq.TargetRef = element.SelectAttr("targetRef").Value
	for _, element := range element.ChildElements() {
		if element.Tag == "documentation" {
			seq.Explain = element.Text()
		} else if element.Tag == "conditionExpression" {
			seq.Expression = element.Text()
			hasExpression = true
		}
	}
	if !hasExpression {
		seq.Expression = ""
	}
	return &seq, nil
}

func (p *xmlParser) ParseFormData(element *etree.Element) (*NodeFormResult, error) {
	var formResult = &NodeFormResult{}
	if id := element.SelectAttr("id"); id != nil {
		formResult.ID = id.Value
	}

	if fieldList := element.SelectElements("formField"); fieldList != nil {
		for _, item := range fieldList {
			var field = &FormFieldResult{}
			var err error
			if properties := item.SelectElement("properties"); properties != nil {
				field.Properties, err = p.ParseProperties(properties)
				if err != nil {
					return nil, err
				}
			}
			if validations := item.SelectElement("validation"); validations != nil {
				field.Validations, err = p.ParseValidations(validations)
				if err != nil {
					return nil, err
				}
			}
			if nodeType := item.SelectAttr("type"); nodeType != nil {
				field.Type = nodeType.Value
				if field.Type == "enum" {
					field.Values, err = p.ParseEnumValues(item)
					if err != nil {
						return nil, err
					}
				}
			}
			if id := item.SelectAttr("id"); id != nil {
				field.ID = id.Value
			}
			if label := item.SelectAttr("label"); label != nil {
				field.Label = label.Value
			}
			if defaultValue := item.SelectAttr("defaultValue"); defaultValue != nil {
				field.DefaultValue = defaultValue.Value
			}
			formResult.Fields = append(formResult.Fields, field)
		}
	}

	return formResult, nil
}

func (p *xmlParser) ParseProperties(element *etree.Element) ([]*FieldProperty, error) {
	var properties = make([]*FieldProperty, 0)
	if propertyList := element.SelectElements("property"); propertyList != nil {
		for _, item := range propertyList {
			var property = &FieldProperty{}
			if id := item.SelectAttr("id"); id != nil {
				property.ID = id.Value
			}
			if value := item.SelectAttr("value"); value != nil {
				property.Value = value.Value
			}
			properties = append(properties, property)
		}
	}
	return properties, nil
}

func (p *xmlParser) ParseValidations(element *etree.Element) ([]*FieldValidation, error) {
	var validations = make([]*FieldValidation, 0)
	if validationList := element.SelectElements("constraint"); validationList != nil {
		for _, item := range validationList {
			var validation = &FieldValidation{}
			if name := item.SelectAttr("name"); name != nil {
				validation.Name = name.Value
			}
			if config := item.SelectAttr("config"); config != nil {
				validation.Config = config.Value
			}
			validations = append(validations, validation)
		}
	}
	return validations, nil
}

func (p *xmlParser) ParseEnumValues(element *etree.Element) ([]*FieldOption, error) {
	var options = make([]*FieldOption, 0)
	if optionList := element.SelectElements("value"); optionList != nil {
		for _, item := range optionList {
			var option = &FieldOption{}
			if id := item.SelectAttr("id"); id != nil {
				option.ID = id.Value
			}
			if name := item.SelectAttr("name"); name != nil {
				option.Name = name.Value
			}
			options = append(options, option)
		}
	}
	return options, nil
}

type nodeInfo struct {
	ProcessCode    string
	Type           string
	Code           string
	Name           string
	CandidateUsers []string
	Properties     []*PropertyResult
	FormResult     *NodeFormResult
}

type sequenceFlow struct {
	ProcessCode string
	XMLName     string
	Code        string
	SourceRef   string
	TargetRef   string
	Explain     string
	Expression  string
}
