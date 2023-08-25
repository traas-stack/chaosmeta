import { styled } from '@umijs/max';

interface OperateType {
  operatetype: 'login' | 'register';
}

export const Container = styled.div`
  height: 100vh;
  background-image: url('https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*mLnfQ6O7lVgAAAAAAAAAAAAADmKmAQ/original');
  background-size: cover;
  font-size: 14px;
  /* .login {
    width: 200px;
    margin: auto;
    padding-top: 100px;
  }
  .operate {
    width: 100%;
    justify-content: flex-end;
  } */
  .seize {
    height: 1px;
  }
  .card {
    width: 420px;
    /* height: 500px; */
    /* position: relative; */
    border-radius: 12px;
    margin-left: 16%;
    margin-top: 7%;
    padding: 0 60px 38px 60px;
    backdrop-filter: blur(8px);
    background-color: rgba(255, 255, 255, 0.7);
    box-shadow: 0 6px 16px -8px rgba(0, 10, 26, 0.13),
      0 9px 28px 0 rgba(0, 10, 26, 0.1), 0 12px 48px 16px rgba(0, 10, 26, 0.08);
    border-radius: 12px;

    .content {
      img {
        padding: 48px 0;
      }
      .title {
        font-weight: 500;
        font-size: 28px;
      }
      .tip {
        color: rgba(0, 10, 26, 0.47);
        margin-bottom: 16px;
      }
      .ant-form-item-explain {
        margin-bottom: 24px;
      }
      button,
      input {
        height: 40px;
        width: 100%;
        font-size: 16px;
      }
      button {
        margin-bottom: 16px;
        margin-top: 8px;
      }
    }
  }
`;

// export const OperateArea = styled.div<OperateType>`
//   width: 100%;
//   display: flex;
//   /* justify-content: space-between; */
//   align-items: center;
//   justify-content: ${(props) => {
//     const { operatetype } = props;
//     return operatetype === 'login' ? 'flex-end' : 'space-between';
//   }};
//   .ant-btn:last-child {
//     margin-left: 8px;
//   }
// `;
