import ShowText from '@/components/ShowText';
import {
  CheckOutlined,
  CloseOutlined,
  DeleteOutlined,
  EditOutlined,
} from '@ant-design/icons';
import { Form, Input, InputNumber, Select, Space } from 'antd';
import React, { useEffect, useState } from 'react';
import { NodeConfigContainer } from '../style';

interface IProps {
  form: any;
  activeCol: any;
  arrangeList: any[];
  setArrangeList: any;
  setActiveCol: any;
}

/**
 * 节点信息配置
 * @param props
 * @returns
 */
const NodeConfig: React.FC<IProps> = (props) => {
  const { form, activeCol, setArrangeList, setActiveCol } = props;
  const [editTitleState, setEditTitleState] = useState<boolean>(false);

  /**
   * 持续时长变化时要同步修改编排节点的时长
   */
  const handleTimeChange = () => {
    const curTime = form.getFieldValue('time');
    setArrangeList((result: any) => {
      const values = JSON.parse(JSON.stringify(result));
      const parentIndex = values?.findIndex(
        (item: { id: any }) => item?.id === activeCol?.parentId,
      );
      if (parentIndex !== -1 && activeCol?.index >= 0) {
        values[parentIndex].children[activeCol.index].second = curTime;
      }
      return values;
    });
  };

  /**
   * 删除当前选中的节点
   */
  const handleDeleteNode = () => {
    setArrangeList((result: any) => {
      const values = JSON.parse(JSON.stringify(result));
      const parentIndex = values?.findIndex(
        (item: { id: any }) => item?.id === activeCol?.parentId,
      );
      if (parentIndex !== -1 && activeCol?.index >= 0) {
        values[parentIndex]?.children?.splice(activeCol?.index, 1);
      }
      return values;
    });
    setActiveCol({ state: false });
  };

  useEffect(() => {
    form.setFieldsValue({
      time: activeCol?.second,
      title: 'CPU燃烧',
    });
  }, [activeCol]);

  return (
    <NodeConfigContainer>
      <Form form={form} layout="vertical">
        <div className="header">
          {/* <Space>
            <span>CPU燃烧</span>
            <EditOutlined
              className="cancel"
              style={{ color: '#1890FF' }}
              onClick={() => {
                setEditTitleState(false);
              }}
            />
          </Space> */}
          <Space>
            {editTitleState ? (
              <Form.Item name={'title'} noStyle>
                <Input placeholder="请输入" style={{ width: '120px' }} />
              </Form.Item>
            ) : (
              <Form.Item name={'title'} noStyle>
                <ShowText value="CPU燃烧" />
              </Form.Item>
            )}
            {editTitleState ? (
              <Space>
                <CloseOutlined
                  className="cancel"
                  style={{ color: '#FF4D4F' }}
                  onClick={() => {
                    setEditTitleState(false);
                  }}
                />
                <CheckOutlined
                  style={{ color: '#1890FF' }}
                  className="confirm"
                  onClick={() => {
                    setEditTitleState(false);
                  }}
                />
              </Space>
            ) : (
              <EditOutlined
                className="edit"
                style={{ color: '#1890FF' }}
                onClick={() => {
                  setEditTitleState(true);
                }}
              />
            )}
          </Space>
          <Space size={24}>
            <DeleteOutlined
              style={{ color: '#FF4D4F' }}
              onClick={() => {
                handleDeleteNode();
              }}
            />
            <CloseOutlined
              onClick={() => {
                setActiveCol({ state: false });
              }}
            />
          </Space>
        </div>
        <div className="form">
          <Form.Item label="节点类型" name="nodeType">
            <Input placeholder="节点类型" disabled />
          </Form.Item>

          <Form.Item
            label="持续时长"
            name="time"
            rules={[{ required: true, message: '请输入持续时长' }]}
          >
            <InputNumber
              placeholder="请输入持续时长"
              addonAfter={'秒'}
              onChange={handleTimeChange}
            />
          </Form.Item>

          <Form.Item
            label="CPU使用率"
            name="CPUuse"
            rules={[{ required: true, message: '请输入CPU使用率' }]}
          >
            <InputNumber
              placeholder="请输入CPU使用率"
              min={0}
              max={100}
              formatter={(value) => `${value}%`}
              parser={(value: any) => value!.replace('%', '')}
            />
          </Form.Item>
          <Form.Item label="CPU满载核数" name="heshu">
            <Select placeholder="请选择"></Select>
          </Form.Item>
          <div className="range">攻击范围</div>
          <Form.Item
            label="Kubernetes Namespace"
            name="ku"
            rules={[{ required: true, message: '请输入' }]}
          >
            <Select placeholder="请选择"></Select>
          </Form.Item>
          <Form.Item label="Kubernetes Label" name="kulebel">
            <Input placeholder="请输入" />
          </Form.Item>
          <Form.Item label="应用" name="app">
            <Select placeholder="请选择"></Select>
          </Form.Item>
          <Form.Item label="筛选模式" name="mode">
            <Select placeholder="请选择"></Select>
          </Form.Item>
          <Form.Item name="mode">
            <Select placeholder="请选择"></Select>
          </Form.Item>
        </div>
      </Form>
    </NodeConfigContainer>
  );
};

export default NodeConfig;
