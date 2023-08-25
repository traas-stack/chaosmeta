import request from '@/utils/request';

/**
 * 创建空间
 * @param body
 * @param options
 * @returns
 */
export async function createSpace(
  body: {
    name: string;
    description?: string;
  },
  options?: { [key: string]: any },
) {
  return request<any>(`/chaosmeta/api/v1/namespaces`, {
    method: 'POST',
    data: body,
    ...(options || {}),
  });
}

/**
 * 获取空间列表
 * @param params
 * @param options
 * @returns
 */
export async function querySpaceList(
  params?: {
    sort?: string;
    name?: string;
    creator?: string;
    page?: number;
    page_size?: number;
  },
  options?: { [key: string]: any },
) {
  return request<any>(`/chaosmeta/api/v1/namespaces/list`, {
    method: 'GET',
    params,
    ...(options || {}),
  });
}

/**
 * 获取空间列表带分类
 * @param params
 * @param options
 * @returns
 */
export async function queryClassSpaceList(
  params?: {
    sort?: string;
    name?: string;
    creator?: string;
    page?: number;
    page_size?: number;
    namespaceClass?: string;
  },
  options?: { [key: string]: any },
) {
  return request<any>(`/chaosmeta/api/v1/namespaces/query`, {
    method: 'GET',
    params,
    ...(options || {}),
  });
}

/**
 * 修改空间基础信息
 * @param body
 * @param options
 * @returns
 */
export async function editSpaceBasic(
  body: {
    id: number;
    name: string;
    description?: string;
  },
  options?: { [key: string]: any },
) {
  return request<any>(`/chaosmeta/api/v1/namespaces/${body?.id}`, {
    method: 'POST',
    data: body,
    ...(options || {}),
  });
}

/**
 * 获取空间详情
 * @param params
 * @param options
 * @returns
 */
export async function querySpaceDetail(
  params: { id: number | string },
  options?: { [key: string]: any },
) {
  return request<any>(`/chaosmeta/api/v1/namespaces/${params?.id}`, {
    method: 'GET',
    params,
    ...(options || {}),
  });
}

/**
 * 删除空间
 * @param params
 * @param options
 * @returns
 */
export async function deleteSpace(
  params: {
    id: number;
  },
  options?: { [key: string]: any },
) {
  return request<any>(`/chaosmeta/api/v1/namespaces/${params.id}`, {
    method: 'DELETE',
    params: { ...params },
    ...(options || {}),
  });
}

/**
 * 空间下成员列表
 * @param params
 * @param options
 * @returns
 */
export async function querySpaceUserList(
  params: {
    id: number;
    sort?: string;
    name?: string;
    username?: string;
    page?: number;
    page_size?: number;
  },
  options?: { [key: string]: any },
) {
  return request<any>(`/chaosmeta/api/v1/namespaces/${params.id}/users`, {
    method: 'GET',
    params,
    ...(options || {}),
  });
}

/**
 * 空间下批量添加成员
 * @param body
 * @param options
 * @returns
 */
export async function spaceAddUser(
  body: {
    id: number | string;
    users: API.Query_AddUser[];
  },
  options?: { [key: string]: any },
) {
  return request<any>(`/chaosmeta/api/v1/namespaces/${body.id}/users/batch`, {
    method: 'POST',
    data: { users: body.users },
    ...(options || {}),
  });
}

/**
 * 空间下移除成员
 * @param params
 * @param options
 * @returns
 */
export async function spaceDeleteUser(
  params: {
    id: number;
    user_id: number;
  },
  options?: { [key: string]: any },
) {
  const { id, user_id } = params;
  return request<any>(`/chaosmeta/api/v1/namespaces/${id}/users/${user_id}`, {
    method: 'DELETE',
    params: { ...params },
    ...(options || {}),
  });
}

/**
 * 空间下批量移除用户
 * @param params
 * @param options
 * @returns
 */
export async function spaceBatchDeleteUser(
  body: {
    id: number;
    user_ids: number[];
  },
  options?: { [key: string]: any },
) {
  const { id, user_ids } = body;
  return request<any>(`/chaosmeta/api/v1/namespaces/${id}/users`, {
    method: 'DELETE',
    data: { user_ids },
    ...(options || {}),
  });
}

/**
 * 更改成员在空间内的权限
 * @param body
 * @param options
 * @returns
 */
