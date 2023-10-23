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
export async function queryPodLIst(
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
 * 获取podnode列表
 * @param params
 * @param options
 * @returns
 */
export async function queryPodNodeList(
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
