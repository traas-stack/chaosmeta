import { styled } from '@umijs/max';

export const Container = styled.div`
  .result-list {
    .search {
      padding-bottom: 0;
      margin-bottom: 24px;
    }
    .table {
      background-color: #fff;
      padding: 16px;
      border-radius: 6px;
    }
    .ellipsis {
      position: relative;
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
      word-break: keep-all;
    }
  }
`;
