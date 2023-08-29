import { styled } from '@umijs/max';

export const Container = styled.div`
  margin-bottom: 24px;
  .content {
    background-color: #fff;
    padding: 16px 24px;
    border-radius: 8px;
    .content-title {
      div {
        font-weight: 500;
        font-size: 16px;
        color: #1c2533;
        margin-right: 24px;
      }
      display: flex;
      margin-bottom: 10px;
      .ant-progress {
        width: 150px;
      }
    }
    .log {
      background-color: #f7f7fa;
      margin-top: 16px;
      border-radius: 6px;
      padding: 16px 24px;
      .ant-tabs-nav::before {
        border: none;
      }
    }
  }
`;

export const ObservationChartContainer = styled.div`
  height: 300px;
  width: 100%;
  background-color: #fff;
  border-radius: 6px;
  padding: 16px 24px;
`;

export const LogConainer = styled.div`
  .log-contet {
    margin-top: 16px;
  }
`;
