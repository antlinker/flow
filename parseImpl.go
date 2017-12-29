package flow

import (
	"context"
	"strings"

	"gitee.com/antlinker/flow/util"
	"github.com/beevik/etree"
)

type Node struct {
	ProcessCode string
	Type        string
	//Id int64  `gorm:"Id,primary_key,AUTO_INCREMENT"`
	Code           string
	Name           string
	CandidateUsers []string
}

type SequenceFlow struct {
	ProcessCode string
	XmlName     string
	Code        string
	SourceRef   string
	TargetRef   string
	Explain     string
	Expression  string
}
type ParserImpl struct {
}

func (ParserImpl) Parse(ctx context.Context, content []byte) (*ParseResult, error) {
	var result *ParseResult = &ParseResult{}
	var err error

	//time.Now()
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
		if element.Tag == "documentation" || element.Tag == "extensionElements" {
			continue
		} else if element.Tag == "sequenceFlow" {
			sequenceFlow, _ := ParseSequenceFlow(element)
			var routerResult RouterResult
			routerResult.Expression = sequenceFlow.Expression
			routerResult.Explain = sequenceFlow.Explain
			routerResult.TargetNodeID = sequenceFlow.TargetRef
			if nodeResult, exist := nodeMap[sequenceFlow.SourceRef]; exist {
				nodeResult.Routers = append(nodeResult.Routers, &routerResult)
			}
		} else {
			node, _ := ParseNode(element)
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
	}

	for _, nodeResult := range nodeMap {
		result.Nodes = append(result.Nodes, nodeResult)
	}
	return result, nil
}

func ParseNode(element *etree.Element) (*Node, error) {
	var node Node

	node.Type = element.Tag
	if name := element.SelectAttr("name"); name != nil {
		node.Name = name.Value
	}
	if id := element.SelectAttr("id"); id != nil {
		node.Code = id.Value
	}
	if candidateUsers := element.SelectAttr("candidateUsers"); candidateUsers != nil {
		candidateUserList := strings.Split(candidateUsers.Value, ";")
		node.CandidateUsers = candidateUserList
		// node = id.Value
	}
	return &node, nil
}

func ParseSequenceFlow(element *etree.Element) (*SequenceFlow, error) {
	hasExpression := false
	var sequenceFlow SequenceFlow
	sequenceFlow.XmlName = element.Tag
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
		sequenceFlow.Expression = "true"
	}
	return &sequenceFlow, nil
}
