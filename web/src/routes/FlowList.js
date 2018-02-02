import React, { PureComponent } from 'react';
import { connect } from 'dva';
import { Link } from 'dva/router';
import { Row, Col, Card, Form, Input, Button, Table, Divider, Modal } from 'antd';
import { formatTimestamp } from '../utils/util';
import styles from './FlowList.less';

const { confirm } = Modal;
@connect(state => ({ flow: state.flow }))
@Form.create()
export default class FlowList extends PureComponent {
  componentDidMount() {
    this.props.dispatch({
      type: 'flow/fetch',
    });
  }

  onDelOKClick = (id) => {
    this.props.dispatch({
      type: 'flow/delete',
      payload: { record_id: id },
    });
  };

  onDelClick = (id) => {
    confirm({
      title: '确认删除该流程吗？',
      content: '删除该流程将会删除与流程相关的节点数据！',
      okText: '确认',
      okType: 'danger',
      cancelText: '取消',
      onOk: this.onDelOKClick.bind(this, id),
    });
  };

  onResetClick = () => {
    this.props.form.resetFields();
    this.props.dispatch({
      type: 'flow/fetch',
      payload: { code: undefined, name: undefined },
    });
  };

  onSearchClick = (e) => {
    if (e) {
      e.preventDefault();
    }
    this.props.form.validateFields((err, values) => {
      if (err) return;
      this.props.dispatch({
        type: 'flow/fetch',
        payload: values,
      });
    });
  };

  onTableChange = (pagination) => {
    this.props.dispatch({
      type: 'flow/fetch',
      pagination: {
        current: pagination.current,
        pageSize: pagination.pageSize,
      },
    });
  };

  renderSearchForm() {
    const { getFieldDecorator } = this.props.form;
    return (
      <Form onSubmit={this.onSearchClick} layout="inline">
        <Row
          gutter={{
            md: 8,
            lg: 24,
            xl: 48,
          }}
        >
          <Col md={8} sm={24}>
            <Form.Item label="流程编号">
              {getFieldDecorator('code')(<Input placeholder="请输入" />)}
            </Form.Item>
          </Col>
          <Col md={8} sm={24}>
            <Form.Item label="流程名称">
              {getFieldDecorator('name')(<Input placeholder="请输入" />)}
            </Form.Item>
          </Col>
          <Col md={8} sm={24}>
            <span>
              <Button type="primary" htmlType="submit">
                查询
              </Button>
              <Button
                style={{
                  marginLeft: 8,
                }}
                onClick={this.onResetClick}
              >
                重置
              </Button>
            </span>
          </Col>
        </Row>
      </Form>
    );
  }

  render() {
    const { flow: { loading, data: { list, pagination } } } = this.props;

    const columns = [
      {
        title: '流程编号',
        dataIndex: 'code',
      },
      {
        title: '流程名称',
        dataIndex: 'name',
      },
      {
        title: '版本号',
        dataIndex: 'version',
      },
      {
        title: '创建时间',
        dataIndex: 'created',
        render: val => <span>{formatTimestamp(val)}</span>,
      },
      {
        title: '操作',
        dataIndex: 'record_id',
        width: 240,
        render: val => (
          <div>
            <Link to={`/flow/view/${val}`}>查看</Link>
            <Divider type="vertical" />
            <Link to={`/flow/copy/${val}`}>复制</Link>
            <Divider type="vertical" />
            <a
              onClick={() => {
                this.onDelClick(val);
              }}
            >
              删除
            </a>
          </div>
        ),
      },
    ];

    const paginationProps = {
      showSizeChanger: true,
      showQuickJumper: true,
      showTotal: (total) => {
        return <span>共{total}条</span>;
      },
      ...pagination,
    };

    return (
      <Card title="流程管理" bordered={false}>
        <div className={styles.tableList}>
          <div className={styles.tableListForm}>{this.renderSearchForm()}</div>
          <div className={styles.tableListOperator}>
            <Link to="/flow/add">
              <Button icon="plus" type="primary">
                新建
              </Button>
            </Link>
          </div>
          <Table
            loading={loading}
            rowKey={record => record.record_id}
            dataSource={list}
            columns={columns}
            pagination={paginationProps}
            onChange={this.onTableChange}
          />
        </div>
      </Card>
    );
  }
}
