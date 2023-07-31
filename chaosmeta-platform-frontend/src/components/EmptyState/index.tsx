import { Button } from 'antd';
import React from 'react';
import { Container } from './style';

interface IProps {
  desc: string;
  title: string;
  btnText: string;
  btnOperate?: () => void;
}
const EmptyState: React.FC<IProps> = (props) => {
  const { desc, title, btnOperate = () => {}, btnText } = props;
  return (
    <Container>
      <div>
        <div className="desc">{desc}</div>
        <div className="title">{title}</div>
        <div>
          <Button type="primary" onClick={btnOperate}>
            {btnText}
          </Button>
        </div>
      </div>
      <div>
        <img src="https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*ImtnQbACXJAAAAAAAAAAAAAADmKmAQ/original" />
      </div>
    </Container>
  );
};

export default React.memo(EmptyState);
