// 路由菜单相关的国际化
const routeMenu = {
  'menu.login': '登录',
  'menu.space': '空间',
  'menu.experimentCreate': '创建实验',
  'menu.space.overview': '空间概览',
  'menu.space.experiment': '实验',
  'menu.space.experimentDetail': '实验详情',
  'menu.space.experimentCreate': '创建实验',
  'menu.space.experimentResult': '实验结果',
  'menu.space.experimentResultDetail': '实验结果详情',
  'menu.space.settings': '空间设置',
  'menu.globalSettings': '全局设置',
  'menu.globalSettings.account': '账号管理',
  'menu.globalSettings.space': '空间管理',
};

// 空间概览
const spaceOverview = {
  'overview.workbench': '工作台',
  'overview.tip': '开始您的实验，只需要3步！',
  'overview.panel.close': '收起',
  'overview.panel.expand': '展开',
  'overview.step1.title': '创建实验',
  'overview.step1.description':
    '可选择实验模版快速构建实验场景，进行基础资源，如cpu燃烧等实验来验证应用系统的可靠性',
  'overview.step2.title': '执行实验',
  'overview.step2.description': '针对配置好的实验可发起攻击',
  'overview.step3.title': '查看实验结果',
  'overview.step3.description':
    '实验过程中可观测系统指标，实验完成后可查看实验结果，系统会自动度量',
  'overview.spaceOverview': '空间总览',
  'overview.statistics.newExperiment': '新增实验',
  'overview.statistics.performingExperiments': '执行实验',
  'overview.statistics.executionFailed': '执行失败',
  'overview.statistics.count': '个',
  'overview.statistics.times': '次',
  'overview.statistics.option.7': '最近7天',
  'overview.statistics.option.30': '最近30天',
  'overview.spaceOverview.tab.more': '查看全部实验',
  'overview.spaceOverview.tab1.title': '最近编辑的实验',
  'overview.spaceOverview.tab1.noAuth.empty.description':
    '当前页面暂无最近编辑的实验',
  'overview.spaceOverview.tab1.noAuth.empty.title':
    '您可以前往实验列表查看实验',
  'overview.spaceOverview.tab1.noAuth.empty.btn': '前往空间列表',
  'overview.spaceOverview.tab2.title': '即将运行的实验',
  'overview.spaceOverview.tab2.noAuth.empty.description':
    '当前页面暂无即将运行的实验',
  'overview.spaceOverview.tab3.title': '最近运行的实验结果',
  'overview.spaceOverview.tab3.noAuth.empty.description':
    '当前暂无最近运行的实验结果',
  'overview.spaceOverview.tab3.noAuth.empty.title':
    '您可以前往实验结果列表查看实验',
  'overview.spaceOverview.tab3.noAuth.empty.btn': '前往实验结果列表',
};

