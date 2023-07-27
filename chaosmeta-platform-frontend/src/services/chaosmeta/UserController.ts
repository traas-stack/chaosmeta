/* eslint-disable */
// 该文件由 OneAPI 自动生成，请勿手动修改！
import { request } from '@umijs/max';

/** 此处后端没有提供注释 GET /api/v1/queryUserList */
export async function queryUserList(
  params: {
    // query
    /** keyword */
    keyword?: string;
    /** current */
    current?: number;
    /** pageSize */
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  return request<API.Result_PageInfo_UserInfo__>('/api/v1/queryUserList', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 此处后端没有提供注释 POST /api/v1/user */
export async function addUser(
  body?: API.UserInfoVO,
  options?: { [key: string]: any },
) {
  return request<API.Result_UserInfo_>('/api/v1/user', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 此处后端没有提供注释 GET /api/v1/user/${param0} */
export async function getUserDetail(
  params: {
    // path
    /** userId */
    userId?: string;
  },
  options?: { [key: string]: any },
) {
  const { userId: param0 } = params;
  return request<API.Result_UserInfo_>(`/api/v1/user/${param0}`, {
    method: 'GET',
    params: { ...params },
    ...(options || {}),
  });
}

/** 此处后端没有提供注释 PUT /api/v1/user/${param0} */
export async function modifyUser(
  params: {
    // path
    /** userId */
    userId?: string;
  },
  body?: API.UserInfoVO,
  options?: { [key: string]: any },
) {
  const { userId: param0 } = params;
  return request<API.Result_UserInfo_>(`/api/v1/user/${param0}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    params: { ...params },
    data: body,
    ...(options || {}),
  });
}

/** 此处后端没有提供注释 DELETE /api/v1/user/${param0} */
export async function deleteUser(
  params: {
    // path
    /** userId */
    userId?: string;
  },
  options?: { [key: string]: any },
) {
  const { userId: param0 } = params;
  return request<API.Result_string_>(`/api/v1/user/${param0}`, {
    method: 'DELETE',
    params: { ...params },
    ...(options || {}),
  });
}

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

// /users/token

export async function tokenState(
  params?: any,
  options?: { [key: string]: any },
) {
  return request<API.Result_PageInfo_UserInfo__>('/users/token', {
    method: 'GET',
    params,
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
  return request<API.Result_PageInfo_UserInfo__>(
    '/chaosmeta/api/v1/users/list',
    {
      method: 'GET',
      params,
      ...(options || {}),
    },
  );
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
  return request<API.Result_PageInfo_UserInfo__>(
    `/chaosmeta/api/v1/users/${params?.name}`,
    {
      method: 'GET',
      params,
      ...(options || {}),
    },
  );
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
  return request<API.Result_PageInfo_UserInfo__>(`/users/token/refresh`, {
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
