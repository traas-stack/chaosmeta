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
      target: 'http://30.46.241.207:8082/',
      changeOrigin: true,
      pathRewrite: { '^/api': '' },
    },
    '/chaosmeta': {
      target: 'http://30.46.241.207:8082/',
      changeOrigin: true,
      pathRewrite: { '^/api': '' },
    },
    '/chaos': {
      target: 'http://antchaos-gz00b-006002015020.sa128-sqa.alipay.net/',
      changeOrigin: true,
      pathRewrite: { '^/api': '' },
    },
  },
});
