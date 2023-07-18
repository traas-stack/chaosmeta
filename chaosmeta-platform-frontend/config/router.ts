export default [
  {
    path: '/',
    redirect: '/space',
  },
  { path: '/*', component: '@/pages/404' },
  {
    name: '空间',
    path: '/space',
    component: './Space',
    icon: 'https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*ySleT74WnD4AAAAAAAAAAAAADmKmAQ/original',
  },
  {
    name: '空间设置',
    path: '/space-setting',
    component: './SpaceSetting',
    icon: 'https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*ySleT74WnD4AAAAAAAAAAAAADmKmAQ/original',
  },
  {
    name: '登录',
    path: '/login',
    component: './Login',
    layout: false,
  },
];
