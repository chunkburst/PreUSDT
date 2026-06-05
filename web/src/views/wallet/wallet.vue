<template>
  <div class="snow-page">
    <div class="snow-inner">
      <a-form ref="searchFormRef" auto-label-width :model="formData.form">
        <a-row :gutter="16">
          <a-col :xs="24" :sm="24" :md="12" :lg="12" :xl="6" :xxl="6">
            <a-form-item field="name" label="钱包名称">
              <a-input v-model="formData.form.name" placeholder="请输入名称" allow-clear />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="24" :md="12" :lg="12" :xl="6" :xxl="6">
            <a-form-item field="qrcode" label="钱包地址">
              <a-input v-model="formData.form.address" placeholder="请输入钱包地址" allow-clear />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="24" :md="12" :lg="12" :xl="6" :xxl="6">
            <a-form-item field="trade_type" label="交易类型">
              <a-select v-model="formData.form.trade_type" placeholder="请选择交易类型" allow-clear allow-search>
                <a-option v-for="item in tradeTypeOptions" :key="item.value" :value="item.value">
                  {{ item.label }}
                </a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="24" :md="12" :lg="12" :xl="6" :xxl="6">
            <a-space class="search-btn" wrap>
              <a-button type="primary" @click="getCommonTableList">
                <template #icon><icon-search /></template>
                查询
              </a-button>
              <a-button @click="onReset">
                <template #icon><icon-refresh /></template>
                重置
              </a-button>
              <a-button type="primary" status="success" @click="onAdd">
                <template #icon><icon-plus /></template>
                新增钱包
              </a-button>
            </a-space>
          </a-col>
        </a-row>
      </a-form>

      <a-table
        row-key="id"
        size="small"
        :bordered="{ cell: true }"
        :scroll="{ x: '100%', y: '100%', minWidth: 1180 }"
        :loading="loading"
        :columns="columns"
        :data="data"
        v-model:selectedKeys="selectedKeys"
        :pagination="pagination"
        @page-change="pageChange"
        @page-size-change="pageSizeChange"
      >
        <template #kind="{ record }">
          <a-tag size="small" :color="record.kind === 'multi' ? 'purple' : 'blue'">
            {{ record.kind === "multi" ? "多链钱包" : "基本钱包" }}
          </a-tag>
        </template>

        <template #trade_types="{ record }">
          <a-space wrap>
            <a-tag v-for="item in record.trade_types_display" :key="item" size="small" color="arcoblue">{{ item }}</a-tag>
          </a-space>
        </template>

        <template #address="{ record }">
          <div class="address-cell">
            <a-typography-text copyable class="address-text">
              {{ record.address }}
            </a-typography-text>
          </div>
        </template>

        <template #status="{ record }">
          <a-tag size="small" :color="record.status === 1 ? 'green' : 'red'">
            {{ record.status === 1 ? "启用" : "停用" }}
          </a-tag>
        </template>

        <template #other_notify="{ record }">
          <a-tag size="small" :color="record.other_notify === 1 ? 'arcoblue' : 'gray'">
            {{ record.other_notify === 1 ? "开启" : "关闭" }}
          </a-tag>
        </template>

        <template #optional="{ record }">
          <a-space wrap>
            <a-button size="mini" type="primary" @click="showDetail(record)">详情</a-button>
            <a-button size="mini" @click="onMod(record)">修改</a-button>
            <a-popconfirm content="确定删除这条数据吗?" type="warning" @ok="onDelete(record)">
              <a-button size="mini" type="primary" status="danger">删除</a-button>
            </a-popconfirm>
          </a-space>
        </template>
      </a-table>
    </div>
  </div>

  <a-modal :width="formDialogWidth" v-model:visible="open" @close="afterClose" @ok="addWallet" @cancel="afterClose">
    <template #title>{{ title }}</template>
    <a-form ref="addFormRef" auto-label-width :layout="formLayout" :rules="rules" :model="addFrom">
      <a-form-item field="name" label="钱包名称" validate-trigger="blur">
        <a-input v-model="addFrom.name" placeholder="请输入钱包名称" allow-clear />
      </a-form-item>
      <a-form-item field="address" label="钱包地址" validate-trigger="blur">
        <a-input v-model="addFrom.address" placeholder="请输入钱包地址" allow-clear />
      </a-form-item>
      <a-form-item field="kind" label="钱包类型">
        <a-radio-group v-model="addFrom.kind">
          <a-radio value="single">基本钱包地址</a-radio>
          <a-radio value="multi">多链钱包地址</a-radio>
        </a-radio-group>
      </a-form-item>
      <a-form-item
        v-if="addFrom.kind === 'single'"
        field="trade_type"
        label="交易类型"
        :rules="[{ required: true, message: '请选择交易类型' }]"
      >
        <a-select v-model="addFrom.trade_type" placeholder="请选择" allow-clear allow-search>
          <a-option v-for="item in tradeTypeOptions" :key="item.value" :value="item.value">
            {{ item.label }}
          </a-option>
        </a-select>
      </a-form-item>
      <template v-else>
        <a-form-item
          field="base"
          label="基础分类"
          :rules="[{ required: true, message: '请选择基础分类' }]"
        >
          <a-select v-model="addFrom.base" placeholder="请选择基础分类" @change="onAddBaseChange">
            <a-option v-for="item in walletBaseTypeOptions" :key="item.value" :value="item.value">
              {{ item.label }}
            </a-option>
          </a-select>
        </a-form-item>
        <a-form-item
          field="trade_types"
          label="支持链"
          :rules="[{ required: true, type: 'array', min: 1, message: '请选择支持链' }]"
        >
          <a-select v-model="addFrom.trade_types" multiple placeholder="请选择支持链" allow-search>
            <a-option v-for="item in addTradeTypeOptions" :key="item.value" :value="item.value">
              {{ item.label }}
            </a-option>
          </a-select>
        </a-form-item>
      </template>
      <a-form-item
        field="other_notify"
        label="其他通知"
        extra="开启后，当该钱包发生非订单转入/转出或 TRON 资源变化时，会通过已配置的通知渠道提醒。普通订单收款通知不依赖此开关。"
      >
        <a-select v-model="addFrom.other_notify" placeholder="请选择" allow-clear>
          <a-option :value="0">关闭</a-option>
          <a-option :value="1">启用</a-option>
        </a-select>
      </a-form-item>
      <a-form-item field="remark" label="备注信息" validate-trigger="blur">
        <a-textarea v-model="addFrom.remark" placeholder="请输入备注信息" allow-clear :auto-size="{ minRows: 2, maxRows: 4 }" />
      </a-form-item>
    </a-form>
  </a-modal>

  <a-modal :width="formDialogWidth" v-model:visible="modOpen" @close="afterModClose" @ok="modWallet" @cancel="afterModClose">
    <template #title>{{ modTitle }}</template>
    <a-form ref="modFormRef" auto-label-width :layout="formLayout" :rules="rules" :model="modFrom">
      <a-form-item field="name" label="钱包名称" validate-trigger="blur">
        <a-input v-model="modFrom.name" placeholder="请输入钱包名称" allow-clear />
      </a-form-item>
      <a-form-item field="address" label="钱包地址" validate-trigger="blur">
        <a-input v-model="modFrom.address" placeholder="请输入钱包地址" allow-clear />
      </a-form-item>
      <a-form-item field="kind" label="钱包类型">
        <a-radio-group v-model="modFrom.kind">
          <a-radio value="single">基本钱包地址</a-radio>
          <a-radio value="multi">多链钱包地址</a-radio>
        </a-radio-group>
      </a-form-item>
      <a-form-item
        v-if="modFrom.kind === 'single'"
        field="trade_type"
        label="交易类型"
        :rules="[{ required: true, message: '请选择交易类型' }]"
      >
        <a-select v-model="modFrom.trade_type" placeholder="请选择" allow-clear allow-search>
          <a-option v-for="item in tradeTypeOptions" :key="item.value" :value="item.value">
            {{ item.label }}
          </a-option>
        </a-select>
      </a-form-item>
      <template v-else>
        <a-form-item
          field="base"
          label="基础分类"
          :rules="[{ required: true, message: '请选择基础分类' }]"
        >
          <a-select v-model="modFrom.base" placeholder="请选择基础分类" @change="onModBaseChange">
            <a-option v-for="item in walletBaseTypeOptions" :key="item.value" :value="item.value">
              {{ item.label }}
            </a-option>
          </a-select>
        </a-form-item>
        <a-form-item
          field="trade_types"
          label="支持链"
          :rules="[{ required: true, type: 'array', min: 1, message: '请选择支持链' }]"
        >
          <a-select v-model="modFrom.trade_types" multiple placeholder="请选择支持链" allow-search>
            <a-option v-for="item in modTradeTypeOptions" :key="item.value" :value="item.value">
              {{ item.label }}
            </a-option>
          </a-select>
        </a-form-item>
      </template>
      <a-form-item field="status" label="收款状态">
        <a-select v-model="modFrom.status" placeholder="请选择" allow-clear>
          <a-option :value="1">启用</a-option>
          <a-option :value="0">停用</a-option>
        </a-select>
      </a-form-item>
      <a-form-item
        field="other_notify"
        label="其他通知"
        extra="开启后，当该钱包发生非订单转入/转出或 TRON 资源变化时，会通过已配置的通知渠道提醒。普通订单收款通知不依赖此开关。"
      >
        <a-select v-model="modFrom.other_notify" placeholder="请选择" allow-clear>
          <a-option :value="0">关闭</a-option>
          <a-option :value="1">开启</a-option>
        </a-select>
      </a-form-item>
      <a-form-item field="remark" label="备注信息" validate-trigger="blur">
        <a-textarea v-model="modFrom.remark" placeholder="请输入备注信息" allow-clear :auto-size="{ minRows: 2, maxRows: 4 }" />
      </a-form-item>
    </a-form>
  </a-modal>

  <a-modal
    :width="detailDialogWidth"
    v-model:visible="detailVisible"
    @close="closeDetail"
    @cancel="closeDetail"
    :footer="false"
    unmount-on-close
  >
    <template #title>
      <div class="detail-modal-title">
        <icon-star />
        <span>钱包详情信息</span>
      </div>
    </template>

    <div class="detail-content">
      <a-card class="detail-card" title="基础信息" :bordered="false">
        <template #extra>
          <a-tag size="medium" :color="detailData.status === 1 ? 'green' : 'red'" class="status-tag">
            <icon-check-circle v-if="detailData.status === 1" />
            <icon-close-circle v-else />
            {{ detailData.status === 1 ? "启用中" : "已停用" }}
          </a-tag>
        </template>

        <a-row :gutter="24">
          <a-col :xs="24" :sm="24" :md="12">
            <div class="detail-item">
              <div class="detail-label">
                <icon-idcard />
                <span>钱包ID</span>
              </div>
              <div class="detail-value">{{ detailData.id }}</div>
            </div>
          </a-col>
          <a-col :xs="24" :sm="24" :md="12">
            <div class="detail-item">
              <div class="detail-label">
                <icon-user />
                <span>钱包名称</span>
              </div>
              <div class="detail-value">{{ detailData.name }}</div>
            </div>
          </a-col>
        </a-row>

        <a-row :gutter="24">
          <a-col :xs="24" :sm="24" :md="12">
            <div class="detail-item">
              <div class="detail-label">
                <icon-apps />
                <span>钱包类型</span>
              </div>
              <div class="detail-value">
                <a-tag :color="detailData.kind === 'multi' ? 'purple' : 'blue'">
                  {{ detailData.kind === "multi" ? "多链钱包" : "基本钱包" }}
                </a-tag>
              </div>
            </div>
          </a-col>
          <a-col v-if="detailData.kind === 'multi'" :xs="24" :sm="24" :md="12">
            <div class="detail-item">
              <div class="detail-label">
                <icon-layers />
                <span>基础分类</span>
              </div>
              <div class="detail-value">{{ baseLabel(detailData.base) }}</div>
            </div>
          </a-col>
        </a-row>

        <a-row :gutter="24">
          <a-col :xs="24" :sm="24" :md="24">
            <div class="detail-item">
              <div class="detail-label">
                <icon-location />
                <span>钱包地址</span>
              </div>
              <div class="detail-value address-value">
                <a-typography-text copyable>{{ detailData.address }}</a-typography-text>
              </div>
            </div>
          </a-col>
        </a-row>

        <a-row :gutter="24">
          <a-col :xs="24" :sm="24" :md="12">
            <div class="detail-item">
              <div class="detail-label">
                <icon-swap />
                <span>交易类型</span>
              </div>
              <div class="detail-value">
                <a-space wrap>
                  <a-tag v-for="item in detailTradeTypes" :key="item" color="blue">{{ item }}</a-tag>
                </a-space>
              </div>
            </div>
          </a-col>
          <a-col :xs="24" :sm="24" :md="12">
            <div class="detail-item">
              <div class="detail-label">
                <icon-notification />
                <span>监控状态</span>
              </div>
              <div class="detail-value">
                <a-tag :color="detailData.other_notify === 1 ? 'arcoblue' : 'gray'">
                  <icon-eye v-if="detailData.other_notify === 1" />
                  <icon-eye-invisible v-else />
                  {{ detailData.other_notify === 1 ? "已开启" : "已关闭" }}
                </a-tag>
              </div>
            </div>
          </a-col>
        </a-row>
      </a-card>

      <a-card class="detail-card" title="备注信息" :bordered="false" v-if="detailData.remark">
        <div class="remark-content">
          <icon-message />
          <span>{{ detailData.remark }}</span>
        </div>
      </a-card>

      <a-card class="detail-card" title="时间信息" :bordered="false" v-if="detailData.created_at || detailData.updated_at">
        <a-row :gutter="24">
          <a-col :xs="24" :sm="24" :md="12" v-if="detailData.created_at">
            <div class="detail-item">
              <div class="detail-label">
                <icon-plus-circle />
                <span>创建时间</span>
              </div>
              <div class="detail-value">{{ detailData.created_at }}</div>
            </div>
          </a-col>
          <a-col :xs="24" :sm="24" :md="12" v-if="detailData.updated_at">
            <div class="detail-item">
              <div class="detail-label">
                <icon-edit />
                <span>更新时间</span>
              </div>
              <div class="detail-value">{{ detailData.updated_at }}</div>
            </div>
          </a-col>
        </a-row>
      </a-card>
    </div>
  </a-modal>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from "vue";
