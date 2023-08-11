import { styled } from '@umijs/max';

export const Container = styled.div`
  /* font-size: 14px;
  color: #ff4d4f; */
  .ellipsis {
    position: relative;
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
    word-break: keep-all;
  }
  /* 头部样式 */
  .ant-page-header {
    height: 72px;
    padding: 0 24px;
    border-bottom: 1px solid #dadde7;
    box-shadow: inset 0 -1px 0 0 rgba(6, 17, 120, 0.1);
    /* title部分 */
    .ant-page-header-heading {
      .ant-page-header-heading-title {
        width: 300px;
        color: rgba(4, 24, 94, 0.85);
        font-size: 16px;
        font-weight: 500;
        .ellipsis {
          width: 260px;
        }
        .cancel {
          color: #ff4d4f;
        }
        .edit,
        .confirm {
          color: #1890ff;
        }
        .tags {
          font-weight: 400;
          span {
            color: rgba(0, 10, 26, 0.47);
          }
        }
      }
      /* 右侧展示部分 */
      .ant-page-header-heading-extra {
        flex: 1;
      }
      .ant-page-header-heading-extra > .ant-space,
      .ant-page-header-heading-extra > .ant-space > .ant-space-item {
        width: 100%;
      }
      .header-extra {
        flex: 1;
        display: flex;
        justify-content: space-between;
        .ant-form-item {
          margin: 0;
          color: rgba(0, 0, 0, 0.65);
        }
        .ant-form-item-label {
          label {
            color: rgba(0, 0, 0, 0.65);
          }
        }
      }
    }
  }
`;
