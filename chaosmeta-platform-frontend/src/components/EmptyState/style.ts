import { styled } from '@umijs/max';

export const Container = styled.div`
  height: 100%;
  padding: 0ch;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #fff;
  border-radius: 6px;
  .desc {
    color: rgba(0, 0, 0, 0.65);
  }
  .title {
    font-size: 18px;
    font-weight: 500;
    color: rgba(0, 0, 0, 0.85);
    margin-bottom: 24px;
    margin-top: 8px;
  }
`;
