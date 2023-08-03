import { history } from '@umijs/max';
import { Button, Result } from 'antd';

export default () => {
  return (
    <Result
      status="404"
      title="404"
      subTitle="抱歉，您访问的页面不存在"
      // subTitle="Sorry, the page you visited does not exist."
      extra={
        <Button
          type="primary"
          onClick={() => {
            history.push('/space/overview');
          }}
        >
          回到主页
        </Button>
      }
    />
  );
};