// 公用的
const publicText = {
  query: '查询',
  reset: '重置',
  pleaseInput: '请输入',
  pleaseSelect: '请选择',
  createExperiment: '创建实验',
  goToMemberManagement: '前往成员管理',
  name: '名称',
  creator: '创建人',
  numberOfExperiments: '实验次数',
  latestExperimentalTime: '最近实验时间',
  status: '状态',
  triggerMode: '触发方式',
  upcomingRunningTime: '即将运行时间',
  lastEditTime: '最近编辑时间',
  lastStartTime: '最近一次实验开始的时间',
  startTime: '发起时间',
  endTime: '结束时间',
  operate: '操作',
  edit: '编辑',
  copy: '复制',
  delete: '删除',
  stop: '停止',
  experiment: '实验',
  experimentList: '实验列表',
  experimentStartTime: '实验开始时间',
  experimentEndTime: '实验结束时间',
  experimentStatus: '实验状态',
  experimentResult: '实验结果',
  experimentResultList: '实验结果列表',
  timeType: '时间类型',
  'timeType.all': '全部',
  'timeType.7': '近7天',
  'timeType.30': '近30天',
  copyText: '复制成功！',
  deleteText: '删除成功！',
  updateText: '更新成功！',
  createText: '创建成功！',
  stopConfirmText: '你确定要停止实验吗？',
  deleteConfirmText: '你确定要删除吗？',
  run: '运行',
  basicInfo: '基本信息',
  label: '标签',
  lastOperationTime: '最近操作时间',
  description: '描述',
  experimentConfig: '实验配置',
  configInfo: '配置信息',
  attackRange: '攻击范围',
  nodeName: '节点名称',
  nodeType: '节点类型',
  waitTime: '等待时长',
  duration: '持续时长',
  totalDuration: '总时长',
  experimentProgress: '实验进度',
  experimentLog: '实验日志',
  experimentName: '实验名称',
  finish: '完成',
  check: '查看',
  inputPlaceholder: '请输入',
  selectPlaceholder: '请选择',
  experimentDescription: '实验描述',
  expression: '表达式',
  confirm: '确认',
  cancel: '取消',
  fault: '故障节点',
  measure: '度量引擎',
  flow: '流量注入',
  wait: '等待时长',
  undone: '未完成',
  limit: '大于0',
  second: '秒',
  minute: '分',
  yes: '是',
  no: '否',
  ruleText: '的取值为',
  moreSpace: '更多空间',
  createSpace: '新建空间',
  keyword: '请输入关键词',
  spaceName: '空间名称',
  spaceDescription: '空间描述',
  spaceDescriptionTip: '请尽量保持空间名称的简洁，不超过64个字符',
  spaceSetting: '空间设置',
  save: '保存',
  memberNumber: '成员数量',
  createTime: '创建时间',
  memberList: '成员列表',
  usernamePlaceholder: '请输入用户名',
  batchOperate: '批量操作',
  addMember: '添加成员',
  selected: '已选择',
  item: '项',
  clear: '清空',
  batchDelete: '批量删除',
  joinTime: '加入时间',
  permission: '权限',
  username: '用户名',
  readonly: '只读',
  write: '读写',
  search: '搜索',
  all: '全部',
  selectAll: '全选',
  submit: '提交',
};

// 实验
const experiment = {
  'experiment.table.noAuth.description':
    '您在该空间是只读权限，暂不支持创建实验。若想创建实验请去成员管理中找空间内有读写权限的成员修改权限',
  'experiment.table.title': '当前空间还没有实验数据',
  'experiment.table.description':
    '请先创建实验，您可选择自己创建实验也可以通过推荐实验来快速构建实验场景，来验证应用系统的可靠性',
  'experiment.delete.title': '确认要删除这个实验吗？',
  'experiment.delete.content':
    '删除实验将会删除该实验的配置，但不会删除历史实验结果！',
};

// 实验结果
const experimentResult = {
  'experimentResult.table.description': '当前暂无实验结果数据',
  'experimentResult.table.title': '您可以前往实验详情页面运行实验',
  'experimentResult.table.noAuth.title': '您可以前往实验列表页面查看实验',
  'experimentResult.table.btn': '前往实验列表',
  'experimentResult.stop.text': '实验已停止',
};

// 新增/编辑实验页面
const addExperiment = {
  'addExperiment.basic.tip': '请完善基本信息',
  'addExperiment.node.tip': '请完善节点信息',
};

// tag抽屉
const tag = {
  'tag.repeat.text': '标签已经存在，请重新输入',
  'tag.repeat.check': '标签已存在',
  'tag.empty.tip': '请添加标签',
  'tag.create.success.tip': '您已成功新建标签',
};

// 节点库
const nodeLibrary = {
  'nodeLibrary.title': '节点库',
};

// 空间下列表
const spaceDropdown = {
  'spaceDropdown.tip': '没有相关空间？查看',
};

// 创建空间
const createSpace = {
  'createSpace.confirm': '创建完成并去配置',
};

// 空间下空间设置
const spaceSetting = {
  'spaceSetting.tab.member': '成员管理',
  'spaceSetting.tab.tag': '标签管理',
};

// 空间下成员管理
const memberManageMent = {
  'memberManageMent.readonly.tip': '在当前空间有编辑权限',
  'memberManageMent.write.tip': '在当前空间只能查看',
  'memberManageMent.delete.title': '您确认要删除当前所选成员吗？',
  'memberManageMent.delete.content': '删除空间内成员，该成员将无法进入该空间！',
  'memberManageMent.delete.success.tip': '您已成功删除所选成员！',
  'memberManageMent.permission.success.tip': '权限修改成功',
};

