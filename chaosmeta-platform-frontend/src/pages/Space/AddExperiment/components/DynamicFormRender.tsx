import DynamicForm from '@/components/DynamicForm';
import { useIntl } from '@umijs/max';
import { useEffect, useState } from 'react';

interface IProps {
  fieldList: any[];
  form?: any;
  nodeType: string;
  readonly?: boolean;
}
/**
 * 动态表单模块区分渲染
 * @returns
 */
const DynamicFormRender: React.FC<IProps> = (props) => {
  const { fieldList, form, nodeType, readonly } = props;
  const intl = useIntl();
  // 度量通用参数表单
  const [measureCommonFields, setMeasureCommonFields] = useState<any[]>([]);
  // 度量判定参数表单
  const [measureJudgeFields, setMeasureJudgeFields] = useState<any[]>([]);
  // 度量执行参数表单
  const [measureExecuteFields, setMeasureExecuteFields] = useState<any[]>([]);
  // 故障参数表单
  const [faultFields, setFaultFields] = useState<any[]>([]);
  // 流量通用参数表单
  const [flowCommonFields, setFlowCommonFields] = useState<any[]>([]);
  // 流量执行参数表单
  const [flowExecuteFields, setflowExecuteFields] = useState<any[]>([]);

  // 流量通用参数keys
  const flowcommonParamsKeys = [
    // 'flowType',
    'duration',
    'parallelism',
    'source',
  ];

  // 度量通用参数keys
  const measureCommonParamsKeys = [
    // 'measureType',
    'duration',
    'interval',
    'successCount',
    'failedCount',
  ];
  // 度量判定参数keys
  const measureJudgeParamsKeys = ['judgeType', 'judgeValue'];

  useEffect(() => {
    if (fieldList?.length > 0) {
      if (nodeType === 'fault') {
        setFaultFields(fieldList);
      }
      if (nodeType === 'measure') {
        // 度量通用参数
        setMeasureCommonFields(() => {
          return fieldList?.filter((item) =>
            measureCommonParamsKeys?.includes(item?.key),
          );
        });
        // 度量判定参数
        setMeasureJudgeFields(() => {
          return fieldList?.filter((item) =>
            measureJudgeParamsKeys?.includes(item?.key),
          );
        });
        // 度量执行参数 - 不属于通用和判定的即为执行参数
        setMeasureExecuteFields(() => {
          return fieldList?.filter(
            (item) =>
              !measureJudgeParamsKeys?.includes(item?.key) &&
              !measureCommonParamsKeys?.includes(item?.key) &&
              item?.key !== 'measureType',
          );
        });
      }
      if (nodeType === 'flow') {
        // 流量通用参数
        setFlowCommonFields(() => {
          return fieldList?.filter((item) =>
            flowcommonParamsKeys?.includes(item?.key),
          );
        });
        // 流量执行参数 - 不属于通用的即为执行参数
        setflowExecuteFields(() => {
          return fieldList?.filter(
            (item) =>
              !flowcommonParamsKeys?.includes(item?.key) &&
              item?.key !== 'flowType',
          );
        });
      }
    }
  }, [fieldList, nodeType]);

  const fieldsFormRender = (fields: any[], title: string) => {
    return (
      fields?.length > 0 && (
        <>
          {readonly ? (
            <div className="subtitle range">
              {intl.formatMessage({ id: title })}
            </div>
          ) : (
            <div className="range">{intl.formatMessage({ id: title })}</div>
          )}

          <DynamicForm
            fieldList={fields}
            parentName={'args_value'}
            form={form}
            readonly={readonly}
          />
        </>
      )
    );
  };

  return (
    <>
      {/* 故障参数 */}
      {nodeType === 'fault' && (
        <>{fieldsFormRender(faultFields, 'faultParameters')}</>
      )}
      {nodeType === 'measure' && (
        <>
          <DynamicForm
            fieldList={fieldList?.filter((item) => item?.key === 'measureType')}
            parentName={'args_value'}
            form={form}
            readonly={readonly}
          />
          {/* 度量通用参数表单 */}
          {fieldsFormRender(measureCommonFields, 'commonParameters')}
          {/* 度量判定参数表单 */}
          {fieldsFormRender(measureJudgeFields, 'judgmentParameters')}
          {/* 度量执行参数表单 */}
          {fieldsFormRender(measureExecuteFields, 'executionParameters')}
        </>
      )}
      {nodeType === 'flow' && (
        <>
          <DynamicForm
            fieldList={fieldList?.filter((item) => item?.key === 'flowType')}
            parentName={'args_value'}
            form={form}
            readonly={readonly}
          />
          {/* 流量通用参数表单 */}
          {fieldsFormRender(flowCommonFields, 'commonParameters')}
          {/* 流量执行参数表单 */}
          {fieldsFormRender(flowExecuteFields, 'executionParameters')}
        </>
      )}
    </>
  );
};
export default DynamicFormRender;
