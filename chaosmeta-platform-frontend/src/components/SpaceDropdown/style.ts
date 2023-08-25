import { styled } from '@umijs/max';

export const SpaceMenu = styled.div`
  width: 190px;
  border-radius: 6px;
  background-color: #fff;
  padding: 8px 0;
  div:first-child {
    /* width: 180px; */
    padding: 0 16px;
  }
  ul.ant-dropdown-menu {
    margin: 8px 0;
    padding: 0;
  }
  .ant-dropdown-menu {
    box-shadow: none;
  }
  .add-space {
    padding: 12px 12px 0 12px;
  }
  .more {
    padding: 0 12px;
    font-size: 12px;
    span:first-child {
      color: rgba(0, 10, 26, 0.47);
    }
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
  overflow: hidden;
  white-space: nowrap;
  cursor: pointer;
`;