export async function spaceModifyUserPermission(
  body: {
    id: number | string;
    user_ids: number[];
    permission: number | string;
  },
  options?: { [key: string]: any },
) {
  const { id, user_ids, permission } = body;
  return request<any>(`/chaosmeta/api/v1/namespaces/${id}/users/permission`, {
    method: 'POST',
    data: {
      permission,
      user_ids,
    },
    ...(options || {}),
  });
}

/**
 * 设置空间内可攻击集群
 * @param body
 * @param options
 * @returns
 */
export async function spaceSettingCluster(
  body: {
    id: number;
    cluster_id: number;
  },
  options?: { [key: string]: any },
) {
  return request<any>(`/chaosmeta/api/v1/namespaces/${body.id}/cluster`, {
    method: 'POST',
    data: body,
    ...(options || {}),
  });
}

/**
 * 查询空间内可攻击集群
 * @param body
 * @param options
 * @returns
 */
export async function querySpaceCluster(
  params: {
    id: number;
  },
  options?: { [key: string]: any },
) {
  return request<any>(`/chaosmeta/api/v1/namespaces/${params.id}/cluster`, {
    method: 'GET',
    params,
    ...(options || {}),
  });
}

/**
 * 添加空间内可攻击hosts
 * @param body
 * @param options
 * @returns
 */
export async function spaceAddHost(
  body: {
    id: number;
    host_ids: number[];
  },
  options?: { [key: string]: any },
) {
  return request<any>(`/chaosmeta/api/v1/namespaces/${body.id}/hosts`, {
    method: 'POST',
    data: body,
    ...(options || {}),
  });
}

/**
 * 删除空间内可攻击hosts
 * @param params
 * @param options
 * @returns
 */
export async function spaceDeleteHost(
  params: {
    id: number;
    host_ids: number[];
  },
  options?: { [key: string]: any },
) {
  const { id } = params;
  return request<any>(`/chaosmeta/api/v1/namespaces/${id}/hosts`, {
    method: 'DELETE',
    params: { ...params },
    ...(options || {}),
  });
}

/**
 * 查询空间内可攻击hosts
 * @param params
 * @param options
 * @returns
 */
export async function querySpaceHosts(
  params: {
    id: number;
  },
  options?: { [key: string]: any },
) {
  return request<any>(`/chaosmeta/api/v1/namespaces/${params.id}/hosts`, {
    method: 'GET',
    params,
    ...(options || {}),
  });
}

/**
 * 查询空间内的标签列表
 * @param params
 * @param options
 * @returns
 */
export async function querySpaceTagList(
  params: {
    id: number;
    name?: string;
    page?: number;
    pageSize?: number;
    creator?: string;
  },
  options?: { [key: string]: any },
) {
  return request<any>(`/chaosmeta/api/v1/namespaces/${params.id}/labels`, {
    method: 'GET',
    params,
    ...(options || {}),
  });
}

/**
 * 空间内添加标签
 * @param body
 * @param options
 * @returns
 */
export async function spaceAddTag(
  body: {
    id: number | string;
    name: string;
    color: string;
  },
  options?: { [key: string]: any },
) {
  const { id, name, color } = body;
  return request<any>(`/chaosmeta/api/v1/namespaces/${id}/labels`, {
    method: 'POST',
    data: {
      name,
      color,
    },
    ...(options || {}),
  });
}

/**
 * 查询标签信息，可用于校验标签是否已添加
 * @param params
 * @param options
 * @returns
 */
export async function querySpaceTagName(
  params: {
    ns_id: number | string;
    name: string;
  },
  options?: { [key: string]: any },
) {
  const { ns_id, name } = params;
  return request<any>(`/chaosmeta/api/v1/namespaces/${ns_id}/labels/${name}`, {
    method: 'GET',
    params,
    ...(options || {}),
  });
}

/**
 * 删除空间内标签
 * @param params
 * @param options
 * @returns
 */
export async function spaceDeleteTag(
  params: {
    id: number;
    ns_id: number | string;
  },
  options?: { [key: string]: any },
) {
  const { id, ns_id } = params;
  return request<any>(`/chaosmeta/api/v1/namespaces/${ns_id}/labels/${id}`, {
    method: 'DELETE',
    params: { ...params },
    ...(options || {}),
  });
}
