import {
  queryFaultNodeItem,
  queryFaultNodeScopes,
  queryFaultNodeTargets,
  queryFlowList,
  queryMeasureList,
} from '@/services/chaosmeta/ExperimentController';
import { getIntlName } from '@/utils/format';
import { useSortable } from '@dnd-kit/sortable';
import { useIntl } from '@umijs/max';
import { Tooltip, Tree } from 'antd';
import React, { useState } from 'react';
import { NodeItem, NodeLibraryContainer } from '../style';

interface IProps {
  disabled?: boolean;
}

// 初始化节点, 度量和流量暂时禁用，后期加入进来
const initTreeData: any[] = [
  { nameCn: '故障节点', key: 'fault', name: 'faulty node' },
  {
    nameCn: '度量引擎',
    key: 'measure',
    name: 'measurement engine',
  },
  {
    nameCn: '流量注入',
    key: 'flow',
    name: 'flow injection',
  },
  {
    nameCn: '其他节点',
    name: 'other nodes',
    key: 'other',
    children: [
      // 其他节点下默认写死等待时长节点
      {
        nameCn: '等待时长',
        name: 'waiting time',
        key: 'wait',
        targetId: 'wait-init',
        isLeaf: true,
        exec_type: 'wait',
        exec_type_name: '等待时长',
        nodeInfoState: true,
      },
    ],
  },
];

/**
 * 节点库
 */
