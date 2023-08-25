import { formatTime } from '@/utils/format';
import { EditOutlined } from '@ant-design/icons';
import React, { ReactNode } from 'react';
import { Container } from './styles';

interface IProps {
  value?: string | number | string[] | number[];
  /**
   * 是否为时间
   * @default false
   */
  isTime?: boolean;
  style?: React.CSSProperties;
  /**
   * 是否需要省略展示
   * @default false
   */
  ellipsis?: boolean;
  /**
   * 是否允许编辑
   * @default false
   */
  isEdit?: boolean;
}

const ShowText: React.FC<IProps> = (props) => {
  const { value, isTime, style, isEdit, ellipsis } = props;
  let renderText: ReactNode = <span>--</span>;
  if (value) {
    if (isTime) {
      renderText = <span>{formatTime(value)}</span>;
    } else if (isEdit) {
      renderText = <span>{value}</span>;
    } else {
      renderText = <span>{value}</span>;
    }
  }
  return (
    <Container>
      <span style={style} className={ellipsis ? 'ellipsis' : ''}>
        {renderText}
      </span>
      {isEdit && (
        <span >
          <EditOutlined />
        </span>
      )}
    </Container>
  );
};

export default React.memo(ShowText);
