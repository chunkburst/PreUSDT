import { ref } from "vue";

export interface WalletDetail {
  id: number;
  group_id?: string;
  name: string;
  address: string;
  trade_type: string;
  trade_types: string[];
  remark: string;
  other_notify: number;
  status: number;
  kind: string;
  base: string;
  created_at?: string;
  updated_at?: string;
}

export const useWalletDetail = () => {
  const detailVisible = ref(false);
  const detailData = ref<WalletDetail>({
    id: 0,
    name: "",
    address: "",
    trade_type: "",
    trade_types: [],
    remark: "",
    other_notify: 0,
    status: 0,
    kind: "single",
    base: "generic"
  });

  const showDetail = (record: WalletDetail) => {
    detailData.value = { ...record, trade_types: record.trade_types || [] };
    detailVisible.value = true;
  };

  const closeDetail = () => {
    detailVisible.value = false;
    detailData.value = {
      id: 0,
      name: "",
      address: "",
      trade_type: "",
      trade_types: [],
      remark: "",
      other_notify: 0,
      status: 0,
      kind: "single",
      base: "generic"
    };
  };

  return {
    detailVisible,
    detailData,
    showDetail,
    closeDetail
  };
};
