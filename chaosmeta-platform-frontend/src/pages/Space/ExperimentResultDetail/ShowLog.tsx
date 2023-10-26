import { LogConainer } from './style';
// 编辑器相关, 顺序不能变更
import AceEditor from 'react-ace';
import 'ace-builds/src-noconflict/mode-java';
import 'ace-builds/src-noconflict/theme-monokai';

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
      <AceEditor
        style={{
          marginTop: '8px',
          borderRadius: '6px',
          maxHeight: '200px',
          backgroundColor: '#292e33',
        }}
        mode="java"
        theme="monokai"
        showGutter={true}
        fontSize={14}
        showPrintMargin={false}
        wrapEnabled
        value={markdown}
        readOnly
        width="100%"
        name="UNIQUE_ID_OF_DIV"
        editorProps={{ $blockScrolling: true }}
        setOptions={{
          enableBasicAutocompletion: true,
          enableLiveAutocompletion: true,
          highlightActiveLine: false,
        }}
      />
    </LogConainer>
  );
};

export default ShowLog;
