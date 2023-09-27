/**
 * 添加成员弹窗
 */

import { spaceAddUser } from '@/services/chaosmeta/SpaceController';
import { getSpaceUserList } from '@/services/chaosmeta/UserController';
import { SearchOutlined } from '@ant-design/icons';
import { history, useIntl, useRequest } from '@umijs/max';
import {
  Button,
  Checkbox,
  Drawer,
  Empty,
  Form,
  Input,
  Radio,
  Space,
  Spin,
  Tag,
  message,
} from 'antd';
import { CheckboxChangeEvent } from 'antd/es/checkbox';
import { CheckboxValueType } from 'antd/es/checkbox/Group';
import React, { useEffect, useState } from 'react';
import { AddUserDrawerContainer } from '../style';
interface IProps {
  open: boolean;
  setOpen: (val: boolean) => void;
  handlePageSearch: () => void;
  spaceId: number | string;
}

const AddMemberModal: React.FC<IProps> = (props) => {
  const { open, setOpen, handlePageSearch, spaceId } = props;
  const [form] = Form.useForm();
  // 半选状态
  const [indeterminate, setIndeterminate] = useState(false);
  // 全选状态
  const [checkAll, setCheckAll] = useState(false);
  // 已选择的
  const [checkedList, setCheckedList] = useState<CheckboxValueType[]>([]);
  // 右侧已选择用户
  const [rightCheckedList, setRightCheckedList] = useState<any[]>([]);
  // 未加入的成员list
  const [noJoinUserList, setNoJoinUserList] = useState<any>([]);
  // 总的成员list
  const [userList, setUserList] = useState<any[]>([]);
  // 分页信息
  const [pageData, setPageData] = useState<any>({});
  const intl = useIntl();

  /**
   * 检索用户
   */
  const queryUserList = useRequest(getSpaceUserList, {
    manual: true,
    debounceInterval: 500,
    formatResult: (res) => res,
    onSuccess: (res) => {
      if (res?.code === 200) {
        setPageData(res?.data);
        const joinList = res?.data?.users?.filter(
          (item: { isJoin: boolean }) => !item.isJoin,
        );
        // 为1时重新赋值userlist，不为1代表滚动加载需拼接
        if (res.data?.page === 1) {
          setUserList(res?.data?.users);
          setNoJoinUserList(joinList);
        } else {
          setNoJoinUserList([...noJoinUserList, ...joinList]);
          setUserList([...userList, ...res?.data?.users]);
          if (checkedList?.length > 0) {
            setCheckAll(false);
            setIndeterminate(true);
          }
        }
      }
    },
  });

  /**
   * 添加成员接口
   */
  const addUser = useRequest(spaceAddUser, {
    manual: true,
    formatResult: (res) => res,
    onSuccess: (res) => {
      if (res?.code === 200) {
        message.success(intl.formatMessage({ id: 'addMember.add.success' }));
        setOpen(false);
        handlePageSearch();
      }
    },
  });

  /**
   * 添加成员
   */
  const handleAddUser = () => {
    form.validateFields().then((values) => {
      const users = rightCheckedList?.map((item) => {
        return {
          id: item.id,
          permission: values.permission,
        };
      });
      const params = {
        id: history.location.query.spaceId as string,
        users,
      };
      addUser.run(params);
    });
  };

  // 全选变化
  const onCheckAllChange = (e: CheckboxChangeEvent) => {
    const idList = noJoinUserList?.map((item: { id: number }) => item.id);
    setCheckedList(e.target.checked ? idList : []);
    setIndeterminate(false);
    setCheckAll(e.target.checked);
    if (e.target.checked) {
      // 全选时过滤出右侧未添加的项，添加到右侧
      const rightIdList = rightCheckedList?.map((item) => item.id);
      const temp = noJoinUserList.filter(
        (item: { id: number; isJoin: boolean }) => {
          return !rightIdList.includes(item.id) && !item.isJoin;
        },
      );
      setRightCheckedList([...rightCheckedList, ...temp]);
    }
  };

  /**
   * 检索用户
   */
  const handleSearchUser = () => {
    const name = form.getFieldValue('name');
    queryUserList.run({ page: 1, page_size: 30, name, id: spaceId });
  };

  useEffect(() => {
    if (open) {
      // 默认检索
      queryUserList.run({ page: 1, page_size: 30, id: spaceId });
    }
  }, [open]);

  /**
   * 滚动加载
   * @param event
   */
  const handleScroll = (event: any) => {
    if (event?.target) {
      const { scrollTop, clientHeight, scrollHeight } = event.target;
      if (
        scrollTop + clientHeight >= scrollHeight - 10 &&
        pageData?.users?.length >= pageData?.pageSize
      ) {
        queryUserList.run({
          page: pageData.page + 1,
          page_size: 20,
          id: spaceId,
        });
      }
    }
  };

  /**
   * 多选框变化时
   * @param target
   * @param id
   */
  const handleCheckChange = (target: any, id: number) => {
    // 将当前选择的人员保存到checkedList， 并添加到右侧已选择
    const { checked } = target;
    if (checked) {
      const curCheckList = [...checkedList, id];
      setCheckedList(curCheckList);
      // 长度相等全选
      if (curCheckList?.length === noJoinUserList?.length) {
        setIndeterminate(false);
        setCheckAll(true);
      } else {
        // 半选
        setIndeterminate(true);
      }
      // 选中左侧且右侧不包含该项，添加到右侧
      if (!rightCheckedList.some((item) => item.id === id)) {
        const curList = noJoinUserList?.filter(
          (item: { id: number }) => item.id === id,
        );
        setRightCheckedList([...rightCheckedList, ...curList]);
      }
      // 取消左侧选中，全选取消，根据check长度是否为0判断是否半选
    } else {
      setCheckAll(false);
      setCheckedList(() => {
        const newList = checkedList?.filter((el) => el !== id);
        if (newList?.length === 0) {
          setIndeterminate(false);
        } else {
          setIndeterminate(true);
        }
        return newList;
      });
    }
  };

  /**
   * 右侧删除已选项
   * @param id
   */
  const handleRightDelete = (id: number) => {
    // 右侧删除时左侧需取消选择
    if (checkedList.includes(id)) {
      setCheckedList(() => {
        const newList = checkedList?.filter((el) => el !== id);
        // 选择的list长度存在且小于全部数据长度时，属于半选状态，否则为全选或全部未选
        if (newList?.length > 0 && newList?.length < noJoinUserList?.length) {
          setIndeterminate(true);
        } else {
          setIndeterminate(false);
        }
        // 选择的list长度和总数据长度不同时，处于未全选状态
        if (newList?.length !== noJoinUserList?.length) {
          setCheckAll(false);
        }
        return newList;
      });
    }
    // 更新右侧数据
    setRightCheckedList((val) => {
      const newList = val.filter((el) => el?.id !== id);
      return newList;
    });
  };

  return (
    <Drawer
      title={intl.formatMessage({ id: 'addMember' })}
      open={open}
      width={800}
      onClose={() => {
        setOpen(false);
      }}
      footer={
        <div style={{ textAlign: 'right' }}>
          <Space>
            <Button>{intl.formatMessage({ id: 'cancel' })}</Button>
            <Button
              type="primary"
              onClick={handleAddUser}
              loading={addUser?.loading}
            >
              {intl.formatMessage({ id: 'confirm' })}
            </Button>
          </Space>
        </div>
      }
    >
      <Form form={form} layout="vertical">
        <AddUserDrawerContainer>
          <div className="left" onScroll={handleScroll}>
            <Form.Item
              label={intl.formatMessage({ id: 'username' })}
              name={'name'}
            >
              <Input
                placeholder={intl.formatMessage({ id: 'keyword' })}
                onChange={() => {
                  handleSearchUser();
                }}
                onPressEnter={handleSearchUser}
                suffix={<SearchOutlined onClick={handleSearchUser} />}
              />
            </Form.Item>
            {userList?.length > 0 ? (
              <Spin spinning={queryUserList?.loading}>
                <div>
                  {form.getFieldValue('name')
                    ? intl.formatMessage({ id: 'search' })
                    : intl.formatMessage({ id: 'all' })}{' '}
                  {intl.formatMessage({ id: 'memberList' })}
                </div>
                <div className="check-all">
                  <Checkbox
                    indeterminate={indeterminate}
                    onChange={onCheckAllChange}
                    checked={checkAll}
                  >
                    {intl.formatMessage({ id: 'selectAll' })}
                  </Checkbox>
                </div>
                {/* <Checkbox.Group value={checkedList} onChange={onChange}> */}
                {userList?.map((item) => {
                  return (
                    <div
                      className={
                        checkedList?.includes(item.id) && !item.isJoin
                          ? 'check-item-active check-item'
                          : 'check-item'
                      }
                      key={item?.id}
                    >
                      <Checkbox
                        disabled={item.isJoin}
                        checked={item.isJoin || checkedList.includes(item.id)}
                        onChange={(val) => {
                          handleCheckChange(val.target, item.id);
                        }}
                        value={item?.id}
                      >
                        {item.name}
                      </Checkbox>
                    </div>
                  );
                })}
                {/* </Checkbox.Group> */}
              </Spin>
            ) : (
              <Spin spinning={queryUserList?.loading}>
                <Empty
                  image={Empty.PRESENTED_IMAGE_SIMPLE}
                  description={intl.formatMessage({
                    id: 'addMember.search.result',
                  })}
                />
              </Spin>
            )}
            <div
              style={{
                textAlign: 'center',
                color: 'rgba(0,0,0,0.65)',
              }}
            >
              {queryUserList.loading && (
                <Spin>{intl.formatMessage({ id: 'addMember.loading' })}</Spin>
              )}
              {pageData?.users?.length < pageData.pageSize && (
                <div style={{ marginTop: '24px' }}>
                  {intl.formatMessage({ id: 'addMember.noMore' })}
                </div>
              )}
            </div>
          </div>
          <div className="right">
            <Form.Item
              label={intl.formatMessage({ id: 'addMember.user.permission' })}
              name={'permission'}
              initialValue={0}
              rules={[
                {
                  required: true,
                  message: intl.formatMessage({ id: 'selectPlaceholder' }),
                },
              ]}
            >
              <Radio.Group>
                <Radio value={0}>
                  {intl.formatMessage({ id: 'readonly' })}
                </Radio>
                <Radio value={1}>{intl.formatMessage({ id: 'write' })}</Radio>
              </Radio.Group>
            </Form.Item>
            <div className="title">
              {intl.formatMessage({ id: 'addMember.selected.user' })}
            </div>
            {rightCheckedList?.length > 0 && (
              <div>
                {rightCheckedList?.map((item) => {
                  return (
                    <Tag
                      closeIcon
                      onClose={() => {
                        handleRightDelete(item.id);
                      }}
                      key={item.id}
                    >
                      {item.name}
                    </Tag>
                  );
                })}
              </div>
            )}
          </div>
        </AddUserDrawerContainer>
      </Form>
    </Drawer>
  );
};

export default React.memo(AddMemberModal);
