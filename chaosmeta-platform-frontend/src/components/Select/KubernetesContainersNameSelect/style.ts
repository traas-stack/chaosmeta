import { styled ,css} from '@umijs/max';
import { Select } from 'antd';


export const GrayText = styled.div`
  color: rgba(0, 0, 0, 0.88);
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  word-break: keep-all;
  width: 35%;
`;

export const OptionRow = styled.div`
  white-space: nowrap;
  text-overflow: ellipsis;
  word-break: keep-all;
  color: #999;
  display: flex;
  align-items: center;
  justify-content: space-between;
  .pods {
    white-space: nowrap;
    text-overflow: ellipsis;
    word-break: keep-all;
    width: 60%;
    overflow: hidden;
  }
`;


export const GroupLabel = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 14px;
  color: rgba(0, 0, 0, 0.88);
  font-weight: 600;
`;