const NodeLibrary: React.FC<IProps> = (props) => {
  const { disabled = false } = props;
  const [treeData, setTreeData] = useState(initTreeData);
  const [expandedKeys, setExpandedKeys] = useState<string[]>([]);
  const intl = useIntl();

  /**
   * 更新节点数据
   * @param list
   * @param key
   * @param children
   * @returns
   */
  const updateTreeData = (
    list: any,
    key: React.Key,
    children: any[],
  ): any[] => {
    return list.map((node: { key: React.Key; children: any }) => {
      if (node.key === key) {
        return {
          ...node,
          children,
        };
      }
      if (node.children) {
        return {
          ...node,
          children: updateTreeData(node.children, key, children),
        };
      }
      return node;
    });
  };

  /**
   * 左侧节点渲染
   */
  const LeftNodeItem = (params: any) => {
    const { itemData, disabledItem } = params;
    // 用于绑定拖拽，listeners需配置到拖动元素上，拖动就靠它
    const { setNodeRef, listeners, isDragging } = useSortable({
      id: itemData?.key,
      disabled: disabledItem,
      // 额外数据，用于悬浮态数据的渲染和判断
      data: {
        dragtype: 'node',
        isNode: true,
        ...itemData,
      },
    });

    return (
      <NodeItem
        ref={setNodeRef}
        $isDragging={isDragging}
        $disabledItem={disabledItem}
      >
        <Tooltip title={getIntlName(itemData)}>
          <div {...listeners} className="temp-item ellipsis">
            <img src="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*rOAzRrDGQoAAAAAAAAAAAAAADmKmAQ/original" />
            {getIntlName(itemData)}
          </div>
        </Tooltip>
      </NodeItem>
    );
  };

  const onLoadData = (params: any) => {
    const { key, id, scopeId } = params;
    if (key === 'other') {
      return Promise.resolve();
    }
    // 故障节点 - 一级节点list查询
    if (key === 'fault') {
      return queryFaultNodeScopes()?.then((res: any) => {
        if (res?.code === 200) {
          const formatList = res?.data?.scopes;
          formatList?.forEach((item: { key: string; id: number }) => {
            item.key = `${key}-${item.id}`;
          });
          setTreeData((origin) => {
            return updateTreeData(origin, key, formatList);
          });
        }
      });
    }
    // 故障节点 - 二级节点list查询
    if (key && !scopeId && key.includes('fault')) {
      return queryFaultNodeTargets({ id }).then((res: any) => {
        if (res?.code === 200) {
          const targetList = res?.data?.targets;
          targetList?.forEach((item: { key: string; id: number }) => {
            item.key = `${key}-${item.id}`;
          });
          setTreeData((origin) => {
            return updateTreeData(origin, key, targetList);
          });
        }
      });
    }
    // 故障节点 - 三级节点list查询
    if (scopeId && key.includes('fault')) {
      return queryFaultNodeItem({ scope_id: scopeId, target_id: id }).then(
        (res: any) => {
          if (res?.code === 200) {
            const faults = res?.data?.faults;
            // 设置一些编排配置需要的参数
            const faultParams = {
              isLeaf: true,
              // 一级节点的id
              scope_id: scopeId,
              // 二级节点的id
              target_id: id,
              exec_type: 'fault',
              exec_type_name: '故障节点',
              // 当前节点信息是否进行配置完成，默认未完成，需要进行配置
              nodeInfoState: false,
            };
            const newFaults = faults?.map((item: any) => {
              return { ...item, key: `${key}-${item.id}`, ...faultParams };
            });
            setTreeData((origin) => {
              return updateTreeData(origin, key, newFaults);
            });
          }
        },
      );
    }
    // 度量引擎
    if (key === 'measure') {
      return queryMeasureList()?.then((res: any) => {
        if (res?.code === 200) {
          const measures = res?.data?.measures;
          // 设置一些编排配置需要的参数
          const measuresParams = {
            isLeaf: true,
            isChildrenNode: true,
            exec_type: 'measure',
            exec_type_name: '度量引擎',
            // 当前节点信息是否进行配置完成，默认未完成，需要进行配置
            nodeInfoState: false,
          };
          const newFaults = measures?.map((item: any) => {
            return { ...item, key: `${key}-${item.id}`, ...measuresParams };
          });
          setTreeData((origin) => {
            return updateTreeData(origin, key, newFaults);
          });
        }
      });
    }
    if (key === 'flow') {
      return queryFlowList()?.then((res: any) => {
        if (res?.code === 200) {
          const flows = res?.data?.flows;
          // 设置一些编排配置需要的参数
          const flowsParams = {
            isLeaf: true,
            isChildrenNode: true,
            exec_type: 'flow',
            exec_type_name: '流量注入',
            // 当前节点信息是否进行配置完成，默认未完成，需要进行配置
            nodeInfoState: false,
          };
          const newFaults = flows?.map((item: any) => {
            return { ...item, key: `${key}-${item.id}`, ...flowsParams };
          });
          setTreeData((origin) => {
            return updateTreeData(origin, key, newFaults);
          });
        }
      });
    }
    return Promise.resolve();
  };

  return (
    <NodeLibraryContainer>
      <div className="wrap">
        <div className="title">
          {intl.formatMessage({ id: 'nodeLibrary.title' })}
        </div>
        <div className="node">
          {/* <div>
            <Input
              placeholder="搜索节点名称"
              onPressEnter={() => {
                // handlePageSearch();
              }}
              bordered={false}
              prefix={
                <SearchOutlined
                  onClick={() => {
                    // handlePageSearch();
                  }}
                />
              }
            />
          </div> */}
          <Tree
            loadData={onLoadData}
            treeData={treeData}
            fieldNames={{ title: 'name' }}
            onSelect={(keys, params) => {
              // 点击名称时也需要展开或收起
              const expanded = params?.node?.expanded;
              if (!expanded) {
                setExpandedKeys([...expandedKeys, params?.node?.key]);
              } else {
                setExpandedKeys((values) => {
                  return values?.filter((item) => item !== params?.node?.key);
                });
              }
            }}
            onExpand={(val: any) => {
              setExpandedKeys(val);
            }}
            expandedKeys={expandedKeys}
            titleRender={(nodeData) => {
              const {
                targetId,
                key,
                disabled: nodeDisabled,
                isChildrenNode,
              } = nodeData;
              if (targetId || isChildrenNode) {
                return (
                  <div className="tree-node">
                    <LeftNodeItem
                      itemData={nodeData}
                      key={key}
                      disabledItem={disabled}
                    />
                  </div>
                );
              }
              return (
                <div style={{ color: nodeDisabled ? 'rgba(0,0,0,0.25)' : '' }}>
                  {getIntlName(nodeData)}
                </div>
              );
            }}
          />
        </div>
      </div>
    </NodeLibraryContainer>
  );
};

export default NodeLibrary;
