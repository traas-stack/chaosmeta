import React from 'react';
import { Container } from './style';

interface IProps {
  desc?: string;
  title?: string;
  imgSrc?: string;
  topTitle?: string;
  btns?: any;
}
const EmptyCustom: React.FC<IProps> = (props) => {
  const {
    desc,
    title,
    imgSrc = 'https://mdn.alipayobjects.com/huamei_d3kmvr/afts/img/A*ImtnQbACXJAAAAAAAAAAAAAADmKmAQ/original',
    topTitle,
    btns,
  } = props;
  return (
    <Container>
      <div className="left">
        <div className="top-title">{topTitle}</div>
        <div className="desc">{desc}</div>
        <div className="title">{title}</div>
        {btns}
      </div>
      <div>
        <img src={imgSrc} />
      </div>
    </Container>
  );
};

export default React.memo(EmptyCustom);
