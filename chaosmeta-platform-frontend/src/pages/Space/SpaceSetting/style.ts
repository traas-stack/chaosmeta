import { styled } from '@umijs/max';

export const Container = styled.div`
  .ant-tabs-nav::before {
    border: none;
  }
`;

// 基本信息
export const BasicInfoContainer = styled.div`
  background-color: #fff;
  padding: 24px;
  border-radius: 6px;
  .ant-form {
    width: 50%;
    min-width: 600px;
  }
`;

// 攻击范围
export const AttackRangeContainer = styled.div`
  .search {
    background-color: #fff;
    padding-top: 24px;
    border-radius: 6px;
    margin-bottom: 16px;
  }
  .table {
    background-color: #fff;
    padding: 16px;
    border-radius: 6px;
    .operate {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 16px;
    }
  }
`;

// 成员列表-添加成员抽屉
export const AddUserDrawerContainer = styled.div`
  margin-top: -24px;
  margin-bottom: -24px;
  /* height: calc(100vh - 108px); */
  display: flex;
  .left {
    width: 50%;
    height: calc(100vh - 108px);
    overflow-y: auto;
    padding: 24px 24px 24px 0;
    border-right: 1px solid #e3e4e6;
    .check-all {
      height: 32px;
      margin-top: 12px;
      padding: 0 8px;
    }
    .check-item {
      padding: 0 8px;
      width: 100%;
      /* height: 100px; */
      height: 32px;
      line-height: 32px;
      border-radius: 6px;
    }
    .check-item-active {
      background-color: #e6f7ff;
    }
  }
  .right {
    width: 50%;
    padding: 24px;
    .title {
      margin-bottom: 12px;
    }
    .ant-tag {
      background-color: #eaebed;
      border: none;
      border-radius: 4px;
      margin-bottom: 4px;
    }
  }
`;

// 标签管理-添加标签抽屉
export const AddTagDrawerContainer = styled.div`
  .label {
    margin-bottom: 12px;
  }
  .tag {
    border-radius: 6px;
    border: 1px solid rgba(0, 10, 26, 0.16);
    padding: 4px 4px 0 4px;
    .ant-tag {
      font-size: 12px;
      color: rgba(0, 10, 26, 0.68);
      white-space: break-spaces;
      margin-bottom: 4px;
      span {
        color: rgba(0, 10, 26, 0.68);
      }
      span:hover {
        color: #000;
      }
    }
    .add {
      border-style: dashed;
      cursor: pointer;
    }
  }
`;

// 标签管理-添加标签抽屉-pop
export const AddTagPopContent = styled.div`
  width: 180px;
  .tip {
    margin-top: 4px;
    font-size: 12px;
    color: rgba(255, 77, 79, 1);
  }
  .ant-form-item {
    margin: 0;
  }
  .tags {
    margin: 8px 0 12px 0;
    span {
      width: 20px;
      height: 20px;
      display: flex;
      justify-content: center;
      align-items: center;
      margin: 0;
      color: #707070;
      cursor: pointer;
    }
  }
`;