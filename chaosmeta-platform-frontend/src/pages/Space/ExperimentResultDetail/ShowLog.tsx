import ReactMarkdown from 'react-markdown';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { atomOneDark } from 'react-syntax-highlighter/dist/esm/styles/hljs';
import { LogConainer } from './style';

interface Props {
  message?: string;
}
const ShowLog = (props: Props) => {
  const { message: markdown } = props;

  return (
    <LogConainer>
      {/* todo -- 后端暂不支持机器相关信息 */}
      {/* <Space style={{ width: '100%', justifyContent: 'space-between' }}>
        <div>
          总共注入15台机器，成功11台，失败4台，失败原因：agent版本过低。
        </div>
        <Select options={options} style={{ width: '200px' }}></Select>
      </Space> */}
      <div className="log-contet">
        <ReactMarkdown
          components={{
            code({ children, ...props }) {
              return (
                <SyntaxHighlighter
                  {...props}
                  customStyle={{ maxHeight: 400, borderRadius: 8, margin: 0 }}
                  style={atomOneDark}
                  language={'jsx'}
                  PreTag="div"
                  showLineNumbers
                >
                  {children as string}
                </SyntaxHighlighter>
              );
            },
          }}
        >
          {`
          ${markdown}
          `}
        </ReactMarkdown>
      </div>
    </LogConainer>
  );
};

export default ShowLog;