// 空间下添加成员
const addMember = {
  'addMember.user.permission': '用户权限',
  'addMember.selected.user': '已选择用户',
  'addMember.noMore': '没有更多了～',
  'addMember.loading': '加载中...',
  'addMember.search.result': '暂未搜索到结果，请尝试切换关键词重新搜索',
  'addMember.add.success': '添加成功',
};

// 空间下标签管理
const tagManageMent = {
  'tagManageMent.title': '标签列表',
  'tagManageMent.add.title': '新建标签',
  'tagManageMent.search.placeholder': '请输入标签名称',
  'tagManageMent.table.empty.title': '您还没有标签数据',
  'tagManageMent.table.empty.description':
    '请新建标签，提前创建好标签在创建实验的时候可以直接选用快速为实验打上标签。',
  'tagManageMent.table.empty.noAuth.title': '当前还没有标签数据',
  'tagManageMent.table.empty.noAuth.description':
    '您在该空间只是只读权限，暂不支持添加标签。若想添加标签请去成员管理中找空间内有读写权限的成员修改权限',
  'tagManageMent.delete.success.tip': '你已成功删除标签',
  'tagManageMent.column.tagColor': '标签颜色',
  'tagManageMent.column.tagStyle': '标签样式',
};

// 登录/密码相关
const login = {
  login: '登录',
  register: '注册',
  signOut: '退出登录',
  updatePassword: '修改密码',
  password: '密码',
  'password.old': '原密码',
  'password.old.placeholder': '请输入原密码',
  'password.new': '新密码',
  'password.new.placeholder': '请输入新密码',
  'password.new.placeholder.again': '请再次输入新密码',
  'password.rule': '密码8-16位中英文大小写及下划线等特殊字符',
  'password.confirm': '请确认密码',
  'password.confirm.again': '确认密码',
  'password.error': '密码不正确',
  'password.new.confirm': '确认新密码',
  'password.success': '密码修改成功，即将跳转到登录页面重新登录',
  'password.placeholder': '请输入密码',
  'password.inconsistent': '两次密码不一致',
  welcome: '欢迎使用',
  'username.rule': '用户名可以使用中英文，长度不超过64个字符',
  notAccount: '还没有账号？',
  haveAccount: '已经有账号？',
  'reister.success': '注册成功，请登录',
};

// 账号相关
const account = {
  'account.title': '账号管理',
  'account.list': '账号列表',
  'account.search.placeholder': '请输入用户名进行搜索',
  role: '角色',
  admin: '管理员',
  generalUser: '普通用户',
  adminDescription: '拥有所有权限',
  generalUserDescription: '可登录查看，同时叠加空间内权限',
  'account.delete.title': '确认要删除当前所选账号吗？',
  'account.delete.content':
    '删除账号用户将无法登录平台，要再次使用只能重新注册！',
  'account.delete.success': '您已成功删除所选成员',
  'account.role.update': '用户角色修改成功',
};

// 空间管理
const spaceManagement = {
  'spaceManagement.title': '空间管理',
  'spaceManagement.spaceName.placeholder': '请输入空间名称',
  'spaceManagement.spaceMember.placeholder': '请输入空间成员',
  'spaceManagement.alert': '可联系空间内具有读写权限的成员添加为空间成员',
  'spaceManagement.member': '空间成员',
  'spaceManagement.tab.all': '全部空间',
  'spaceManagement.tab.related': '我相关的',
  'spaceManagement.delete.title': '确认要删除当前所选空间吗？',
  'spaceManagement.delete.success': '您已成功删除所选空间',
  'spaceManagement.noAuth.tip': '你没有该空间的权限，请联系读写成员',
  'spaceManagement.noAuth.readonly.tip': '只读用户暂无法使用此功能',
  'spaceManagement.write': '读写成员',
  'spaceManagement.experimentCount': '实验数量',

  
}

export default {
  ...routeMenu,
  ...spaceOverview,
  ...experiment,
  ...experimentResult,
  ...publicText,
  ...addExperiment,
  ...tag,
  ...nodeLibrary,
  ...spaceDropdown,
  ...createSpace,
  ...spaceSetting,
  ...memberManageMent,
  ...addMember,
  ...tagManageMent,
  ...login,
  ...account,
  ...spaceManagement
};
