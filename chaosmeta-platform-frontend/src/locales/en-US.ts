// 路由菜单相关的国际化
const routeMenu = {
  'menu.login': 'login',
  'menu.space': 'space',
  'menu.experimentCreate': 'create experiment',
  'menu.space.overview': 'space overview',
  'menu.space.experiment': 'experiment',
  'menu.space.experimentDetail': 'experiment detail',
  'menu.space.experimentCreate': 'create experiment',
  'menu.space.experimentResult': 'experiment results',
  'menu.space.experimentResultDetail': 'experiment result details',
  'menu.space.settings': 'space setting',
  'menu.globalSettings': 'global settings',
  'menu.globalSettings.account': 'account management',
  'menu.globalSettings.space': 'space management',
};

// 空间概览
const spaceOverview = {
  'overview.workbench': 'workbench',
  'overview.tip': 'Start your experiment in just 3 steps!',
  'overview.panel.close': 'close',
  'overview.panel.expand': 'expand',
  'overview.step1.title': 'create experiment',
  'overview.step1.description':
    'You can choose experimental templates to quickly build experimental scenarios and conduct basic resources such as CPU combustion experiments to verify the reliability of the application system',
  'overview.step2.title': 'perform experiments',
  'overview.step2.description':
    'Attacks can be launched against configured experiments',
  'overview.step3.title': 'view experimental results',
  'overview.step3.description':
    'System indicators can be observed during the experiment, and the experimental results can be viewed after the experiment is completed. The system will automatically measure',
  'overview.spaceOverview': 'space overview',
  'overview.statistics.newExperiment': 'new experiment',
  'overview.statistics.performingExperiments': 'performing experiments',
  'overview.statistics.executionFailed': 'execution failed',
  'overview.statistics.count': ' ',
  'overview.statistics.times': 'times',
  'overview.statistics.option.7': 'Last 7 days',
  'overview.statistics.option.30': 'Last 30 days',
  'overview.spaceOverview.tab.more': 'View all experiments',
  'overview.spaceOverview.tab1.title': 'Recently edited experiments',
  'overview.spaceOverview.tab1.noAuth.empty.description':
    'There are no recently edited experiments on the current page.',
  'overview.spaceOverview.tab1.noAuth.empty.title':
    'You can view experiments by going to the experiment list',
  'overview.spaceOverview.tab1.noAuth.empty.btn': 'Go to space list',
  'overview.spaceOverview.tab2.title': 'Upcoming experiments',
  'overview.spaceOverview.tab2.noAuth.empty.description':
    'There are no experiments about to be run on the current page.',
  'overview.spaceOverview.tab3.title': 'Results of recently run experiments',
  'overview.spaceOverview.tab3.noAuth.empty.description':
    'There are currently no recently run experimental results.',
  'overview.spaceOverview.tab3.noAuth.empty.title':
    'You can go to the experiment results list to view the experiment',
  'overview.spaceOverview.tab3.noAuth.empty.btn':
    'Go to experiment results list',
};

