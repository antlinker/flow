export default {

  namespace: 'flow',

  state: {
    listLoading: false,
    data: {
      list: [],
      pagination: {},
    },
  },

  effects: {
    * fetch({ payload }, { call, put }) {  // eslint-disable-line

    },
  },

  reducers: {
    saveListLoading(state, action) {
      return { ...state, listLoading: action.payload };
    },
    saveData(state, action) {
      return { ...state, data: action.payload };
    }
  },

};