import { getWalletListAPI, delWalletAPI, addWalletAPI, modWalletAPI } from "@/api/modules/wallet/index";
import { List, FormData, Pagination, AddForm, ModForm, WalletBaseOption } from "./config";
import { Notification } from "@arco-design/web-vue";
import { useUserInfoStore } from "@/store/modules/user-info";
import { useWalletDetail } from "./detail";
import { useLayoutModel } from "@/hooks/useLayoutModel";

const userStores = useUserInfoStore();
const { detailVisible, detailData, showDetail, closeDetail } = useWalletDetail();
const { dialogWidth, formLayout } = useLayoutModel();
const formDialogWidth = computed(() => dialogWidth("40%"));
const detailDialogWidth = computed(() => dialogWidth("720px"));

const tradeTypeMap = computed<Record<string, string>>(() => userStores.trade_type || {});
const tradeTypeOptions = computed(() => Object.entries(tradeTypeMap.value).map(([value, label]) => ({ value, label })));
const walletBaseTypeOptions = computed<WalletBaseOption[]>(() => userStores.wallet_base_types || []);
const walletBaseTypeMap = computed<Record<string, WalletBaseOption>>(() => {
  const map: Record<string, WalletBaseOption> = {};
  walletBaseTypeOptions.value.forEach(item => {
    map[item.value] = item;
  });
  return map;
});

