import React, {PureComponent} from 'react';
import {connect} from 'dva';
import {Card, Form, Layout} from 'antd';

var BpmnModeler = require('bpmn-js/lib/Modeler')
var propertiesPanelModule = require('bpmn-js-properties-panel')
var propertiesProviderModule = require('bpmn-js-properties-panel/lib/provider/camunda')
var camundaModdleDescriptor = require('camunda-bpmn-moddle/resources/camunda')

var newDiagramXML = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<bpmn2:definitions xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns:bpmn2=\"http://www.omg.org/spec/BPMN/20100524/MODEL\" xmlns:bpmndi=\"http://www.omg.org/spec/BPMN/20100524/DI\" xmlns:dc=\"http://www.omg.org/spec/DD/20100524/DC\" xmlns:di=\"http://www.omg.org/spec/DD/20100524/DI\" xsi:schemaLocation=\"http://www.omg.org/spec/BPMN/20100524/MODEL BPMN20.xsd\" id=\"sample-diagram\" targetNamespace=\"http://bpmn.io/schema/bpmn\">\n  <bpmn2:process id=\"Process_1\" isExecutable=\"false\">\n    <bpmn2:startEvent id=\"StartEvent_1\"/>\n  </bpmn2:process>\n  <bpmndi:BPMNDiagram id=\"BPMNDiagram_1\">\n    <bpmndi:BPMNPlane id=\"BPMNPlane_1\" bpmnElement=\"Process_1\">\n      <bpmndi:BPMNShape id=\"_BPMNShape_StartEvent_2\" bpmnElement=\"StartEvent_1\">\n        <dc:Bounds height=\"36.0\" width=\"36.0\" x=\"412.0\" y=\"240.0\"/>\n      </bpmndi:BPMNShape>\n    </bpmndi:BPMNPlane>\n  </bpmndi:BPMNDiagram>\n</bpmn2:definitions>";

@connect(state => ({flow: state.flow}))
@Form.create()
export default class FlowDesigner extends PureComponent {
  componentDidMount() {
    var bpmnModeler = new BpmnModeler({
      container: '#js-canvas',
      propertiesPanel: {
        parent: '#js-properties-panel'
      },
      additionalModules: [
        propertiesPanelModule, propertiesProviderModule
      ],
      moddleExtensions: {
        camunda: camundaModdleDescriptor
      }
    });

    bpmnModeler.importXML(newDiagramXML, function(err) {
      if (err) {
        console.error(err);
      }
    });
  }

  render() {
    return (<Card title="新建流程">
      <Layout style={{
          position: 'fixed',
          top: 56,
          left: 0,
          bottom: 0,
          right: 0
      }}>
        <Layout>
          <Layout.Content style={{
              overflow: 'auto',
              position: 'absolute',
              top: 0,
              left: 0,
              bottom: 0,
              right: 260
          }}>
            <div id="js-canvas" style={{
                height: '100%'
            }}></div>
          </Layout.Content>
        </Layout>
        <Layout.Sider breakpoint="md" width={260} style={{
            background: '#fff',
            overflow: 'auto',
            position: 'absolute',
            top: 0,
            bottom: 0,
            right: 0
          }}>
          <div id="js-properties-panel"></div>
        </Layout.Sider>
      </Layout>
    </Card>);
  }
}
