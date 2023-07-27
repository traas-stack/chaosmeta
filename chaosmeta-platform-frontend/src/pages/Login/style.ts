import { styled } from '@umijs/max';

interface OperateType {
  operatetype: 'login' | 'register';
}

export const Container = styled.div`
  background-image: url('https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*wt5SSrooBdIAAAAAAAAAAAAADmKmAQ/original') !important;
  background-size: cover;
  height: 100vh;
  font-size: 14px;
  /* display: flex;
  justify-content: center;
  align-items: center; */
  .login {
    width: 200px;
    margin: auto;
    padding-top: 100px;
  }
  .operate {
    width: 100%;
    justify-content: flex-end;
  }
`;

export const OperateArea = styled.div<OperateType>`
  width: 100%;
  display: flex;
  /* justify-content: space-between; */
  align-items: center;
  justify-content: ${(props) => {
    const { operatetype } = props;
    return operatetype === 'login' ? 'flex-end' : 'space-between';
  }};
  .ant-btn:last-child {
    margin-left: 8px;
  }
`;
