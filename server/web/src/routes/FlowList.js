import React, {PureComponent} from 'react';
import {connect} from 'dva';
import {Link} from 'dva/router';
import {
  Row,
  Col,
  Card,
  Form,
  Input,
  Button,
  Table,
  Divider
} from 'antd';
import styles from './FlowList.less';

@connect(state => ({flow: state.flow}))
@Form.create()
export default class FlowList extends PureComponent {
  componentDidMount() {}

  renderSearchForm() {
    const {getFieldDecorator} = this.props.form;
    return (<Form onSubmit={this.handleSearch} layout="inline">
      <Row gutter={{
          md: 8,
          lg: 24,
          xl: 48
        }}>
        <Col md={8} sm={24}>
          <Form.Item label="流程名称">
            {getFieldDecorator('name')(<Input placeholder="请输入"/>)}
          </Form.Item>
        </Col>
        <Col md={8} sm={24}>
          <span>
            <Button type="primary" htmlType="submit">查询</Button>
            <Button style={{
                marginLeft: 8
              }} onClick={this.handleFormReset}>重置</Button>
          </span>
        </Col>
      </Row>
    </Form>);
  }

  render() {
    const {
      flow: {
        listLoading,
        data: {
          list,
          pagination
        }
      }
    } = this.props;

    const columns = [
      {
        title: '流程编号',
        dataIndex: 'code'
      }, {
        title: '流程名称',
        dataIndex: 'name'
      }, {
        title: '版本号',
        dataIndex: 'version'
      }, {
        title: '创建时间',
        dataIndex: 'created'
      }, {
        title: '操作',
        dataIndex: 'record_id',
        width: 120,
        render: val => (<div>
          <Link to={`/flow/${val}`}>
            查看
          </Link>
          <Divider type="vertical"/>
          <a onClick={() => {
              this.handleDelClick(val);
            }}>复制</a>
          <Divider type="vertical"/>
          <a onClick={() => {
              this.handleDelClick(val);
            }}>删除</a>
        </div>)
      }
    ];

    const paginationProps = {
      showSizeChanger: true,
      showQuickJumper: true,
      showTotal: (total) => {
        return <span>共{total}条</span>;
      },
      ...pagination
    };

    return (<Card title="流程管理" bordered={false}>
      <div className={styles.tableList}>
        <div className={styles.tableListForm}>
          {this.renderSearchForm()}
        </div>
        <div className={styles.tableListOperator}>
          <Link to="/flow_designer">
            <Button icon="plus" type="primary">
              新建
            </Button>
          </Link>
        </div>
        <Table loading={listLoading} rowKey={record => record.record_id} dataSource={list} columns={columns} pagination={paginationProps}/>
      </div>
    </Card>);
  }
}
