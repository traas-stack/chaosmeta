import { tagColors } from '@/constants';
import {
  querySpaceTagName,
  spaceAddTag,
} from '@/services/chaosmeta/SpaceController';
import { CheckOutlined, PlusOutlined } from '@ant-design/icons';
import { useIntl, useRequest } from '@umijs/max';
import {
  Button,
  Drawer,
  Form,
  Input,
  Popover,
  Space,
  Tag,
  message,
} from 'antd';
import React, { useState } from 'react';
import { AddTagDrawerContainer, AddTagPopContent } from '../style';
interface IProps {
  open: boolean;
  setOpen: (val: boolean) => void;
  spaceId: string;
  handlePageSearch: () => void;
}
const AddTagDrawer: React.FC<IProps> = (props) => {
  const { open, setOpen, spaceId, handlePageSearch } = props;
  const [form] = Form.useForm();
  // 添加标签的展示
  const [addTagOpen, setAddTagOpen] = useState<boolean>(false);
  // 是否有重复标签存在的提示
  const [isTip, setIsTip] = useState<boolean>(false);
  // 当前选中标签的类型
  const [checkedTagType, setCheckedTagType] = useState<string>('default');
  // 选择要添加的标签list
  const [addTagList, setAddTagList] = useState<any>([]);
  const handleClose = () => {
    setOpen(false);
  };
  const [submitLoading, setSubmitLoading] = useState<boolean>(false);
  const intl = useIntl();
  /**
   * 查询标签信息，用于校验标签是否已经添加过
   */
  const handleCheckTag = useRequest(querySpaceTagName, {
    manual: true,
    formatResult: (res) => res,
  });

  // 添加标签校验
  const handleTagCheck = () => {
    form.validateFields().then(async (values) => {
      const name = values?.tagName;
      // 当前标签名称是否已输入过
      if (addTagList.some((item: { name: string }) => item.name === name)) {
        setIsTip(true);
        return;
      }
      const res = await handleCheckTag?.run({ ns_id: spaceId, name });
      //  data中有返回值证明当前标签已经添加过
      if (res.data?.name) {
        setIsTip(true);
        return;
      }
      form.setFieldValue('tagName', undefined);
      const temp = {
        color: checkedTagType,
        name,
      };
      setAddTagList([...addTagList, temp]);
      setAddTagOpen(false);
    });
  };

  /**
   * 提交
   */
  const handleSubmit = async () => {
    if (!addTagList?.length) {
      message.info(intl.formatMessage({ id: 'tag.empty.tip' }));
      return;
    }
    setSubmitLoading(true);
    const queryList: any[] = [];
    // 多个时循环调用
    addTagList?.forEach(
      (item: { id: string | number; name: string; color: string }) => {
        const curQuery = spaceAddTag({ ...item, id: spaceId });
        queryList.push(curQuery);
      },
    );
    // 调用成功之后关闭弹窗
    const result = await Promise.all(queryList);
    if (result.every((item) => item.code === 200)) {
      message.success(intl.formatMessage({ id: 'tag.create.success.tip' }));
      setSubmitLoading(false);
      setOpen(false);
      handlePageSearch();
    }
  };

  const content = () => {
    return (
      <AddTagPopContent>
        <Form form={form}>
          <Form.Item
            name={'tagName'}
            rules={[
              {
                message: intl.formatMessage({ id: 'inputPlaceholder' }),
                required: true,
              },
            ]}
          >
            <Input
              placeholder={intl.formatMessage({ id: 'inputPlaceholder' })}
              onChange={() => {
                setIsTip(false);
              }}
            />
          </Form.Item>
        </Form>

        {isTip && (
          <div className="tip">
            {intl.formatMessage({ id: 'tag.repeat.text' })}
          </div>
        )}
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
        <div style={{ textAlign: 'right' }}>
          <Space>
            <Button
              size="small"
              onClick={() => {
                setAddTagOpen(false);
              }}
            >
              {intl.formatMessage({ id: 'cancel' })}
            </Button>
            <Button
              size="small"
              type="primary"
              onClick={() => {
                handleTagCheck();
              }}
              loading={handleCheckTag?.loading}
            >
              {intl.formatMessage({ id: 'confirm' })}
            </Button>
          </Space>
        </div>
      </AddTagPopContent>
    );
  };

  // 获取标签对应颜色
  const getTagColor = (type: string) => {
    const color = tagColors?.filter((item) => item.type === type)[0]?.color;
    return color;
  };

  return (
    <Drawer
      open={open}
      onClose={handleClose}
      title={intl.formatMessage({ id: 'tagManageMent.add.title' })}
      width={480}
      footer={
        <div style={{ textAlign: 'right' }}>
          <Space>
            <Button onClick={handleClose}>
              {intl.formatMessage({ id: 'cancel' })}
            </Button>
            <Button
              type="primary"
              onClick={handleSubmit}
              loading={submitLoading}
            >
              {intl.formatMessage({ id: 'confirm' })}
            </Button>
          </Space>
        </div>
      }
    >
      <AddTagDrawerContainer>
        <div className="label">{intl.formatMessage({ id: 'label' })}</div>
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
              <PlusOutlined /> {intl.formatMessage({ id: 'label' })}
            </Tag>
          </Popover>
        </div>
      </AddTagDrawerContainer>
    </Drawer>
  );
};
export default React.memo(AddTagDrawer);
