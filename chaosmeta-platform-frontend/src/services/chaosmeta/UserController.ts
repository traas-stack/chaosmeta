import request from '@/utils/request';

/**
 * 用户注册
 * @param body
 * @param options
 * @returns
 */
export async function register(
  body: {
    name: string;
    password: string;
  },
  options?: any,
) {
  return request<any>('/users/token/create', {
    method: 'POST',
    data: body,
    headers: {
      'Content-Type': 'application/json',
    },
    ...(options || {}),
  });
}

/**
 * 用户登录
 * @param body
 * @param options
 * @returns
 */
export async function login(
  body: {
    name: string;
    password: string;
  },
  options?: any,
) {
  return request<any>('/users/token/login', {
    method: 'POST',
    data: body,
    headers: {
      'Content-Type': 'application/json',
    },
    ...(options || {}),
  });
}

/**
 * 获取用户列表
 * @param params
 * @param options
 * @returns
 */
export async function getUserList(
  params?: {
    sort?: string;
    name?: string;
    role?: string;
    offset?: number;
    limit?: number;
  },
  options?: { [key: string]: any },
) {
  return request<any>('/chaosmeta/api/v1/users/list', {
    method: 'GET',
    params,
    ...(options || {}),
  });
}

/**
 * 获取单个用户信息
 * @param params
 * @param options
 * @returns
 */
export async function getUserInfo(
  params?: {
    name?: string;
  },
  options?: { [key: string]: any },
) {
  return request<any>(`/chaosmeta/api/v1/users/${params?.name}`, {
    method: 'GET',
    params,
    ...(options || {}),
  });
}

/**
 * 更新token
 * @param body
 * @param options
 * @returns
 */
export async function updateToken(
  body?: any,
  options?: { [key: string]: any },
) {
  return request<any>(`/users/token/refresh`, {
    method: 'POST',
    data: body,
    ...(options || {}),
  });
}

/**
 * 用户修改密码
 * @param body
 * @param options
 * @returns
 */
export async function updatePassword(
  body?: {
    password: string;
  },
  options?: { [key: string]: any },
) {
  return request<any>(`/chaosmeta/api/v1/users/password`, {
    method: 'POST',
    data: body,
    ...(options || {}),
  });
}

/**
 * 删除单个用户/账号
 * @param params
 * @param options
 * @returns
 */
export async function deleteUser(
  params: {
    // path
    /** userId */
    userId?: string;
  },
  options?: { [key: string]: any },
) {
  const { userId } = params;
  return request<any>(`/chaosmeta/api/v1/users/${userId}`, {
    method: 'DELETE',
    params: { ...params },
    ...(options || {}),
  });
}

/**
 * 批量删除用户/账号
 * @param params
 * @param options
 * @returns
 */
export async function batchDeleteUser(
  body: {
    /** user_ids */
    user_ids?: string[];
  },
  options?: { [key: string]: any },
) {
  return request<any>(`/chaosmeta/api/v1/users`, {
    method: 'DELETE',
    data: { ...body },
    ...(options || {}),
  });
}

/**
 * 批量修改用户身份
 * @param body
 * @param options
 * @returns
 */
export async function changeUserRole(
  body?: {
    user_ids: number[];
    role: string;
  },
  options?: { [key: string]: any },
) {
  return request<any>(`/chaosmeta/api/v1/users/role`, {
    method: 'POST',
    data: body,
    ...(options || {}),
  });
}

/** 根据表单key获取表单
@param sceneType
@return Result<FormEntity>
 GET /chaos/form/list */
export async function list(
  params: {
    // query
    /** 攻击场景类型 */
    sceneType?: any;
  },
  options?: { [key: string]: any },
) {
  return request<any>('/chaos/form/list', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}
