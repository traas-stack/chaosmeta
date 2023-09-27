import { defineConfig } from '@umijs/max';
import routes from './router';

export default defineConfig({
  title: 'chaosmeta',
  antd: {},
  access: {},
  model: {},
  initialState: {},
  request: {},
  layout: {
    title: 'chaosmeta',
  },
  hash: true,
  historyWithQuery: {},
  routes,
  npmClient: 'yarn',
  styledComponents: {},
  proxy: {
    '/users': {
      target: 'http://127.0.0.1:8082/',
      changeOrigin: true,
      pathRewrite: { '^/api': '' },
    },
    '/chaosmeta': {
      target: 'http://127.0.0.1:8082/',
      changeOrigin: true,
      pathRewrite: { '^/api': '' },
    },
  },
});
