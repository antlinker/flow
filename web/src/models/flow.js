import { queryPage } from "../services/flow";

export default {
  namespace: "flow",
  state: {
    loading: false,
    data: {
      list: [],
      pagination: {}
    },
    search: {}
  },
  effects: {
    *fetch({ payload, pagination }, { call, put, select }) {
      yield put({
        type: "changeLoading",
        payload: true
      });

      let search = yield select(state => state.flow.search);

      if (payload) {
        search = { ...search, ...payload };
        yield put({
          type: "saveSearch",
          payload: payload
        });
      }

      if (pagination) {
        search = { ...search, ...pagination };
      }

      const response = yield call(queryPage, search);
      yield put({
        type: "saveData",
        payload: response
      });

      yield put({
        type: "changeLoading",
        payload: false
      });
    }
  },
  reducers: {
    changeLoading(state, action) {
      return {
        ...state,
        loading: action.payload
      };
    },
    saveSearch(state, action) {
      return {
        ...state,
        search: action.payload
      };
    },
    saveData(state, action) {
      return {
        ...state,
        data: action.payload
      };
    }
  }
};
