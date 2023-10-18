import request from '@/utils/request';

/**
 * 获取实验列表
 * @param params
 * @param options
 * @returns
 */
export async function queryExperimentList(
  params?: {
    sort?: string;
    name?: string;
    creator?: string;
    page?: number;
    page_size?: number;
    create_time?: string;
    time_search_field?: string;
    start_time?: string;
    end_time?: string;
    namespace_id: string;
  },
  options?: { [key: string]: any },
) {
  return request<any>(`/chaosmeta/api/v1/experiments`, {
    method: 'GET',
    params,
    ...(options || {}),
  });
}

/**
 * 获取实验详情
 * @param params
 * @param options
 * @returns
 */
export async function queryExperimentDetail(
  params?: {
    uuid?: string;
  },
  options?: { [key: string]: any },
) {
  return request<any>(`/chaosmeta/api/v1/experiments/${params?.uuid}`, {
    method: 'GET',
    params,
    ...(options || {}),
  });
}

/**
 * 创建实验
 * @param body
 * @param options
 * @returns
 */
export async function createExperiment(
  body: any,
  options?: { [key: string]: any },
) {
  return request<any>(`/chaosmeta/api/v1/experiments`, {
    method: 'POST',
    data: body,
    ...(options || {}),
  });
}

/**
 * 更新实验
 * @param body
 * @param options
 * @returns
 */
export async function updateExperiment(
  body: any,
  options?: { [key: string]: any },
) {
  const { uuid } = body;
  return request<any>(`/chaosmeta/api/v1/experiments/${uuid}`, {
    method: 'POST',
    data: body,
    ...(options || {}),
  });
}

/**
 * 运行实验
 * @param body
 * @param options
 * @returns
 */
export async function runExperiment(
  body: {
    uuid: string;
  },
  options?: { [key: string]: any },
) {
  const { uuid } = body;
  return request<any>(`/chaosmeta/api/v1/experiments/${uuid}/start`, {
    method: 'POST',
    data: body,
    ...(options || {}),
  });
}

/**
 * 删除实验
 * @param params
 * @param options
 * @returns
 */
export async function deleteExperiment(
  params: {
    uuid: string;
  },
  options?: { [key: string]: any },
) {
  return request<any>(`/chaosmeta/api/v1/experiments/${params.uuid}`, {
    method: 'DELETE',
    params,
    ...(options || {}),
  });
}

/**
 * 获取实验结果列表
 * @param params
 * @param options
 * @returns
 */
export async function queryExperimentResultList(
  params?: {
    sort?: string;
    name?: string;
    creator_name?: string;
    page?: number;
    page_size?: number;
    experiment_uuid?: string;
    status?: string;
    time_search_field?: string;
    start_time?: string;
    end_time?: string;
    namespace_id: string;
  },
  options?: { [key: string]: any },
) {
  return request<any>(`/chaosmeta/api/v1/experiments/results`, {
    method: 'GET',
    params,
    ...(options || {}),
  });
}

/**
 * 停止实验结果
 * @param body
 * @param options
 * @returns
 */
export async function stopExperimentResult(
  body: {
    // 实验结果的uuid
    uuid: string;
  },
  options?: { [key: string]: any },
) {
  return request<any>(`/chaosmeta/api/v1/experiments/${body?.uuid}/stop`, {
    method: 'POST',
    data: body,
    ...(options || {}),
  });
}

/**
 * 获取实验结果详情
 * @param params
 * @param options
 * @returns
 */
export async function queryExperimentResultDetail(
  params?: {
    uuid?: string;
  },
  options?: { [key: string]: any },
) {
  return request<any>(`/chaosmeta/api/v1/experiments/results/${params?.uuid}`, {
    method: 'GET',
    params,
    ...(options || {}),
  });
}

/**
 * 获取实验结果的编排节点列表
 * @param params
 * @param options
 * @returns
 */
export async function queryExperimentResultArrangeNodeList(
  params?: {
    uuid?: string;
  },
  options?: { [key: string]: any },
) {
  return request<any>(
    `/chaosmeta/api/v1/experiments/results/${params?.uuid}/nodes`,
    {
      method: 'GET',
      params,
      ...(options || {}),
    },
  );
}

/**
 * 获取实验结果的编排节点详情
 * @param params
 * @param options
 * @returns
 */
export async function queryExperimentResultArrangeNodeDetail(
  params?: {
    uuid?: string;
    node_id?: string;
  },
  options?: { [key: string]: any },
) {
  return request<any>(
    `/chaosmeta/api/v1/experiments/results/${params?.uuid}/nodes/${params?.node_id}`,
    {
      method: 'GET',
      params,
      ...(options || {}),
    },
  );
}

/**
 * 获取实验结果的编排节点单实例的执行详情
 * @param params
 * @param options
 * @returns
 */