// 公用的
const publicText = {
  query: 'query',
  reset: 'reset',
  pleaseInput: 'please input',
  pleaseSelect: 'please select',
  createExperiment: 'create experiment',
  goToMemberManagement: 'Go to member management',
  name: 'name',
  creator: 'creator',
  numberOfExperiments: 'number of experiments',
  latestExperimentalTime: 'last experiment time',
  status: 'state',
  triggerMode: 'trigger mode',
  upcomingRunningTime: 'upcoming running time',
  lastEditTime: 'last edit time',
  lastStartTime: 'The time when the most recent experiment started',
  startTime: 'start time',
  endTime: 'end time',
  operate: 'operate',
  edit: 'edit',
  copy: 'copy',
  delete: 'delete',
  stop: 'stop',
  experiment: 'experiment',
  experimentList: 'Experiment list',
  experimentStartTime: 'experiment start time',
  experimentEndTime: 'experiment end time',
  experimentStatus: 'experiment status',
  experimentResult: 'experiment results',
  experimentResultList: 'experiment result list',
  timeType: 'time type',
  'timeType.all': 'all',
  'timeType.7': 'last 7 days',
  'timeType.30': 'last 30 days',
  copyText: 'copy successfully!',
  deleteText: 'successfully deleted!',
  updateText: 'update completed!',
  createText: 'created successfully!',
  stopConfirmText: 'Are you sure you want to stop the experiment?',
  deleteConfirmText: 'Are you sure you want to delete?',
  run: 'run',
  basicInfo: 'Basic Information',
  label: 'label',
  lastOperationTime: 'last operation time',
  description: 'description',
  experimentConfig: 'Experimental configuration',
  configInfo: 'Configuration information',
  attackRange: 'Attack range',
  nodeName: 'node name',
  nodeType: 'node type',
  waitTime: 'waiting time',
  duration: 'duration',
  totalDuration: 'total duration',
  experimentProgress: 'Experiment progress',
  experimentLog: 'Experiment log',
  experimentName: 'Experiment name',
  finish: 'finish',
  check: 'check',
  inputPlaceholder: 'please input',
  selectPlaceholder: 'please select',
  experimentDescription: 'experiment description',
  expression: 'expression',
  confirm: 'Confirm',
  cancel: 'Cancel',
  fault: 'faulty node',
  measure: 'measurement engine',
  flow: 'flow injection',
  wait: 'waiting time',
  undone: 'undone',
  limit: 'greater than 0',
  second: 'second',
  minute: 'minute',
  yes: 'yes',
  no: 'no',
  ruleText: 'the value is',
  moreSpace: 'more space',
  createSpace: 'create a new space',
  keyword: 'please enter a keyword',
  spaceName: 'space name',
  spaceDescription: 'space description',
  spaceDescriptionTip:
    'please try to keep the space name concise, no more than 64 characters',
  spaceSetting: 'space setting',
  save: 'save',
  memberNumber: 'number of members',
  createTime: 'create time',
  memberList: 'Member list',
  usernamePlaceholder: 'please enter user name',
  batchOperate: 'batch operation',
  addMember: 'add member',
  selected: 'selected',
  item: 'item',
  clear: 'clear',
  batchDelete: 'batch delete',
  joinTime: 'join time',
  permission: 'permission',
  username: 'username',
  readonly: 'read only',
  write: 'read and write',
  search: 'search',
  all: 'all',
  selectAll: 'select all',
  submit: 'submit',
};

// 实验
const experiment = {
  'experiment.table.noAuth.description':
    'You have read-only permission in this space, and the creation of experiments is not currently supported. If you want to create an experiment, please go to the member management to find a member with read and write permissions in the space to modify the permissions.',
  'experiment.table.title':
    'There is no experimental data in the current space',
  'experiment.table.description':
    'Please create an experiment first. You can choose to create your own experiments or quickly build experimental scenarios by recommending experiments to verify the reliability of the application system',
  'experiment.delete.title': 'Are you sure you want to delete this experiment?',
  'experiment.delete.content':
    'Deleting an experiment will delete the configuration of the experiment, but will not delete the historical experiment results!',
};

// 实验结果
const experimentResult = {
  'experimentResult.table.description':
    'There is currently no experimental result data',
  'experimentResult.table.title':
    'You can run the experiment by going to the experiment details page',
  'experimentResult.table.noAuth.title':
    'You can view experiments by going to the experiment list page',
  'experimentResult.table.btn': 'Go to experiment list',
  'experimentResult.stop.text': 'The experiment has stopped',
};

// 新增/编辑实验页面
const addExperiment = {
  'addExperiment.basic.tip': 'Please complete basic information',
  'addExperiment.node.tip': 'Please complete the node information',
};
// tag抽屉
const tag = {
  'tag.repeat.text': 'Tag already exists, please re-enter it',
  'tag.repeat.check': 'Tag already exists',
  'tag.empty.tip': 'Please add tags',
  'tag.create.success.tip': 'You have successfully created a new label',
};

// 节点库
const nodeLibrary = {
  'nodeLibrary.title': 'node library',
};

// 空间下列表
const spaceDropdown = {
  'spaceDropdown.tip': 'no related space? view ',
};

// 创建空间
const createSpace = {
  'createSpace.confirm': 'Complete creation and go to configuration',
};

// 空间下空间设置
const spaceSetting = {
  'spaceSetting.tab.member': 'Member management',
  'spaceSetting.tab.tag': 'Tag management',
};

// 空间下成员管理
const memberManageMent = {
  'memberManageMent.readonly.tip':
    'have editing permissions in the current space',
  'memberManageMent.write.tip': 'can only be viewed in the current space',
  'memberManageMent.delete.title':
    'Are you sure you want to delete the currently selected member?',
  'memberManageMent.delete.content':
    'If you delete a member in a space, the member will not be able to enter the space!',
  'memberManageMent.delete.success.tip':
    'You have successfully deleted the selected members!',
  'memberManageMent.permission.success.tip':
    'Permissions modified successfully!',
};