const tradeLabel = (value: string) => tradeTypeMap.value[value] || value;
const baseLabel = (value: string) => walletBaseTypeMap.value[value]?.label || value || "--";
const buildTradeTypeOptions = (base: string) => {
  const option = walletBaseTypeMap.value[base];
  if (!option) return [];
  return option.trade_types.map(value => ({ value, label: tradeLabel(value) }));
};

const addTradeTypeOptions = computed(() => buildTradeTypeOptions(addFrom.value.base));
const modTradeTypeOptions = computed(() => buildTradeTypeOptions(modFrom.value.base));
const detailTradeTypes = computed(() => {
  const values = detailData.value.trade_types?.length ? detailData.value.trade_types : detailData.value.trade_type ? [detailData.value.trade_type] : [];
  return values.map(tradeLabel);
});

const formatWalletRow = (record: List) => ({
  ...record,
  trade_types: record.trade_types || (record.trade_type ? [record.trade_type] : []),
  trade_types_display: (record.trade_types || (record.trade_type ? [record.trade_type] : [])).map(tradeLabel)
});

const formData = reactive<FormData>({
  form: { name: "", trade_type: "", address: "" },
  search: false
});

const selectedKeys = ref<string[]>([]);
const loading = ref(false);
const data = reactive<any[]>([]);
const pagination = ref<Pagination>({
  showPageSize: true,
  showTotal: true,
  current: 1,
  pageSize: 10,
  total: 10
});

