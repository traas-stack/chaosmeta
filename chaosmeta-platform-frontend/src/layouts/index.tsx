import { ConfigProvider } from 'antd';
import { Outlet } from 'umi';
import { Container } from './style';

export default function Layout() {
  return (
    <ConfigProvider>
      <Container>
        <Outlet />
      </Container>
    </ConfigProvider>
  );
}