// 空间下添加成员
const addMember = {
  'addMember.user.permission': 'user permissions',
  'addMember.selected.user': 'user selected',
  'addMember.noMore': 'no more～',
  'addMember.loading': 'loading...',
  'addMember.search.result':
    'No results found yet, please try switching keywords and searching again.',
  'addMember.add.success': 'added successfully',
};

// 空间下标签管理
const tagManageMent = {
  'tagManageMent.title': 'Tag list',
  'tagManageMent.add.title': 'create new label',
  'tagManageMent.search.placeholder': 'please enter a label name',
  'tagManageMent.table.empty.title': "You don't have label data yet",
  'tagManageMent.table.empty.description':
    'Please create a new label. If you create a label in advance, you can directly choose to quickly label the experiment when creating an experiment.',
  'tagManageMent.table.empty.noAuth.title': 'There is currently no label data',
  'tagManageMent.table.empty.noAuth.description':
    'You only have read-only permissions in this space, and adding tags is not currently supported. If you want to add a tag, please go to the member management to find a member with read and write permissions in the space to modify the permissions',
  'tagManageMent.delete.success.tip': 'You have successfully deleted the tag',
  'tagManageMent.column.tagColor': 'tag color',
  'tagManageMent.column.tagStyle': 'tag style',
};

// 登录相关
const login = {
  login: 'Sign in',
  register: 'Sign up',
  signOut: 'sign out',
  updatePassword: 'modify password',
  password: 'password',
  'password.old': 'original password',
  'password.old.placeholder': 'Please enter the original password',
  'password.new': 'new password',
  'password.new.placeholder': 'Please enter a new password',
  'password.new.placeholder.again': 'Please enter new password again',
  'password.rule':
    'Password 8-16 characters Chinese, English, uppercase and lowercase characters and special characters such as underscores',
  'password.confirm': 'Please confirm your password',
  'password.confirm.again': 'Confirm Password',
  'password.error': 'Incorrect password',
  'password.new.confirm': 'Confirm the new password',
  'password.success':
    'The password has been changed successfully. You will be redirected to the login page to log in again.',
  'password.placeholder': 'Please enter password',
  'password.inconsistent': 'Two passwords are inconsistent',
  welcome: 'welcome to',
  'username.rule':
    'The username can be in Chinese and English, and the length should not exceed 64 characters.',
  notAccount: 'Don’t have an account yet?',
  haveAccount: 'Already have an account?',
  'reister.success': 'Registration successful, please log in'
};
// 账号相关
const account = {
  'account.title': 'Account management',
  'account.list': 'Account list',
  'account.search.placeholder': 'Please enter username to search',
  role: 'role',
  admin: 'admin',
  generalUser: 'general user',
  adminDescription: 'Have all permissions',
  generalUserDescription: 'You can log in to view, and superimpose the permissions in the space.',
  'account.delete.title': 'Are you sure you want to delete the currently selected account?',
  'account.delete.content': 'Users who delete their accounts will not be able to log in to the platform, and they can only re-register if they want to use it again!',
  'account.delete.success': 'You have successfully deleted the selected members',
  'account.role.update': 'User role modified successfully'
}

// 空间管理
const spaceManagement = {
  'spaceManagement.title': 'Space management',
  'spaceManagement.spaceName.placeholder': 'Please enter the space name',
  'spaceManagement.spaceMember.placeholder': 'Please enter space members',
  'spaceManagement.alert': 'Members with read and write permissions in the contact space can be added as space members',
  'spaceManagement.member': 'space member',
  'spaceManagement.tab.all': 'All spaces',
  'spaceManagement.tab.related': 'My related',
  'spaceManagement.delete.title': 'Are you sure you want to delete the currently selected space?',
  'spaceManagement.delete.success': 'You have successfully deleted the selected space',
  'spaceManagement.noAuth.tip': 'You do not have permission to this space, please contact a read-write member',
  'spaceManagement.noAuth.readonly.tip': 'Read-only users are currently unable to use this feature',
  'spaceManagement.write': 'read and write members',
  'spaceManagement.experimentCount': 'number of experiments',
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
