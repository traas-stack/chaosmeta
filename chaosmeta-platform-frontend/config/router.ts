export default [
  {
    path: '/',
    redirect: '/login',
  },
  { path: '/*', component: '@/pages/404' },
  {
    name: 'login',
    path: '/login',
    component: './Login',
    layout: false,
  },
  {
    name: 'experimentCreate',
    path: '/space/experiment/add',
    component: './Space/AddExperiment',
    layout: false,
  },
  {
    name: 'space',
    path: '/space',
    key: '/space',
    hideInBreadcrumb: true,
    routes: [
      {
        path: '/space',
        redirect: './overview',
      },
      {
        name: 'overview',
        path: '/space/overview',
        component: './Space/SpaceOverview',
        icon: 'https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*ySleT74WnD4AAAAAAAAAAAAADmKmAQ/original',
      },
      {
        name: 'experiment',
        path: '/space/experiment',
        component: './Space/Experiment',
        icon: 'https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*LTdSQbHlFP8AAAAAAAAAAAAADmKmAQ/original',
      },
      {
        name: 'experimentCreate',
        path: '/space/experiment/choose',
        component: './Space/ChooseExperiment',
        hideInMenu: true,
      },
      {
        name: 'experimentDetail',
        path: '/space/experiment/detail',
        component: './Space/ExperimentDetail',
        hideInMenu: true,
      },
      {
        name: 'experimentResultDetail',
        path: '/space/experiment-result/detail',
        component: './Space/ExperimentResultDetail',
        hideInMenu: true,
      },
      {
        name: 'experimentResult',
        path: '/space/experiment-result',
        component: './Space/ExperimentResult',
        parentKeys: [''],
        icon: 'https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*bABiRoluWWUAAAAAAAAAAAAADmKmAQ/original',
      },
      {
        name: 'settings',
        path: `/space/setting`,
        component: './Space/SpaceSetting',
        icon: 'https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*KesMQo37t4sAAAAAAAAAAAAADmKmAQ/original',
      },
    ],
  },
  {
    name: 'globalSettings',
    path: '/setting',
    hideInBreadcrumb: true,
    routes: [
      {
        path: '/setting',
        redirect: './account',
      },
      {
        name: 'account',
        path: '/setting/account',
        component: './GlobalSetting/Account',
        icon: 'https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*Yvf-TKO1tPAAAAAAAAAAAAAADmKmAQ/original',
      },
      {
        name: 'space',
        path: '/setting/space',
        component: './GlobalSetting/Space',
        icon: 'https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*8FUVTpn7RXIAAAAAAAAAAAAADmKmAQ/original',
      },
      // {
      //   name: 'Agent管理',
      //   path: '/setting/agent',
      //   component: './GlobalSetting/Agent',
      //   icon: 'https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*tFm6TIYpRC8AAAAAAAAAAAAADmKmAQ/original',
      // },
    ],
  },
];