const columns = [
  { title: "ID", align: "center", dataIndex: "id", width: 80 },
  { title: "名称", align: "center", dataIndex: "name", width: 180 },
  { title: "钱包类型", align: "center", dataIndex: "kind", slotName: "kind", width: 120 },
  { title: "交易类型", align: "center", dataIndex: "trade_types_display", slotName: "trade_types", width: 260 },
  { title: "钱包地址", align: "center", dataIndex: "address", slotName: "address", width: 320, ellipsis: true },
  { title: "收款状态", dataIndex: "status", align: "center", slotName: "status", width: 100 },
  { title: "其它通知", dataIndex: "other_notify", align: "center", slotName: "other_notify", width: 100 },
  { title: "操作", slotName: "optional", align: "center", fixed: "right", width: 200 }
];

const rules = {
  name: [{ required: true, message: "请输入钱包名称" }],
  address: [{ required: true, message: "请输入钱包地址" }]
};


const addFormRef = ref();
const modFormRef = ref();
const title = ref("");
const modTitle = ref("");
const open = ref(false);
const modOpen = ref(false);

const addFrom = ref<AddForm>({
  name: "",
  address: "",
  trade_type: "",
  trade_types: [],
  kind: "single",
  base: "generic",
  remark: "",
  other_notify: 0
});

