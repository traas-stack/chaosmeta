import { styled } from '@umijs/max';

export const Container = styled.div`
  height: 100%;
  display: flex;
  justify-content: space-evenly;
  align-items: center;
  background-color: #fff;
  border-radius: 6px;
  text-align: left;
  .left {
    max-width: 380px;
    .desc {
      color: rgba(0, 0, 0, 0.65);
    }
    .title,
    .top-title {
      font-size: 18px;
      font-weight: 500;
      color: rgba(0, 0, 0, 0.85);
    }
    .title {
      margin-bottom: 24px;
      margin-top: 8px;
    }
    .top-title {
      margin-bottom: 8px;
    }
  }
`;
