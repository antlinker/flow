import { routerRedux } from 'dva/router';
import { notification, message } from 'antd';
import { queryPage, get, del, save } from '../services/flow';

export default {
  namespace: 'flow',
  state: {
    loading: false,
    data: {
      list: [],
      pagination: {},
    },
    search: {},
    bpmnModeler: undefined,
    formAction: '',
    formTitle: '',
    submitVisible: true,
    formData: {},
    submitting: false,
  },
  effects: {
    *fetch({ payload, pagination }, { call, put, select }) {
      yield put({
        type: 'changeLoading',
        payload: true,
      });

      let search = yield select(state => state.flow.search);

      if (payload) {
        search = { ...search, ...payload };
        yield put({
          type: 'saveSearch',
          payload,
        });
      }

      if (pagination) {
        search = { ...search, ...pagination };
      }

      const response = yield call(queryPage, search);
      yield put({
        type: 'saveData',
        payload: response,
      });

      yield put({
        type: 'changeLoading',
        payload: false,
      });
    },
    *loadForm({ payload, bpmnModeler }, { call, put }) {
      yield [
        put({
          type: 'saveBpmnModeler',
          payload: bpmnModeler,
        }),
        put({
          type: 'saveFormAction',
          payload: payload.action,
        }),
        put({
          type: 'changeSubmitVisible',
          payload: true,
        }),
      ];

      if (payload.action === 'add') {
        yield put({
          type: 'saveFormTitle',
          payload: '新建流程',
        });
      } else {
        if (payload.action === 'copy') {
          yield put({
            type: 'saveFormTitle',
            payload: '复制流程',
          });
        } else if (payload.action === 'view') {
          yield put({
            type: 'saveFormTitle',
            payload: '查看流程',
          });
          yield put({
            type: 'changeSubmitVisible',
            payload: false,
          });
        }

        const response = yield call(get, { record_id: payload.id });
        yield put({
          type: 'saveFormData',
          payload: response,
        });
        bpmnModeler.importXML(response.xml, (err) => {
          if (err) {
            notification.error({ message: `设计器加载失败:${err}` });
          }
        });
      }
    },
    *delete({ payload }, { call, put }) {
      const response = yield call(del, payload);
      if (response === 'ok') {
        message.success('删除成功');
        yield put({
          type: 'fetch',
        });
      }
    },
    *submit({ payload }, { call, put }) {
      yield put({
        type: 'changeSubmitting',
        payload: true,
      });

      const response = yield call(save, payload);
      if (response === 'ok') {
        yield put({
          type: 'changeSubmitting',
          payload: false,
        });
        message.success('保存成功');
        yield put(routerRedux.push('/'));
      }
    },
  },
  reducers: {
    changeLoading(state, action) {
      return { ...state, loading: action.payload };
    },
    saveSearch(state, action) {
      return { ...state, search: action.payload };
    },
    saveData(state, action) {
      return { ...state, data: action.payload };
    },
    saveBpmnModeler(state, action) {
      return { ...state, bpmnModeler: action.payload };
    },
    saveFormAction(state, action) {
      return { ...state, formAction: action.payload };
    },
    saveFormTitle(state, action) {
      return { ...state, formTitle: action.payload };
    },
    changeSubmitVisible(state, action) {
      return { ...state, submitVisible: action.payload };
    },
    saveFormData(state, action) {
      return { ...state, formData: action.payload };
    },
    changeSubmitting(state, action) {
      return { ...state, submitting: action.payload };
    },
  },
};