const modFrom = ref<ModForm>({
  id: 0,
  name: "",
  address: "",
  trade_type: "",
  trade_types: [],
  kind: "single",
  base: "generic",
  remark: "",
  other_notify: 0,
  status: 1
});

const normalizeSingleForm = (form: AddForm | ModForm) => {
  form.base = "generic";
  form.trade_types = [];
};

const normalizeMultiForm = (form: AddForm | ModForm) => {
  form.trade_type = "";
};

const onAddBaseChange = () => {
  addFrom.value.trade_types = addFrom.value.trade_types.filter(item => addTradeTypeOptions.value.some(option => option.value === item));
};

const onModBaseChange = () => {
  modFrom.value.trade_types = modFrom.value.trade_types.filter(item => modTradeTypeOptions.value.some(option => option.value === item));
};

const pageChange = (page: number) => {
  pagination.value.current = page;
  getCommonTableList();
};

const pageSizeChange = (pageSize: number) => {
  pagination.value.pageSize = pageSize;
  getCommonTableList();
};

const onReset = () => {
  formData.form = { name: "", trade_type: "", address: "" };
  getCommonTableList();
};

const getCommonTableList = async () => {
  try {
    loading.value = true;
    const res = await getWalletListAPI({
      page: pagination.value.current,
      size: pagination.value.pageSize,
      sort: "desc",
      keyword: "",
      name: formData.form.name,
      trade_type: formData.form.trade_type,
      address: formData.form.address,
      status: 99
    });

    data.length = 0;
    data.push(...(res.data || []).map(formatWalletRow));
    pagination.value.total = res.total;
  } finally {
    loading.value = false;
  }
};

const onDelete = async (record: List) => {
  try {
    await delWalletAPI({ id: record.id });
    getCommonTableList();
    Notification.success("删除成功");
  } catch (error) {
    Notification.error(error as any);
  }
};

const onAdd = () => {
  title.value = "新增钱包";
  open.value = true;
};

const onMod = (record: List) => {
  modTitle.value = "修改钱包";
  modFrom.value = {
    id: record.id,
    name: record.name,
    address: record.address,
    trade_type: record.trade_type || "",
    trade_types: record.trade_types || (record.trade_type ? [record.trade_type] : []),
    kind: record.kind || ((record.trade_types || []).length > 1 ? "multi" : "single"),
    base: record.base || "generic",
    remark: record.remark || "",
    other_notify: record.other_notify || 0,
    status: record.status
  };
  if (modFrom.value.kind === "single") {
    normalizeSingleForm(modFrom.value);
  } else {
    normalizeMultiForm(modFrom.value);
  }
  modOpen.value = true;
};

const afterClose = () => {
  addFormRef.value?.resetFields();
  addFrom.value = {
    name: "",
    address: "",
    trade_type: "",
    trade_types: [],
    kind: "single",
    base: "generic",
    remark: "",
    other_notify: 0
  };
};

const afterModClose = () => {
  modFormRef.value?.resetFields();
  modFrom.value = {
    id: 0,
    name: "",
    address: "",
    trade_type: "",
    trade_types: [],
    kind: "single",
    base: "generic",
    remark: "",
    other_notify: 0,
    status: 1
  };
};

