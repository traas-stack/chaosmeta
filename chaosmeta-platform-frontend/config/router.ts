export default [
  {
    path: '/',
    redirect: '/space/overview',
  },
  { path: '/*', component: '@/pages/404' },
  {
    name: '登录',
    path: '/login',
    component: './Login',
    layout: false,
  },
  {
    name: '空间',
    path: '/space',
    hideInBreadcrumb: true,
    routes: [
      {
        path: '/space',
        redirect: './overview',
      },
      {
        name: '空间',
        path: '/space/overview',
        component: './Space/SpaceOverview',
        icon: 'https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*ySleT74WnD4AAAAAAAAAAAAADmKmAQ/original',
      },
      {
        name: '实验',
        path: '/space/experiment',
        component: './Space/Experiment',
        icon: 'https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*LTdSQbHlFP8AAAAAAAAAAAAADmKmAQ/original',
      },
      {
        name: '实验结果',
        path: '/space/result',
        component: './Space/ExperimentResult',
        icon: 'https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*bABiRoluWWUAAAAAAAAAAAAADmKmAQ/original',
      },
      {
        name: '空间设置',
        path: `/space/setting`,
        component: './Space/SpaceSetting',
        icon: 'https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*KesMQo37t4sAAAAAAAAAAAAADmKmAQ/original',
      },
    ],
  },
  {
    name: '全局设置',
    path: '/setting',
    hideInBreadcrumb: true,
    routes: [
      {
        path: '/setting',
        redirect: './account',
      },
      {
        name: '账号管理',
        path: '/setting/account',
        component: './GlobalSetting/Account',
        icon: 'https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*ySleT74WnD4AAAAAAAAAAAAADmKmAQ/original',
      },
      {
        name: '空间管理',
        path: '/setting/space',
        component: './GlobalSetting/Space',
        icon: 'https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*ySleT74WnD4AAAAAAAAAAAAADmKmAQ/original',
      },
      {
        name: 'Agent管理',
        path: '/setting/agent',
        component: './GlobalSetting/Agent',
        icon: 'https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*ySleT74WnD4AAAAAAAAAAAAADmKmAQ/original',
      },
    ],
  },
];
