import React, { PureComponent } from "react";
import { connect } from "dva";
import { Card, Form, Layout, Button, notification, Modal } from "antd";

var BpmnModeler = require("bpmn-js/lib/Modeler");
var propertiesPanelModule = require("bpmn-js-properties-panel");
var propertiesProviderModule = require("bpmn-js-properties-panel/lib/provider/camunda");
var camundaModdleDescriptor = require("camunda-bpmn-moddle/resources/camunda");
var fileDownload = require("js-file-download");

const { confirm } = Modal;
@connect(state => ({ flow: state.flow }))
@Form.create()
export default class FlowCard extends PureComponent {
  componentDidMount() {
    var bpmnModeler = new BpmnModeler({
      container: "#js-canvas",
      propertiesPanel: {
        parent: "#js-properties-panel"
      },
      additionalModules: [propertiesPanelModule, propertiesProviderModule],
      moddleExtensions: {
        camunda: camundaModdleDescriptor
      }
    });

    bpmnModeler.createDiagram(function(err) {
      if (err) {
        notification.error({ message: "设计器加载失败" });
        return console.error(err);
      }
    });

    this.props.dispatch({
      type: "flow/loadForm",
      payload: this.props.match.params,
      bpmnModeler: bpmnModeler
    });
  }

  onSaveOKClick = () => {
    const that = this;
    this.props.flow.bpmnModeler.saveXML({ format: true }, (err, xml) => {
      if (err) {
        notification.error({ message: "保存XML文件失败" });
        return console.error(err);
      }

      that.props.dispatch({
        type: "flow/submit",
        payload: { xml }
      });
    });
  };

  onSaveClick = () => {
    confirm({
      title: "确认保存该流程吗？",
      content: "流程保存之后不允许修改！",
      okText: "确认",
      okType: "danger",
      cancelText: "取消",
      onOk: this.onSaveOKClick.bind(this)
    });
  };

  onExportXMLClick = () => {
    this.props.flow.bpmnModeler.saveXML({ format: true }, (err, xml) => {
      if (err) {
        notification.error({ message: "保存XML文件失败" });
        return console.error(err);
      }
      fileDownload(xml, "diagram.xml");
    });
  };

  onExportSVGClick = () => {
    this.props.flow.bpmnModeler.saveSVG((err, svg) => {
      if (err) {
        notification.error({ message: "保存SVG文件失败" });
        return console.error(err);
      }
      fileDownload(svg, "diagram.svg");
    });
  };

  renderSubmit = () => {
    const { flow: { submitting, submitVisible } } = this.props;
    if (submitVisible) {
      return (
        <Button
          icon="save"
          type="primary"
          loading={submitting}
          onClick={this.onSaveClick}
        >
          保存
        </Button>
      );
    }
  };

  render() {
    const { formTitle, formData } = this.props.flow;

    return (
      <Card
        title={formData.name ? formData.name + " - " + formTitle : formTitle}
        extra={<a onClick={this.props.history.goBack}>返回</a>}
      >
        <Layout
          style={{
            position: "fixed",
            top: 56,
            left: 0,
            bottom: 0,
            right: 0
          }}
        >
          <Layout.Content
            style={{
              backgroundSize: "50px 50px",
              backgroundImage:
                "linear-gradient(to right, gainsboro 1px, transparent 1px), linear-gradient(to bottom, gainsboro 1px, transparent 1px)",
              overflow: "auto",
              position: "absolute",
              top: 0,
              left: 0,
              bottom: 50,
              right: 260
            }}
          >
            <div
              id="js-canvas"
              style={{
                height: "100%"
              }}
            />
          </Layout.Content>
          <Layout.Sider
            breakpoint="md"
            width={260}
            style={{
              background: "#fff",
              overflow: "auto",
              position: "absolute",
              top: 0,
              bottom: 50,
              right: 0
            }}
          >
            <div id="js-properties-panel" />
          </Layout.Sider>
          <Layout.Footer
            style={{
              background: "#fff",
              position: "absolute",
              height: 50,
              left: 0,
              bottom: 0,
              right: 0,
              padding: 0,
              textAlign: "center"
            }}
          >
            <Form layout="inline">
              <Form.Item>
                {this.renderSubmit()}
                <Button
                  icon="download"
                  type="dashed"
                  onClick={this.onExportXMLClick}
                  style={{
                    marginLeft: 8
                  }}
                >
                  导出XML
                </Button>
                <Button
                  icon="download"
                  type="dashed"
                  onClick={this.onExportSVGClick}
                  style={{
                    marginLeft: 8
                  }}
                >
                  导出SVG
                </Button>
              </Form.Item>
            </Form>
          </Layout.Footer>
        </Layout>
      </Card>
    );
  }
}
