import React from 'react';
import { Router, Route, Switch } from 'dva/router';
import dynamic from 'dva/dynamic';
import { LocaleProvider } from 'antd';
import zhCN from 'antd/lib/locale-provider/zh_CN';

// 动态加载组件
function dynamicWrapper(app, models, component) {
  return dynamic({
    app,
    models: () => models.map(m => import(`./models/${m}.js`)),
    component: () => import(`./routes/${component}`),
  });
}

function RouterConfig({ history, app }) {
  return (
    <LocaleProvider locale={zhCN}>
      <Router history={history}>
        <Switch>
          <Route path="/" exact component={dynamicWrapper(app, ['flow'], 'FlowList')} />
          <Route path="/flow/:action/:id" component={dynamicWrapper(app, ['flow'], 'FlowCard')} />
          <Route path="/flow/:action" component={dynamicWrapper(app, ['flow'], 'FlowCard')} />
        </Switch>
      </Router>
    </LocaleProvider>
  );
}

export default RouterConfig;
