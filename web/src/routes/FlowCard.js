import React, { PureComponent } from 'react';
import { connect } from 'dva';
import { Card, Form, Layout, Button, notification, Modal } from 'antd';

import AceEditor from 'react-ace';

import 'brace/mode/xml';
import 'brace/theme/github';

const BpmnModeler = require('bpmn-js/lib/Modeler');
const propertiesPanelModule = require('bpmn-js-properties-panel');
const propertiesProviderModule = require('bpmn-js-properties-panel/lib/provider/camunda');
const camundaModdleDescriptor = require('camunda-bpmn-moddle/resources/camunda');
const fileDownload = require('js-file-download');

const { confirm } = Modal;
@connect(state => ({ flow: state.flow }))
@Form.create()
export default class FlowCard extends PureComponent {
  state = {
    editorVisible: false,
    bpmnXML: '',
  };

  componentDidMount() {
    const bpmnModeler = new BpmnModeler({
      container: '#js-canvas',
      propertiesPanel: {
        parent: '#js-properties-panel',
      },
      additionalModules: [propertiesPanelModule, propertiesProviderModule],
      moddleExtensions: {
        camunda: camundaModdleDescriptor,
      },
    });

    bpmnModeler.createDiagram((err) => {
      if (err) {
        notification.error({ message: `设计器加载失败:${err}` });
      }
    });

    this.props.dispatch({
      type: 'flow/loadForm',
      payload: this.props.match.params,
      bpmnModeler,
    });
  }

  onEditorChange = (val) => {
    this.setState({ bpmnXML: val });
  };

  onModalCancelClick = () => {
    this.setState({ editorVisible: false });
  };

  onModalOKClick = () => {
    this.setState({ editorVisible: false });
    this.props.flow.bpmnModeler.importXML(this.state.bpmnXML, (err) => {
      if (err) {
        notification.error({ message: `设计器加载失败:${err}` });
      }
    });
  };

  onSaveOKClick = () => {
    const that = this;
    this.props.flow.bpmnModeler.saveXML({ format: true }, (err, xml) => {
      if (err) {
        notification.error({ message: `保存XML文件失败:${err}` });
        return;
      }

      that.props.dispatch({
        type: 'flow/submit',
        payload: { xml },
      });
    });
  };

  onSaveClick = () => {
    confirm({
      title: '确认保存该流程吗？',
      content: '流程保存之后不允许修改！',
      okText: '确认',
      okType: 'danger',
      cancelText: '取消',
      onOk: this.onSaveOKClick.bind(this),
    });
  };

  onEditXMLClick = () => {
    const that = this;
    this.setState({ editorVisible: true });
    this.props.flow.bpmnModeler.saveXML({ format: true }, (err, xml) => {
      if (err) {
        notification.error({ message: `导出XML失败：${err}` });
        return;
      }
      that.setState({ bpmnXML: xml });
    });
  };

  onExportXMLClick = () => {
    this.props.flow.bpmnModeler.saveXML({ format: true }, (err, xml) => {
      if (err) {
        notification.error({ message: `导出XML失败：${err}` });
        return;
      }
      fileDownload(xml, 'diagram.xml');
    });
  };

  onExportSVGClick = () => {
    this.props.flow.bpmnModeler.saveSVG((err, svg) => {
      if (err) {
        notification.error({ message: `导出SVG失败：${err}` });
        return;
      }
      fileDownload(svg, 'diagram.svg');
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
          style={{
            marginRight: 8,
          }}
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
        title={formData.name ? `${formData.name} - ${formTitle}` : formTitle}
        extra={<a onClick={this.props.history.goBack}>返回</a>}
      >
        <Layout
          style={{
            position: 'fixed',
            top: 56,
            left: 0,
            bottom: 0,
            right: 0,
          }}
        >
          <Layout.Content
            style={{
              backgroundSize: '50px 50px',
              backgroundImage:
                'linear-gradient(to right, gainsboro 1px, transparent 1px), linear-gradient(to bottom, gainsboro 1px, transparent 1px)',
              overflow: 'auto',
              position: 'absolute',
              top: 0,
              left: 0,
              bottom: 50,
              right: 260,
            }}
          >
            <div
              id="js-canvas"
              style={{
                height: '100%',
              }}
            />
          </Layout.Content>
          <Layout.Sider
            breakpoint="md"
            width={260}
            style={{
              background: '#fff',
              overflow: 'auto',
              position: 'absolute',
              top: 0,
              bottom: 50,
              right: 0,
            }}
          >
            <div id="js-properties-panel" />
          </Layout.Sider>
          <Layout.Footer
            style={{
              background: '#fff',
              position: 'absolute',
              height: 50,
              left: 0,
              bottom: 0,
              right: 0,
              padding: 0,
              textAlign: 'center',
            }}
          >
            <Form layout="inline">
              <Form.Item>
                {this.renderSubmit()}
                <Button icon="edit" type="dashed" onClick={this.onEditXMLClick}>
                  编辑XML
                </Button>
                <Button
                  icon="download"
                  type="dashed"
                  onClick={this.onExportXMLClick}
                  style={{
                    marginLeft: 8,
                  }}
                >
                  导出XML
                </Button>
                <Button
                  icon="download"
                  type="dashed"
                  onClick={this.onExportSVGClick}
                  style={{
                    marginLeft: 8,
                  }}
                >
                  导出SVG
                </Button>
              </Form.Item>
            </Form>
          </Layout.Footer>
        </Layout>

        <Modal
          title="编辑XML"
          width={1000}
          visible={this.state.editorVisible}
          onOk={this.onModalOKClick}
          onCancel={this.onModalCancelClick}
        >
          <AceEditor
            mode="xml"
            theme="github"
            name="xmlEditor"
            width="950px"
            height="450px"
            onChange={this.onEditorChange}
            fontSize={14}
            showPrintMargin
            showGutter
            highlightActiveLine
            value={this.state.bpmnXML}
            setOptions={{
              enableBasicAutocompletion: false,
              enableLiveAutocompletion: false,
              enableSnippets: false,
              showLineNumbers: true,
              tabSize: 2,
            }}
          />
        </Modal>
      </Card>
    );
  }
}
