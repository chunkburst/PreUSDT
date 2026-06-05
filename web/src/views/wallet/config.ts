interface WalletBaseOption {
  value: string;
  label: string;
  trade_types: string[];
}

interface List {
  id: number;
  group_id?: string;
  name: string;
  address: string;
  status: number;
  createTime?: string;
  trade_type?: string;
  trade_types?: string[];
  remark?: string;
  other_notify?: number;
  kind?: string;
  base?: string;
  created_at?: string;
  updated_at?: string;
}

interface FormData {
  form: {
    name: string;
    address: string;
    trade_type: string;
  };
  search: boolean;
}

interface AddForm {
  name: string;
  address: string;
  trade_type: string;
  trade_types: string[];
  kind: string;
  base: string;
  remark: string;
  other_notify: number;
}

interface ModForm {
  id: number;
  name: string;
  status: number;
  address: string;
  trade_type: string;
  trade_types: string[];
  kind: string;
  base: string;
  remark: string;
  other_notify: number;
}

interface Pagination {
  showPageSize: boolean;
  showTotal: boolean;
  current: number;
  pageSize: number;
  total: number;
}

export type { List, FormData, Pagination, AddForm, ModForm, WalletBaseOption };
