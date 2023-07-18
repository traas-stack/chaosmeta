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
    '/api': {
      target: 'http://jsonplaceholder.typicode.com/',
      changeOrigin: true,
      pathRewrite: { '^/api': '' },
    },
  },
});
