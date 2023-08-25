import { Spin } from 'antd';

const PageLoading: React.FC = () => {
  return (
    <div style={{ marginTop: '200px', textAlign: 'center' }}>
      <Spin size="large" />
    </div>
  );
};

export default PageLoading;
