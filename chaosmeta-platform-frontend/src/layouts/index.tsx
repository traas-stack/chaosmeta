import { ConfigProvider } from 'antd';
import { Outlet } from 'umi';
import { Container } from './style';

export default function Layout(props: any) {
  
  console.log(props, 'prosp')
  return (
    <ConfigProvider>
      <Container>
        <Outlet />
      </Container>
    </ConfigProvider>
  );
}
