import { styled } from '@umijs/max';

export const SpaceMenu = styled.div`
  border-radius: 6px;
  background-color: #fff;
  padding: 8px 0;
  div:first-child {
    width: 180px;
    padding: 0 16px;
  }
  ul.ant-dropdown-menu {
    margin: 8px 0;
    padding: 0;
    .ant-dropdown-menu-item-selected {
      /* background-color: #edeeef; */
      /* border-radius: 0; */
      /* color: rgba(0, 0, 0, 0.88); */
    }
  }
  .ant-dropdown-menu {
    box-shadow: none;
  }
`;

export const SpaceContent = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 6px;
  padding: 0 16px;
  height: 52px;
  color: #000a1a;
  font-weight: 500;
  border-bottom: 1px solid rgba(0, 10, 26, 0.07);
`;
