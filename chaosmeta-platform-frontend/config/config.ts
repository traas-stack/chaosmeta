import { defineConfig } from 'umi';
import routes from './router';

export default defineConfig({
  title: 'chaosmeta',
  routes,
  hash: true,
  historyWithQuery: {},
  history: {
    type: 'hash',
  },
  npmClient: 'yarn',
  proxy: {
    '/api': {
      target: 'http://jsonplaceholder.typicode.com/',
      changeOrigin: true,
      pathRewrite: { '^/api': '' },
    },
  },
});
