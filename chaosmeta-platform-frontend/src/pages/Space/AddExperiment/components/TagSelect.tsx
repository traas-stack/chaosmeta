import { tagColors } from '@/constants';
import {
  querySpaceTagList,
  querySpaceTagName,
  spaceAddTag,
} from '@/services/chaosmeta/SpaceController';
import { CheckOutlined, PlusOutlined } from '@ant-design/icons';
import { useRequest } from '@umijs/max';
import { Button, Form, Input, Popover, Space, Tag, message } from 'antd';
import { useState } from 'react';
import { AddTagDrawerContainer, AddTagPopContent } from '../style';

interface Props {
  spaceId: string;
  addTagList: any[];
  setAddTagList: any;
}

const TagSelect = (props: Props) => {
  const { spaceId, addTagList, setAddTagList } = props;
  const [form] = Form.useForm();
  // 添加标签的展示
  const [addTagOpen, setAddTagOpen] = useState<boolean>(false);

  // 是否有重复标签存在的提示
  const [isTip, setIsTip] = useState<boolean>(false);
  // 当前选中标签的类型
  const [checkedTagType, setCheckedTagType] = useState<string>('default');
  const [searchTagList, setSearchTagList] = useState([]);
  // 获取标签对应颜色
  const getTagColor = (type: string) => {
    const color = tagColors?.filter((item) => item.type === type)[0]?.color;
    return color;
  };

  /**
   * 查询标签信息，用于校验标签是否已经添加过
   */
  const handleCheckTag = useRequest(querySpaceTagName, {
    manual: true,
    formatResult: (res) => res,
  });

  /**
   * 新建标签
   */
  const handleAddTag = useRequest(spaceAddTag, {
    manual: true,
    formatResult: (res) => res,
  });

  /**
   * 分页接口
   */
  const getTagList = useRequest(querySpaceTagList, {
    manual: true,
    debounceInterval: 300,
    formatResult: (res) => res,
    onSuccess: (res) => {
      if (res?.code === 200) {
        setSearchTagList(res.data?.labels || []);
      }
    },
  });

  /**
   * 更新addTaglist
   * @param temp
   * @returns
   */
  const hanldeUpdateTagList = (temp: any) => {
    if (addTagList?.some((item: { id: any }) => item.id === temp.id)) {
      message.info('标签已存在');
      return;
    }
    setAddTagList((origin: any) => {
      return [...origin, temp];
    });
  };

  // 添加标签校验
  const handleTagCheck = () => {
    form.validateFields().then(async (values: { tagName: any }) => {
      const name = values?.tagName;
      // 当前标签名称是否已输入过
      if (addTagList?.some((item: { name: string }) => item.name === name)) {
        setIsTip(true);
        return;
      }
      const res = await handleCheckTag?.run({ ns_id: spaceId, name });
      //  data中有返回值证明当前标签已经添加过
      // 将当前选择的标签放入到addlist中备选
      if (res.data?.name) {
        hanldeUpdateTagList(res?.data);
      } else {
        // 当检索的标签不存在时，新建标签
        const addRes: any = await handleAddTag?.run({
          color: checkedTagType,
          name: name,
          id: spaceId,
        });
        if (addRes?.data?.id) {
          const temp = { id: addRes?.data?.id, name, color: checkedTagType };
          hanldeUpdateTagList(temp);
        }
      }
      setAddTagOpen(false);
    });
  };

  const content = () => {
    return (
      <AddTagPopContent>
        <Form form={form}>
          <Form.Item
            name={'tagName'}
            rules={[{ message: '请输入', required: true }]}
          >
            <Input
              placeholder="请输入"
              onChange={(value) => {
                const name = value.target?.value;
                // 获取标签列表
                if (name) {
                  getTagList?.run({
                    id: Number(spaceId),
                    name,
                    page: 1,
                    pageSize: 10,
                  });
                } else {
                  setTimeout(() => {
                    setSearchTagList([]);
                  }, 600);
                }
                setIsTip(false);
              }}
            />
          </Form.Item>
        </Form>

        {isTip && <div className="tip">标签已经存在，请重新输入</div>}
        <Space size={12} className="tags">
          {tagColors?.map((item) => {
            return (
              <Tag
                key={item.type}
                color={item.color}
                onClick={() => {
                  setCheckedTagType(item.type);
                }}
              >
                {checkedTagType === item.type && (
                  <CheckOutlined color="#707070" />
                )}
              </Tag>
            );
          })}
        </Space>
        {searchTagList?.map((el: any) => {
          const temp = tagColors?.filter((item) => item.type === el?.color)[0];
          return (
            <Tag
              color={temp?.type}
              key={el?.id}
              style={{ marginBottom: '4px' }}
              onClick={() => {
                form.setFieldValue('tagName', el.name);
              }}
            >
              {el?.name}
            </Tag>
          );
        })}
        <div style={{ textAlign: 'right' }}>
          <Space>
            <Button
              size="small"
              onClick={() => {
                setAddTagOpen(false);
              }}
            >
              取消
            </Button>
            <Button
              size="small"
              type="primary"
              onClick={() => {
                handleTagCheck();
              }}
              loading={handleCheckTag?.loading}
            >
              确定
            </Button>
          </Space>
        </div>
      </AddTagPopContent>
    );
  };

  return (
    <>
      <AddTagDrawerContainer>
        <div className="label">标签</div>
        <div className="tag">
          {addTagList?.map((item: { name: string; color: string }) => {
            return (
              <Tag
                key={item.name}
                color={getTagColor(item.color)}
                closeIcon
                onClose={() => {
                  setAddTagList(() => {
                    const newList = addTagList?.filter(
                      (el: { name: string }) => el.name !== item.name,
                    );
                    return newList;
                  });
                }}
              >
                {item.name}
              </Tag>
            );
          })}
          <Popover
            content={content()}
            trigger="click"
            placement="bottomLeft"
            open={addTagOpen}
          >
            <Tag
              className="add"
              onClick={() => {
                setAddTagOpen(true);
              }}
            >
              <PlusOutlined /> 标签
            </Tag>
          </Popover>
        </div>
      </AddTagDrawerContainer>
    </>
  );
};

export default TagSelect;
