import { styled } from '@umijs/max';

// 半透明背景区域
export const Area = styled.div`
  background-color: rgba(255, 255, 255, 0.5);
  padding: 24px;
  border-radius: 6px;
  .area-operate {
    display: flex;
    justify-content: space-between;
    margin-bottom: 16px;
    .ant-input-affix-wrapper {
      min-width: 230px;
    }
    .ant-input-suffix {
      color: rgba(0, 0, 0, 0.4);
      cursor: pointer;
    }
    .title {
      font-size: 16px;
      font-weight: 500;
      color: rgba(0, 0, 0, 0.85);
    }
  }
  .area-content {
    background-color: #fff;
    padding: 16px;
    border-radius: 6px;
  }
`;