export async function queryExperimentResultArrangeNodeExecuteDetail(
  params: {
    uuid: string;
    node_id: string;
    subtask_id: string;
  },
  options?: { [key: string]: any },
) {
  const { uuid, node_id, subtask_id } = params;
  return request<any>(
    `/chaosmeta/api/v1/experiments/results/${uuid}/nodes/${node_id}/subtasks/${subtask_id}`,
    {
      method: 'GET',
      params,
      ...(options || {}),
    },
  );
}

/**
 * 删除实验结果
 * @param params
 * @param options
 * @returns
 */
export async function deleteExperimentResult(
  params: {
    uuid: number;
  },
  options?: { [key: string]: any },
) {
  return request<any>(`/chaosmeta/api/v1/experiments/results/${params?.uuid}`, {
    method: 'DELETE',
    params: { ...params },
    ...(options || {}),
  });
}

/**
 * 批量删除实验结果
 * @param params
 * @param options
 * @returns
 */
export async function batchDeleteExperimentResult(
  params: {
    result_uuids: string[];
  },
  options?: { [key: string]: any },
) {
  return request<any>(`/chaosmeta/api/v1/experiments/results`, {
    method: 'DELETE',
    params: { ...params },
    ...(options || {}),
  });
}

/**
 * 故障节点下 - 查询故障注入scope - 一级节点
 * @param params
 * @param options
 * @returns
 */
export async function queryFaultNodeScopes(
  params?: any,
  options?: { [key: string]: any },
) {
  return request<any>(`/chaosmeta/api/v1/injects/scopes`, {
    method: 'GET',
    params,
    ...(options || {}),
  });
}

/**
 * 故障节点下 - 查询故障注入target - 二级节点 需根据一级节点id查询
 * @param params
 * @param options
 * @returns
 */
export async function queryFaultNodeTargets(
  params?: {
    id: number;
  },
  options?: { [key: string]: any },
) {
  return request<any>(
    `/chaosmeta/api/v1/injects/scopes/${params?.id}/targets`,
    {
      method: 'GET',
      params,
      ...(options || {}),
    },
  );
}

/**
 * 故障节点下 - 查询故障注入能力列表-节点 - 三级节点 需根据一级&二级节点id查询
 * @param params
 * @param options
 * @returns
 */
export async function queryFaultNodeItem(
  params: {
    scope_id: number;
    target_id: number;
  },
  options?: { [key: string]: any },
) {
  const { scope_id, target_id } = params;
  return request<any>(
    `/chaosmeta/api/v1/injects/scopes/${scope_id}/targets/${target_id}/faults`,
    {
      method: 'GET',
      params,
      ...(options || {}),
    },
  );
}

/**
 * 故障节点下 - 根据节点id查询节点表单配置信息 - 动态表单
 * @param params
 * @param options
 * @returns
 */
export async function queryFaultNodeFields(
  params: {
    id: number;
  },
  options?: { [key: string]: any },
) {
  const { id } = params;
  return request<any>(`/chaosmeta/api/v1/injects/faults/${id}/args`, {
    method: 'GET',
    params,
    ...(options || {}),
  });
}

/**
 * 度量注入能力列表 - 一级节点
 * @param params
 * @param options
 * @returns
 */
export async function queryMeasureList(
  params?: any,
  options?: { [key: string]: any },
) {
  return request<any>(`/chaosmeta/api/v1/injects/measures`, {
    method: 'GET',
    params,
    ...(options || {}),
  });
}

/**
 * 度量注入能力列表 - 查询度量注入能力节点 - 根据节点id查询节点表单配置信息 - 动态表单
 * @param params
 * @param options
 * @returns
 */
export async function queryMeasureNodeFields(
  params?: {
    id: number;
  },
  options?: { [key: string]: any },
) {
  return request<any>(`/chaosmeta/api/v1/injects/measures/${params?.id}/args`, {
    method: 'GET',
    params,
    ...(options || {}),
  });
}

/**
 * 流量注入能力列表 - 一级节点
 * @param params
 * @param options
 * @returns
 */
export async function queryFlowList(
  params?: any,
  options?: { [key: string]: any },
) {
  return request<any>(`/chaosmeta/api/v1/injects/flows`, {
    method: 'GET',
    params,
    ...(options || {}),
  });
}

/**
 * 流量注入能力列表 - 查询流量注入能力列表节点 - 根据节点id查询节点表单配置信息 - 动态表单
 * @param params
 * @param options
 * @returns
 */
export async function queryFlowNodeFields(
  params?: {
    id: number;
  },
  options?: { [key: string]: any },
) {
  return request<any>(`/chaosmeta/api/v1/injects/flows/${params?.id}/args`, {
    method: 'GET',
    params,
    ...(options || {}),
  });
}
