export interface Column {
  title?: string | React.ReactNode;
  width?: number;
  dataIndex?: string;
  key?: string;
  fixed?: string;
  ellipsis?: boolean;
  defaultSortOrder?: string;
  sorter?: any;
  align?: string;
  summary?: boolean; // 是否是统计列
  sortOrder?: string | boolean;
  render?: (text?: any, record?: any, index?: number) => void;
}