const buildSubmitData = (form: AddForm | ModForm) => {
  if (form.kind === "multi") {
    return {
      ...form,
      trade_type: "",
      trade_types: [...form.trade_types],
      other_notify: form.other_notify ?? 0
    };
  }

  return {
    ...form,
    base: "generic",
    trade_types: [],
    other_notify: form.other_notify ?? 0
  };
};

const addWallet = async () => {
  if (addFrom.value.kind === "single") normalizeSingleForm(addFrom.value);
  else normalizeMultiForm(addFrom.value);

  const state = await addFormRef.value.validate();
  if (state) return;

  try {
    await addWalletAPI(buildSubmitData(addFrom.value));
    open.value = false;
    getCommonTableList();
    Notification.success("添加成功");
  } catch (error) {
    Notification.error(error as any);
  }
};

const modWallet = async () => {
  if (modFrom.value.kind === "single") normalizeSingleForm(modFrom.value);
  else normalizeMultiForm(modFrom.value);

  const state = await modFormRef.value.validate();
  if (state) return;

  try {
    await modWalletAPI(buildSubmitData(modFrom.value));
    modOpen.value = false;
    getCommonTableList();
    Notification.success("修改成功");
  } catch (error) {
    Notification.error(error as any);
  }
};

getCommonTableList();
</script>

<style lang="scss" scoped>
.search-btn {
  margin-bottom: 20px;
}

.address-cell {
  max-width: 250px;

  .address-text {
    font-family: "Monaco", "Menlo", "Consolas", monospace;
    font-size: 12px;
    word-break: break-all;
    line-height: 1.4;

    :deep(.arco-typography-operation-copy) {
      color: $color-link;
      margin-left: 4px;

      &:hover {
        color: rgb(var(--link-5));
      }
    }
  }
}

.detail-modal-title {
  display: flex;
  align-items: center;
  gap: 8px;

  span {
    font-weight: 600;
    font-size: 16px;
  }
}

.detail-content {
  padding: 8px 0;

  .detail-card {
    margin-bottom: 16px;
    box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);

    &:last-child {
      margin-bottom: 0;
    }

    :deep(.arco-card-header) {
      border-bottom: 1px solid var(--color-border-2);
      padding: 16px 20px 12px;

      .arco-card-header-title {
        font-weight: 600;
        color: var(--color-text-1);
      }
    }

    :deep(.arco-card-body) {
      padding: 20px;
    }
  }

  .detail-item {
    margin-bottom: 20px;

    &:last-child {
      margin-bottom: 8px;
    }

    .detail-label {
      display: flex;
      align-items: center;
      gap: 6px;
      margin-bottom: 8px;
      font-size: 13px;
      color: var(--color-text-3);
      font-weight: 500;

      .arco-icon {
        font-size: 14px;
        color: var(--color-text-4);
      }
    }

    .detail-value {
      font-size: 14px;
      color: var(--color-text-1);
      font-weight: 500;
      min-height: 22px;
      display: flex;
      align-items: center;

      :deep(.arco-typography-operation-copy) {
        color: $color-link;
        margin-left: 4px;

        &:hover {
          color: rgb(var(--link-5));
        }
      }

      &.address-value {
        word-break: break-all;
        font-family: "Monaco", "Menlo", monospace;
        font-size: 13px;

        :deep(.arco-typography) {
          font-family: inherit;
          font-size: inherit;
        }
      }
    }
  }

  .status-tag {
    display: flex;
    align-items: center;
    gap: 4px;
    font-weight: 500;

    .arco-icon {
      font-size: 12px;
    }
  }

  .remark-content {
    display: flex;
    align-items: flex-start;
    gap: 8px;
    padding: 12px 16px;
    background-color: var(--color-fill-2);
    border-radius: 6px;
    line-height: 1.6;

    .arco-icon {
      margin-top: 2px;
      color: var(--color-text-3);
      flex-shrink: 0;
    }

    span {
      color: var(--color-text-2);
    }
  }
}

@media (max-width: 1200px) {
  :deep(.arco-table-th),
  :deep(.arco-table-td) {
    padding: 8px 6px !important;
    font-size: 12px;
  }
}

@media (max-width: 768px) {
  :deep(.arco-modal) {
    width: 95vw !important;
    margin: 10px;
  }

  :deep(.arco-table-th),
  :deep(.arco-table-td) {
    padding: 6px 4px !important;
    font-size: 11px;
  }

  .detail-content {
    .detail-card :deep(.arco-card-body) {
      padding: 16px;
    }

    .detail-item .detail-value.address-value {
      font-size: 12px;
    }
  }
}
</style>
