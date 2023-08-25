import { styled } from '@umijs/max';

export const Container = styled.span`
  position: relative;
  display: flex;
  .ellipsis {
    position: relative;
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
    word-break: keep-all;
  }
`;
