import { envType } from '@/constants';
import request from '@/utils/request';

/**
 * 获取namespace列表
 * @param params
 * @param options
 * @returns
 */
export async function queryNamespaceList(
  params?: {
    page?: number;
    page_size?: number;
  },
  options?: { [key: string]: any },
) {
  return request<any>(
    `/chaosmeta/api/v1/kubernetes/cluster/${envType}/namespaces`,
    {
      method: 'GET',
      params,
      ...(options || {}),
    },
  );
}

/**
 * 获取pod列表
 * @param params
 * @param options
 * @returns
 */
export async function queryPodNameList(
  params?: {
    page?: number;
    page_size?: number;
    namespace?: string;
  },
  options?: { [key: string]: any },
) {
  return request<any>(
    `/chaosmeta/api/v1/kubernetes/cluster/${envType}/namespace/${params?.namespace}/pods`,
    {
      method: 'GET',
      params,
      ...(options || {}),
    },
  );
}

/**
 * 获取nodename列表
 * @param params
 * @param options
 * @returns
 */
export async function queryNodeNameList(
  params?: {
    page?: number;
    page_size?: number;
    namespace?: string;
  },
  options?: { [key: string]: any },
) {
  return request<any>(`/chaosmeta/api/v1/kubernetes/cluster/${envType}/nodes`, {
    method: 'GET',
    params,
    ...(options || {}),
  });
}

/**
 * 获取故障节点详情信息，属于那个target下的
 * @param params
 * @param options
 * @returns
 */
export async function queryFaultNodeDetail(
  params?: {
    targetId?: number;
  },
  options?: { [key: string]: any },
) {
  return request<any>(
    `/chaosmeta/api/v1/injects/scopes/target/${params?.targetId}`,
    {
      method: 'GET',
      params,
      ...(options || {}),
    },
  );
}


/**
 * 获取deploymentName列表
 * @param params
 * @param options
 * @returns
 */
export async function queryDeploymentNameList(
  params?: {
    page?: number;
    page_size?: number;
    namespace?: string;
  },
  options?: { [key: string]: any },
) {
  return request<any>(
    `/chaosmeta/api/v1/kubernetes/cluster/${envType}/namespace/${params?.namespace}/deployments`,
    {
      method: 'GET',
      params,
      ...(options || {}),
    },
  );
}

/**
 * 获取containersName列表
 * @param namespace
 * @param data
 * @param options
 * @returns
 */
export async function queryContainersNameList(
  namespace: string,
  data?: {
    target_label?: string;
    target_pods?: string;
  },
  options?: { [key: string]: any },
) {
  return request<any>(
    `/chaosmeta/api/v1/kubernetes/cluster/${envType}/namespace/${namespace}//containers`,
    {
      method: 'POST',
      data: data || {},
      ...(options || {}),
    },
  );
}
