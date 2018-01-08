package flow

import (
	"context"
	"strings"

	"gitee.com/antlinker/flow/util"
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
	if err := doc.ReadFromBytes(content); err != nil {
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
	return &node, nil
}

func (p *xmlParser) ParsesequenceFlow(element *etree.Element) (*sequenceFlow, error) {
	hasExpression := false
	var sequenceFlow sequenceFlow
	sequenceFlow.XMLName = element.Tag
	sequenceFlow.Code = element.SelectAttr("id").Value
	sequenceFlow.SourceRef = element.SelectAttr("sourceRef").Value
	sequenceFlow.TargetRef = element.SelectAttr("targetRef").Value
	for _, element := range element.ChildElements() {
		if element.Tag == "documentation" {
			sequenceFlow.Explain = element.Text()
		} else if element.Tag == "conditionExpression" {
			sequenceFlow.Expression = element.Text()
			hasExpression = true
		}
	}
	if !hasExpression {
		sequenceFlow.Expression = ""
	}
	return &sequenceFlow, nil
}

type nodeInfo struct {
	ProcessCode    string
	Type           string
	Code           string
	Name           string
	CandidateUsers []string
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
