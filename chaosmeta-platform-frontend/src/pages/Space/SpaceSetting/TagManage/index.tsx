/**
 * 成员管理tab
 */
import { Area } from '@/components/CommonStyle';
import { SearchOutlined } from '@ant-design/icons';
import { Button, Input, Space, Tag } from 'antd';
import React from 'react';

interface IProps {}

const TagManage: React.FC<IProps> = () => {
  return (
    <Area>
      <div className="area-operate">
        <div className="title">标签列表</div>
        <Space>
          <Input
            placeholder="请输入标签名称"
            suffix={
              <SearchOutlined
                onClick={() => {
                  console.log('-==');
                }}
              />
            }
          />
          <Button>批量删除</Button>
          <Button type="primary">新建标签</Button>
        </Space>
      </div>
      <div className="tab-content">
        <div className='tag-content'>
          <Tag>测试标签</Tag>
        </div>
      </div>
    </Area>
  );
};

export default React.memo(TagManage);
