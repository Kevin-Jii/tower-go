<template>
  <div class="flex min-h-0 flex-1 flex-col gap-4">
    <div class="flex flex-col md:flex-row md:items-end gap-3 justify-between">
      <h2 class="page-title">门店记账</h2>
      <div class="flex flex-col sm:flex-row flex-wrap gap-2 w-full md:w-auto">
        <a-date-picker v-model="accountDate" value-format="YYYY-MM-DD" class="w-full sm:w-40" :allow-clear="false"
          :disabled-date="disabledFutureDate" />
        <BaseSelect v-model="paymentStatusFilter" class="w-full sm:w-32" :options="paymentStatusFilterOptions" />
        <BaseButton variant="primary" @click="reloadAll">查询</BaseButton>
        <BaseButton v-permission="'store:account:add'" variant="primary" @click="openCreate">快速记账</BaseButton>
        <BaseButton v-permission="'store:account:add'" variant="secondary" @click="openCustomCreate">自定义记账</BaseButton>
        <BaseButton v-permission="'store:account:edit'" variant="secondary" @click="openConsumableProductManage">消耗品维护
        </BaseButton>
      </div>
    </div>

    <BaseTable :columns="columns" :data="(list as unknown) as Record<string, unknown>[]" :loading="loading"
      min-width="1080px" class="min-h-0 flex-1">
      <template #cell-channel="{ row }">
        <span class="store-account-channel">
          <span class="store-account-channel__icon" :class="channelIconClass((row as StoreAccount).channel)">
            {{ channelIconText((row as StoreAccount).channel) }}
          </span>
          <span class="store-account-channel__label">{{ channelLabel((row as StoreAccount).channel) }}</span>
        </span>
      </template>
      <template #cell-member="{ row }">
        {{ memberLabel(row as StoreAccount) }}
      </template>
      <template #cell-payment_status="{ row }">
        <span class="store-account-pay-status">
          <span class="store-account-pay-status__dot" :class="paymentStatusDotClass((row as StoreAccount).payment_status)" />
          <span>{{ paymentStatusLabel((row as StoreAccount).payment_status) }}</span>
        </span>
      </template>
      <template #cell-operator="{ row }">
        {{ operatorLabel(row as StoreAccount) }}
      </template>
      <template #cell-net_income_amount="{ row }">
        {{ formatMoney((row as StoreAccount).net_income_amount) }}
      </template>
      <template #cell-account_date="{ row }">
        {{ formatDate((row as StoreAccount).account_date) }}
      </template>
      <template #cell-actions="{ row }">
        <BaseTableRowActions :actions="accountRowActions(row as StoreAccount)" />
      </template>
    </BaseTable>

    <div class="flex shrink-0 justify-end">
      <BasePagination :page="page" :page-size="pageSize" :total="total" @update:page="(p) => (page = p)"
        @update:page-size="(s) => (pageSize = s)" />
    </div>

    <BaseDialog v-model="customCreateDlg" title="自定义记账" max-width="min(720px, 96vw)">
      <p class="m-0 mb-3 text-xs text-[var(--color-text-3)]">
        商品明细为手写描述，不关联系统商品与库存；请填写单价与小计依据。渠道、会员、支付状态、备注与快速记账一致。
      </p>
      <div class="space-y-4">
        <BaseFormItem label="渠道" required>
          <BaseSelect :model-value="customForm.channel" :options="channelOptions" placeholder="请选择销售渠道"
            @update:model-value="onCustomChannelChange" />
        </BaseFormItem>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <BaseFormItem label="订单编号">
            <BaseInput v-model="customForm.order_no" placeholder="如 美团1号、淘宝闪购1号" />
          </BaseFormItem>
          <BaseFormItem v-if="isTakeoutChannel(customForm.channel)" label="收入金额">
            <BaseNumberInput v-model="customForm.income_amount" :min="0" :step="0.01" placeholder="留空按明细合计" />
          </BaseFormItem>
        </div>
        <BaseFormItem label="绑定会员">
          <BaseSelect v-model="customForm.member_id" :options="memberOptionsWithNone" placeholder="可选，默认不绑定" />
        </BaseFormItem>
        <BaseFormItem label="支付状态">
          <BaseSelect v-model="customForm.payment_status" :options="paymentStatusOptions" />
        </BaseFormItem>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <BaseFormItem label="跑腿订单">
            <BaseSwitch v-model="customForm.is_errand_order" :active-value="1" :inactive-value="0" />
          </BaseFormItem>
          <BaseFormItem v-if="customForm.is_errand_order === 1" label="跑腿费用" required>
            <BaseNumberInput v-model="customForm.errand_fee" :min="0" :step="0.01" />
          </BaseFormItem>
          <BaseFormItem label="抹零金额">
            <BaseNumberInput v-model="customForm.round_amount" :min="0" :step="0.01" />
          </BaseFormItem>
          <BaseFormItem label="赠酒">
            <BaseSwitch v-model="customForm.is_gift_wine" :active-value="1" :inactive-value="0" />
          </BaseFormItem>
          <template v-if="customForm.is_gift_wine === 1">
            <BaseFormItem label="赠酒商品" required>
              <a-cascader v-model="customForm.gift_wine_product_path" :options="productCascaderOptions"
                placeholder="先选分类，再选商品" allow-clear :path-mode="true" :check-strictly="false"
                @change="onCustomGiftWineProductChange" />
            </BaseFormItem>
            <BaseFormItem label="赠酒规格" required>
              <BaseSelect v-model="customForm.gift_wine_unit" :options="giftWineUnitOptions(customForm)"
                :disabled="giftWineUnitOptions(customForm).length <= 1" placeholder="规格" />
            </BaseFormItem>
            <BaseFormItem label="赠酒数量" required>
              <BaseNumberInput v-model="customForm.gift_wine_quantity" :min="0.01" :step="1" />
            </BaseFormItem>
            <BaseFormItem label="赠酒成本">
              <BaseInput :model-value="formatMoney(giftWineCostAmount(customForm))" disabled />
            </BaseFormItem>
          </template>
        </div>
        <div class="flex items-center justify-between gap-2">
          <span class="text-sm font-medium text-slate-700">商品明细（任意描述）</span>
          <BaseButton variant="secondary" size="sm" @click="addCustomLine">加一行</BaseButton>
        </div>
        <div v-for="(line, idx) in customForm.lines" :key="idx"
          class="rounded border border-[var(--color-border-2)] p-3 flex flex-col gap-3">
          <BaseFormItem label="明细描述" required class="w-full">
            <BaseTextarea v-model="line.description" :rows="2" placeholder="可填写任意商品或服务说明" />
          </BaseFormItem>
          <div class="flex flex-wrap items-end gap-2">
            <BaseFormItem label="数量" required class="w-28">
              <BaseNumberInput v-model="line.quantity" :min="0.01" :step="0.01" />
            </BaseFormItem>
            <BaseFormItem label="单位" required class="w-28">
              <BaseInput v-model="line.unit" placeholder="如 瓶、次、项" />
            </BaseFormItem>
            <BaseFormItem label="单价" required class="w-32">
              <BaseNumberInput v-model="line.price" :min="0.01" :step="0.01" />
            </BaseFormItem>
            <BaseFormItem label="行备注" class="min-w-[140px] flex-1">
              <BaseInput v-model="line.line_remark" placeholder="可选" />
            </BaseFormItem>
            <BaseButton variant="ghost" size="sm" class="shrink-0" @click="removeCustomLine(idx)">移除</BaseButton>
          </div>
        </div>
        <BaseFormItem label="整单备注">
          <BaseTextarea v-model="customForm.remark" :rows="2" />
        </BaseFormItem>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="customCreateDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="saving" @click="submitCustomCreate">保存</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="createDlg" title="快速记账（多商品）" max-width="min(720px, 96vw)">
      <div class="space-y-4">
        <BaseFormItem label="渠道" required>
          <BaseSelect :model-value="cForm.channel" :options="channelOptions" placeholder="请选择销售渠道"
            @update:model-value="onCreateChannelChange" />
        </BaseFormItem>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <BaseFormItem label="订单编号">
            <BaseInput v-model="cForm.order_no" placeholder="如 美团1号、淘宝闪购1号" />
          </BaseFormItem>
          <BaseFormItem v-if="isTakeoutChannel(cForm.channel)" label="收入金额">
            <BaseNumberInput v-model="cForm.income_amount" :min="0" :step="0.01" placeholder="留空按明细合计" />
          </BaseFormItem>
        </div>
        <BaseFormItem label="绑定会员">
          <BaseSelect v-model="cForm.member_id" :options="memberOptionsWithNone" placeholder="可选，默认不绑定" />
        </BaseFormItem>
        <BaseFormItem label="支付状态">
          <BaseSelect v-model="cForm.payment_status" :options="paymentStatusOptions" />
        </BaseFormItem>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <BaseFormItem label="跑腿订单">
            <BaseSwitch v-model="cForm.is_errand_order" :active-value="1" :inactive-value="0" />
          </BaseFormItem>
          <BaseFormItem v-if="cForm.is_errand_order === 1" label="跑腿费用" required>
            <BaseNumberInput v-model="cForm.errand_fee" :min="0" :step="0.01" />
          </BaseFormItem>
          <BaseFormItem label="抹零金额">
            <BaseNumberInput v-model="cForm.round_amount" :min="0" :step="0.01" />
          </BaseFormItem>
          <BaseFormItem label="赠酒">
            <BaseSwitch v-model="cForm.is_gift_wine" :active-value="1" :inactive-value="0" />
          </BaseFormItem>
          <template v-if="cForm.is_gift_wine === 1">
            <BaseFormItem label="赠酒商品" required>
              <a-cascader v-model="cForm.gift_wine_product_path" :options="productCascaderOptions"
                placeholder="先选分类，再选商品" allow-clear :path-mode="true" :check-strictly="false"
                @change="onCreateGiftWineProductChange" />
            </BaseFormItem>
            <BaseFormItem label="赠酒规格" required>
              <BaseSelect v-model="cForm.gift_wine_unit" :options="giftWineUnitOptions(cForm)"
                :disabled="giftWineUnitOptions(cForm).length <= 1" placeholder="规格" />
            </BaseFormItem>
            <BaseFormItem label="赠酒数量" required>
              <BaseNumberInput v-model="cForm.gift_wine_quantity" :min="0.01" :step="1" />
            </BaseFormItem>
            <BaseFormItem label="赠酒成本">
              <BaseInput :model-value="formatMoney(giftWineCostAmount(cForm))" disabled />
            </BaseFormItem>
          </template>
        </div>
        <div class="flex items-center justify-between gap-2">
          <span class="text-sm font-medium text-slate-700">商品明细</span>
          <BaseButton variant="secondary" size="sm" @click="addCreateLine">加一行</BaseButton>
        </div>
        <div v-for="(line, idx) in cForm.lines" :key="idx"
          class="rounded border border-[var(--color-border-2)] p-3 flex flex-wrap items-end gap-2">
          <BaseFormItem label="商品" required class="min-w-[220px] flex-1">
            <a-cascader v-model="line.product_path" :options="productCascaderOptions" placeholder="先选分类，再选商品"
              allow-clear :path-mode="true" :check-strictly="false" @change="onCreateProductChange(idx)" />
          </BaseFormItem>
          <BaseFormItem label="数量" required class="w-28">
            <BaseNumberInput v-model="line.quantity" :min="0.01" :step="0.01" />
          </BaseFormItem>
          <BaseFormItem label="单位" class="w-28">
            <BaseSelect v-model="line.unit" :options="lineUnitOptions(line)"
              :disabled="lineUnitOptions(line).length <= 1" placeholder="单位" />
          </BaseFormItem>
          <BaseButton variant="ghost" size="sm" @click="removeCreateLine(idx)">移除</BaseButton>
        </div>
        <BaseFormItem label="备注">
          <BaseTextarea v-model="cForm.remark" :rows="2" />
        </BaseFormItem>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="createDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="saving" @click="submitCreate">保存</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="editDlg" title="编辑记账" max-width="min(440px, 96vw)">
      <div class="space-y-4">
        <BaseFormItem label="渠道">
          <BaseSelect v-model="eForm.channel" :options="channelOptions" placeholder="请选择销售渠道" />
        </BaseFormItem>
        <BaseFormItem label="订单编号">
          <BaseInput v-model="eForm.order_no" placeholder="如 美团1号、淘宝闪购1号" />
        </BaseFormItem>
        <BaseFormItem v-if="isTakeoutChannel(eForm.channel)" label="收入金额">
          <BaseNumberInput v-model="eForm.income_amount" :min="0" :step="0.01" placeholder="留空不修改" />
        </BaseFormItem>
        <BaseFormItem label="绑定会员">
          <BaseSelect v-model="eForm.member_id" :options="memberOptionsWithNone" placeholder="可选，默认不绑定" />
        </BaseFormItem>
        <BaseFormItem label="支付状态">
          <BaseSelect v-model="eForm.payment_status" :options="paymentStatusOptions" />
        </BaseFormItem>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <BaseFormItem label="跑腿订单">
            <BaseSwitch v-model="eForm.is_errand_order" :active-value="1" :inactive-value="0" />
          </BaseFormItem>
          <BaseFormItem v-if="eForm.is_errand_order === 1" label="跑腿费用" required>
            <BaseNumberInput v-model="eForm.errand_fee" :min="0" :step="0.01" />
          </BaseFormItem>
          <BaseFormItem label="抹零金额">
            <BaseNumberInput v-model="eForm.round_amount" :min="0" :step="0.01" />
          </BaseFormItem>
          <BaseFormItem label="赠酒">
            <BaseSwitch v-model="eForm.is_gift_wine" :active-value="1" :inactive-value="0" />
          </BaseFormItem>
          <template v-if="eForm.is_gift_wine === 1">
            <BaseFormItem label="赠酒商品" required>
              <a-cascader v-model="eForm.gift_wine_product_path" :options="productCascaderOptions"
                placeholder="先选分类，再选商品" allow-clear :path-mode="true" :check-strictly="false"
                @change="onEditGiftWineProductChange" />
            </BaseFormItem>
            <BaseFormItem label="赠酒规格" required>
              <BaseSelect v-model="eForm.gift_wine_unit" :options="giftWineUnitOptions(eForm)"
                :disabled="giftWineUnitOptions(eForm).length <= 1" placeholder="规格" />
            </BaseFormItem>
            <BaseFormItem label="赠酒数量" required>
              <BaseNumberInput v-model="eForm.gift_wine_quantity" :min="0.01" :step="1" />
            </BaseFormItem>
            <BaseFormItem label="赠酒成本">
              <BaseInput :model-value="formatMoney(giftWineCostAmount(eForm))" disabled />
            </BaseFormItem>
          </template>
        </div>
        <BaseFormItem label="标签编码">
          <BaseInput v-model="eForm.tag_code" />
        </BaseFormItem>
        <BaseFormItem label="标签名称">
          <BaseInput v-model="eForm.tag_name" />
        </BaseFormItem>
        <BaseFormItem label="备注">
          <BaseTextarea v-model="eForm.remark" :rows="2" />
        </BaseFormItem>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="editDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="saving" @click="submitEdit">保存</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="viewDlg" title="记账详情" max-width="min(720px, 96vw)">
      <div v-if="viewAccount" class="max-h-[70vh] overflow-y-auto space-y-4 pr-1">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-x-4 gap-y-2 text-sm">
          <div><span class="text-[var(--color-text-3)]">单号</span>：{{ viewAccount.account_no }}</div>
          <div><span class="text-[var(--color-text-3)]">记账日期</span>：{{ formatDate(viewAccount.account_date) }}</div>
          <div><span class="text-[var(--color-text-3)]">渠道</span>：{{ channelLabel(viewAccount.channel) }}</div>
          <div><span class="text-[var(--color-text-3)]">会员</span>：{{ memberLabel(viewAccount) }}</div>
          <div><span class="text-[var(--color-text-3)]">支付状态</span>：{{ paymentStatusLabel(viewAccount.payment_status) }}
          </div>
          <div><span class="text-[var(--color-text-3)]">订单号</span>：{{ viewAccount.order_no || '-' }}</div>
          <div><span class="text-[var(--color-text-3)]">总金额</span>：{{ formatMoney(viewAccount.total_amount) }}</div>
          <div><span class="text-[var(--color-text-3)]">其他支出</span>：{{ formatMoney(viewAccount.other_expense_amount) }}
          </div>
          <div><span class="text-[var(--color-text-3)]">抹零金额</span>：{{ formatMoney(viewAccount.round_amount) }}</div>
          <div><span class="text-[var(--color-text-3)]">赠酒</span>：{{ Number(viewAccount.is_gift_wine || 0) === 1 ?
            '是'
            : '否' }}</div>
          <div v-if="Number(viewAccount.is_gift_wine || 0) === 1" class="sm:col-span-2">
            <span class="text-[var(--color-text-3)]">赠酒商品</span>：{{ viewAccount.gift_wine_product_name ||
              (viewAccount.gift_wine_product_id ? `商品#${viewAccount.gift_wine_product_id}` : '-') }} /
            {{ viewAccount.gift_wine_unit || '-' }} × {{ viewAccount.gift_wine_quantity || 0 }}
          </div>
          <div><span class="text-[var(--color-text-3)]">赠酒成本</span>：{{ formatMoney(viewAccount.gift_wine_cost_amount) }}
          </div>
          <div><span class="text-[var(--color-text-3)]">跑腿订单</span>：{{ Number(viewAccount.is_errand_order || 0) === 1 ?
            '是'
            : '否' }}</div>
          <div><span class="text-[var(--color-text-3)]">跑腿费用</span>：{{ formatMoney(viewAccount.errand_fee) }}</div>
          <div><span class="text-[var(--color-text-3)]">商品成本</span>：{{ formatMoney(accountItemCost(viewAccount)) }}
          </div>
          <div><span class="text-[var(--color-text-3)]">耗材金额</span>：{{ formatMoney(accountConsumableAmount(viewAccount))
            }}
          </div>
          <div><span class="text-[var(--color-text-3)]">净收入</span>：{{ formatMoney(viewAccount.net_income_amount) }}
          </div>
          <div><span class="text-[var(--color-text-3)]">明细条数</span>：{{ viewAccount.item_count ??
            (viewAccount.items?.length
              ?? 0) }}</div>
          <div class="sm:col-span-2">
            <span class="text-[var(--color-text-3)]">标签</span>：{{ viewAccount.tag_name || viewAccount.tag_code || '-' }}
          </div>
          <div class="sm:col-span-2"><span class="text-[var(--color-text-3)]">备注</span>：{{ viewAccount.remark || '-' }}
          </div>
          <div class="sm:col-span-2 text-xs text-[var(--color-text-3)]">
            创建时间：{{ viewAccount.created_at ? formatDateTime(viewAccount.created_at) : '-' }}
          </div>
        </div>

        <div>
          <p class="m-0 mb-2 text-sm font-medium text-slate-800">商品明细</p>
          <div v-if="(viewAccount.items?.length ?? 0) === 0" class="text-sm text-[var(--color-text-3)]">暂无商品明细</div>
          <div v-else class="overflow-x-auto rounded border border-[var(--color-border-2)]">
            <table class="w-full min-w-[560px] border-collapse text-sm">
              <thead>
                <tr class="bg-[var(--color-fill-2)]">
                  <th class="border-b border-[var(--color-border-2)] px-2 py-2 text-left font-medium">商品</th>
                  <th class="border-b border-[var(--color-border-2)] px-2 py-2 text-left w-24">规格</th>
                  <th class="border-b border-[var(--color-border-2)] px-2 py-2 text-right w-20">数量</th>
                  <th class="border-b border-[var(--color-border-2)] px-2 py-2 text-center w-16">单位</th>
                  <th class="border-b border-[var(--color-border-2)] px-2 py-2 text-right w-24">单价</th>
                  <th class="border-b border-[var(--color-border-2)] px-2 py-2 text-right w-24">小计</th>
                  <th class="border-b border-[var(--color-border-2)] px-2 py-2 text-left min-w-[80px]">备注</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="it in viewAccount.items" :key="it.id">
                  <td class="border-b border-[var(--color-border-2)] px-2 py-1.5">
                    {{ it.product_name || (it.product_id ? `商品#${it.product_id}` : '—') }}
                  </td>
                  <td class="border-b border-[var(--color-border-2)] px-2 py-1.5">{{ it.spec || '-' }}</td>
                  <td class="border-b border-[var(--color-border-2)] px-2 py-1.5 text-right">{{ it.quantity }}</td>
                  <td class="border-b border-[var(--color-border-2)] px-2 py-1.5 text-center">{{ it.unit || '-' }}</td>
                  <td class="border-b border-[var(--color-border-2)] px-2 py-1.5 text-right">{{ formatMoney(it.price) }}
                  </td>
                  <td class="border-b border-[var(--color-border-2)] px-2 py-1.5 text-right">{{ formatMoney(it.amount)
                    }}
                  </td>
                  <td class="border-b border-[var(--color-border-2)] px-2 py-1.5">{{ it.remark || '-' }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div v-if="(viewAccount.consumables?.length ?? 0) > 0">
          <p class="m-0 mb-2 text-sm font-medium text-slate-800">消耗品</p>
          <div class="overflow-x-auto rounded border border-[var(--color-border-2)]">
            <table class="w-full min-w-[480px] border-collapse text-sm">
              <thead>
                <tr class="bg-[var(--color-fill-2)]">
                  <th class="border-b border-[var(--color-border-2)] px-2 py-2 text-left font-medium">商品</th>
                  <th class="border-b border-[var(--color-border-2)] px-2 py-2 text-right w-20">数量</th>
                  <th class="border-b border-[var(--color-border-2)] px-2 py-2 text-center w-16">单位</th>
                  <th class="border-b border-[var(--color-border-2)] px-2 py-2 text-right w-24">单价</th>
                  <th class="border-b border-[var(--color-border-2)] px-2 py-2 text-right w-24">小计</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="c in viewAccount.consumables" :key="c.id">
                  <td class="border-b border-[var(--color-border-2)] px-2 py-1.5">{{ c.product_name ||
                    `商品#${c.product_id}`
                    }}</td>
                  <td class="border-b border-[var(--color-border-2)] px-2 py-1.5 text-right">{{ c.quantity }}</td>
                  <td class="border-b border-[var(--color-border-2)] px-2 py-1.5 text-center">{{ c.unit || '-' }}</td>
                  <td class="border-b border-[var(--color-border-2)] px-2 py-1.5 text-right">{{ formatMoney(c.price) }}
                  </td>
                  <td class="border-b border-[var(--color-border-2)] px-2 py-1.5 text-right">{{ formatMoney(c.amount) }}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div class="rounded border border-[var(--color-border-2)] bg-[var(--color-fill-1)] px-3 py-2 text-sm">
          <span class="text-[var(--color-text-3)]">净利润口径：</span>
          销售额 {{ formatMoney(viewAccount.total_amount) }} - 其他支出 {{ formatMoney(viewAccount.other_expense_amount) }} -
          跑腿费用
          {{ formatMoney(viewAccount.errand_fee) }} - 商品成本 {{ formatMoney(accountItemCost(viewAccount)) }} - 耗材金额
          {{ formatMoney(accountConsumableAmount(viewAccount)) }} - 赠酒成本
          {{ formatMoney(viewAccount.gift_wine_cost_amount) }} - 抹零金额 {{ formatMoney(viewAccount.round_amount) }} =
          <span class="font-semibold text-emerald-700">{{ formatMoney(accountNetProfitBreakdown(viewAccount)) }}</span>
        </div>
      </div>
      <p v-else class="m-0 text-sm text-[var(--color-text-3)]">加载中…</p>
      <template #footer>
        <BaseButton variant="ghost" @click="viewDlg = false">关闭</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="consumableDlg" title="绑定消耗品" max-width="min(720px, 96vw)">
      <div class="space-y-4">
        <p class="m-0 text-sm text-slate-600">记账单：{{ consumableTarget?.account_no }}（绑定后会计入成本并扣减净利润）</p>
        <div class="rounded border border-slate-200 bg-slate-50 px-3 py-2 text-sm text-slate-700">
          消耗品成本合计：<span class="font-semibold text-slate-900">{{ formatMoney(consumableBindTotal) }}</span>
        </div>
        <div class="flex items-center justify-between gap-2">
          <span class="text-sm font-medium text-slate-700">消耗品明细</span>
          <div class="flex gap-2">
            <BaseButton variant="secondary" size="sm" @click="addConsumableLine">选择消耗品</BaseButton>
            <BaseButton variant="secondary" size="sm" @click="addCustomConsumableLine">自定义</BaseButton>
          </div>
        </div>
        <div v-for="(line, idx) in consumableLines" :key="idx"
          class="rounded border border-[var(--color-border-2)] p-3 flex flex-wrap items-end gap-2">
          <template v-if="line.kind === 'custom'">
            <BaseFormItem label="名称" required class="min-w-[220px] flex-1">
              <BaseInput v-model="line.name" placeholder="自定义消耗品名称" />
            </BaseFormItem>
            <BaseFormItem label="金额" required class="w-32">
              <BaseNumberInput v-model="line.amount" :min="0.01" :step="0.01" />
            </BaseFormItem>
          </template>
          <BaseFormItem v-else label="消耗品" required class="min-w-[220px] flex-1">
            <BaseSelect v-model="line.consumable_product_id" :options="consumableProductOptions" placeholder="选择消耗品" />
          </BaseFormItem>
          <BaseFormItem v-if="line.kind !== 'custom'" label="数量" required class="w-28">
            <BaseNumberInput v-model="line.quantity" :min="0.01" :step="1" :hide-button="false" />
          </BaseFormItem>
          <BaseFormItem v-if="line.kind !== 'custom'" label="成本小计" class="w-32">
            <BaseInput :model-value="formatMoney(consumableLineAmount(line))" disabled />
          </BaseFormItem>
          <BaseButton variant="ghost" size="sm" @click="removeConsumableLine(idx)">移除</BaseButton>
        </div>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="consumableDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="consumableSaving" @click="submitConsumables">保存消耗品</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="consumableProductDlg" title="消耗品维护" max-width="min(760px, 96vw)">
      <div class="flex max-h-[70vh] min-h-0 flex-col gap-3">
        <div class="flex flex-wrap items-center justify-between gap-2">
          <BaseInput v-model="consumableProductKeyword" class="w-56" placeholder="名称 / 备注" clearable
            @enter="reloadConsumableProducts" />
          <BaseButton variant="primary" @click="openConsumableProductCreate">新增消耗品</BaseButton>
        </div>
        <BaseTable :columns="consumableProductColumns"
          :data="(consumableProductRows as unknown) as Record<string, unknown>[]" :loading="consumableProductsLoading"
          min-width="620px" height="360px">
          <template #cell-cost_price="{ row }">{{ formatMoney((row as StoreAccountConsumableProduct).cost_price)
            }}</template>
          <template #cell-actions="{ row }">
            <BaseTableRowActions :actions="consumableProductActions(row as StoreAccountConsumableProduct)"
              :max-inline="2" />
          </template>
        </BaseTable>
        <div class="flex shrink-0 justify-end">
          <BasePagination :page="consumableProductPage" :page-size="consumableProductPageSize"
            :total="consumableProductTotal" @update:page="(p) => (consumableProductPage = p)"
            @update:page-size="(s) => (consumableProductPageSize = s)" />
        </div>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="consumableProductDlg = false">关闭</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="consumableProductFormDlg" :title="consumableProductEditId ? '编辑消耗品' : '新增消耗品'"
      max-width="min(460px, 96vw)">
      <div class="space-y-4">
        <BaseFormItem label="名称" required>
          <BaseInput v-model="consumableProductForm.name" placeholder="请输入消耗品名称" />
        </BaseFormItem>
        <BaseFormItem label="成本价" required>
          <BaseNumberInput v-model="consumableProductForm.cost_price" :min="0" :step="0.01" />
        </BaseFormItem>
        <BaseFormItem label="备注说明">
          <BaseTextarea v-model="consumableProductForm.remark" :rows="2" />
        </BaseFormItem>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="consumableProductFormDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="consumableProductSaving" @click="submitConsumableProduct">保存
        </BaseButton>
      </template>
    </BaseDialog>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import {
  BaseButton,
  BaseDialog,
  BaseFormItem,
  BaseInput,
  BaseNumberInput,
  BasePagination,
  BaseSelect,
  BaseSwitch,
  BaseTable,
  BaseTableRowActions,
  BaseTextarea,
} from '@/components/base'
import type { BaseTableColumn, TableRowAction } from '@/components/base/types'
import {
  bindStoreAccountConsumables,
  createStoreAccountConsumableProduct,
  createStoreAccount,
  deleteStoreAccountConsumableProduct,
  getStoreAccount,
  listStoreAccountConsumableProducts,
  listStoreAccounts,
  updateStoreAccountConsumableProduct,
  updateStoreAccount,
} from '@/api/storeAccount'
import { listDictDataByTypeCode } from '@/api/dict'
import { batchListProductUnitSpecs } from '@/api/supplierProduct'
import { listPurchasableProducts } from '@/api/storeSupplier'
import { listMembers } from '@/api/member'
import type { DictData, MemberRow, ProductUnitSpec, StoreAccount, StoreAccountConsumableProduct } from '@/api/types'
import { toast } from '@/feedback/toast'
import { confirmDialog } from '@/feedback/confirm'
import { useUserStore } from '@/store/user'
import { currentStoreIdOrUndefined } from '@/utils/currentStore'

const qc = useQueryClient()
const userStore = useUserStore()
const tenantStoreId = computed(() => currentStoreIdOrUndefined(userStore.userInfo, userStore.tenantId))

function dateText(t: Date): string {
  const y = t.getFullYear()
  const m = String(t.getMonth() + 1).padStart(2, '0')
  const d = String(t.getDate()).padStart(2, '0')
  return `${y}-${m}-${d}`
}

function businessDate(source: Date): Date {
  const t = new Date(source)
  if (t.getHours() < 5) {
    t.setDate(t.getDate() - 1)
  }
  t.setHours(0, 0, 0, 0)
  return t
}

const accountDate = ref(dateText(businessDate(new Date())))
const paymentStatusFilter = ref(0)

function disabledFutureDate(current?: Date): boolean {
  if (!current) return false
  const today = businessDate(new Date())
  const target = new Date(current)
  target.setHours(0, 0, 0, 0)
  return target.getTime() > today.getTime()
}

const page = ref(1)
const pageSize = ref(10)
const listKey = computed(
  () => ['store-accounts', tenantStoreId.value, page.value, pageSize.value, accountDate.value, paymentStatusFilter.value] as const,
)

const { data: pageData, isLoading: loading } = useQuery({
  queryKey: listKey,
  queryFn: () =>
    listStoreAccounts({
      page: page.value,
      page_size: pageSize.value,
      store_id: tenantStoreId.value,
      start_date: accountDate.value,
      end_date: accountDate.value,
      payment_status: paymentStatusFilter.value || undefined,
    }),
})

const list = computed(() => pageData.value?.list ?? [])
const total = computed(() => pageData.value?.total ?? 0)

const consumableProductPage = ref(1)
const consumableProductPageSize = ref(10)
const consumableProductKeyword = ref('')
const consumableProductParams = computed(() => ({
  page: consumableProductPage.value,
  page_size: consumableProductPageSize.value,
  store_id: tenantStoreId.value,
  keyword: consumableProductKeyword.value.trim() || undefined,
}))
const { data: consumableProductPageData, isLoading: consumableProductsLoading } = useQuery({
  queryKey: computed(() => ['store-account-consumable-products', consumableProductParams.value] as const),
  queryFn: () => listStoreAccountConsumableProducts(consumableProductParams.value),
})
const { data: allConsumableProductPageData } = useQuery({
  queryKey: computed(() => ['store-account-consumable-products-all', tenantStoreId.value] as const),
  queryFn: () => listStoreAccountConsumableProducts({ page: 1, page_size: 500, store_id: tenantStoreId.value }),
})
const consumableProductRows = computed(() => consumableProductPageData.value?.list ?? [])
const consumableProductTotal = computed(() => consumableProductPageData.value?.total ?? 0)
const allConsumableProducts = computed(() => allConsumableProductPageData.value?.list ?? [])
const consumableProductOptions = computed(() =>
  allConsumableProducts.value.map((item) => ({ label: `${item.name}（${formatMoney(item.cost_price)}）`, value: item.id })),
)
const consumableProductMap = computed(() => {
  const map = new Map<number, StoreAccountConsumableProduct>()
  for (const item of allConsumableProducts.value) map.set(item.id, item)
  return map
})

const { data: productData } = useQuery({
  queryKey: computed(() => ['store-account-products', tenantStoreId.value] as const),
  queryFn: () =>
    listPurchasableProducts({
      store_id: tenantStoreId.value,
    }),
})
const productList = computed(() => productData.value ?? [])
const productById = computed(() => {
  const map = new Map<number, (typeof productList.value)[number]>()
  for (const p of productList.value) map.set(p.id, p)
  return map
})
const productIdsKey = computed(() =>
  productList.value
    .map((p) => p.id)
    .sort((a, b) => a - b)
    .join(','),
)
const { data: unitSpecsData } = useQuery({
  queryKey: computed(() => ['store-account-product-unit-specs', productIdsKey.value] as const),
  queryFn: async () => {
    const ids = productList.value.map((p) => p.id)
    if (!ids.length) return [] as ProductUnitSpec[]
    return batchListProductUnitSpecs(ids)
  },
  enabled: computed(() => productList.value.length > 0),
})
const specsByProduct = computed(() => {
  const map = new Map<number, ProductUnitSpec[]>()
  for (const s of unitSpecsData.value ?? []) {
    if (!s.is_enabled) continue
    if (!map.has(s.product_id)) map.set(s.product_id, [])
    map.get(s.product_id)!.push(s)
  }
  for (const [, arr] of map.entries()) {
    arr.sort((a, b) => Number(a.factor_to_base) - Number(b.factor_to_base))
  }
  return map
})
const { data: unitData } = useQuery({
  queryKey: ['dict-data', 'product_unit'],
  queryFn: () => listDictDataByTypeCode('product_unit'),
})
const unitDict = computed(() => unitData.value ?? ([] as DictData[]))
const { data: channelData } = useQuery({
  queryKey: ['dict-data', 'sales_channel'],
  queryFn: () => listDictDataByTypeCode('sales_channel'),
})
const channelOptions = computed(() => (channelData.value ?? []).map((d) => ({ label: d.label, value: d.value })))
const channelDictMap = computed(() => {
  const map = new Map<string, string>()
  for (const d of channelData.value ?? []) {
    map.set(String(d.value), d.label || String(d.value))
  }
  return map
})
const takeoutChannelTokens = ['takeout', 'waimai', 'meituan', 'eleme', 'elm', 'taobao', 'tb', 'flash', 'shangou', 'jd', 'jingdong', '外卖', '美团', '饿了么', '淘宝', '闪购', '京东']
function isTakeoutChannel(channel: string | undefined): boolean {
  const raw = String(channel || '').trim()
  if (!raw) return false
  const label = channelDictMap.value.get(raw) || ''
  const text = `${raw} ${label}`.toLowerCase()
  return takeoutChannelTokens.some((token) => text.includes(token.toLowerCase()))
}
const { data: membersPageData } = useQuery({
  queryKey: ['store-account-members'],
  queryFn: () => listMembers({ page: 1, page_size: 200 }),
})
const memberList = computed(() => membersPageData.value?.list ?? ([] as MemberRow[]))
const memberOptions = computed(() =>
  memberList.value.map((m) => ({
    label: `${m.phone}${m.name ? `（${m.name}）` : ''}`,
    value: m.id,
  })),
)
const memberOptionsWithNone = computed(() => [{ label: '不绑定会员', value: 0 }, ...memberOptions.value])
const memberMap = computed(() => {
  const map = new Map<number, MemberRow>()
  for (const m of memberList.value) {
    map.set(m.id, m)
  }
  return map
})

const productCascaderOptions = computed(() => {
  const grouped = new Map<string, { id: number; name: string }[]>()
  for (const p of productList.value) {
    const cat = p.category?.name?.trim() || '未分类'
    if (!grouped.has(cat)) grouped.set(cat, [])
    grouped.get(cat)!.push({ id: p.id, name: p.name })
  }
  let idx = 0
  return Array.from(grouped.entries()).map(([cat, products]) => {
    idx += 1
    return {
      label: cat,
      value: `cat-${idx}`,
      children: products.map((p) => ({
        label: `${p.name}（#${p.id}）`,
        value: p.id,
      })),
    }
  })
})

function reloadAll(): void {
  page.value = 1
  void qc.invalidateQueries({ queryKey: ['store-accounts'] })
}

watch([accountDate, paymentStatusFilter], () => {
  page.value = 1
})

watch([page, pageSize], () => {
  void qc.invalidateQueries({ queryKey: ['store-accounts'] })
})

watch(
  () => tenantStoreId.value,
  () => {
    void qc.invalidateQueries({ queryKey: ['store-accounts'] })
  },
)

const columns: BaseTableColumn[] = [
  { key: 'channel', label: '渠道', prop: 'channel', width: '100px' },
  { key: 'member', label: '会员', width: '150px', ellipsis: true },
  { key: 'payment_status', label: '支付状态', width: '96px' },
  { key: 'total_amount', label: '销售额', prop: 'total_amount', width: '96px' },
  { key: 'operator', label: '操作人', width: '80px', ellipsis: true },
  { key: 'actions', label: '操作', width: '190px', minWidth: '140px', align: 'right', fixed: 'right' },
]

const consumableProductColumns: BaseTableColumn[] = [
  { key: 'id', label: 'ID', prop: 'id', width: '80px' },
  { key: 'name', label: '名称', prop: 'name', minWidth: '180px', ellipsis: true },
  { key: 'cost_price', label: '成本价', width: '110px', align: 'right' },
  { key: 'remark', label: '备注说明', prop: 'remark', minWidth: '180px', ellipsis: true },
  { key: 'actions', label: '操作', width: '120px', align: 'right' },
]

function formatDate(v: string): string {
  if (!v) return '-'
  return String(v).slice(0, 10)
}

function channelLabel(v: string | undefined): string {
  const key = String(v || '').trim()
  if (!key) return '-'
  return channelDictMap.value.get(key) || key
}

function channelMeta(v: string | undefined): { text: string; tone: string } {
  const key = String(v || '').trim().toLowerCase()
  const label = channelLabel(v).toLowerCase()
  const text = `${key} ${label}`
  if (text.includes('wechat') || text.includes('微信')) return { text: '微', tone: 'wechat' }
  if (text.includes('meituan') || text.includes('美团')) return { text: '美', tone: 'meituan' }
  if (text.includes('eleme') || text.includes('饿了么') || text.includes('elm')) return { text: '饿', tone: 'eleme' }
  if (text.includes('douyin') || text.includes('抖音')) return { text: '抖', tone: 'douyin' }
  if (text.includes('taobao') || text.includes('淘宝') || text.includes('shangou') || text.includes('闪购')) return { text: '淘', tone: 'taobao' }
  if (text.includes('xiaohongshu') || text.includes('小红书')) return { text: '红', tone: 'redbook' }
  if (text.includes('offline') || text.includes('线下') || text.includes('门店')) return { text: '店', tone: 'offline' }
  return { text: '其', tone: 'other' }
}

function channelIconText(v: string | undefined): string {
  return channelMeta(v).text
}

function channelIconClass(v: string | undefined): string {
  return `store-account-channel__icon--${channelMeta(v).tone}`
}

function memberLabel(row: StoreAccount): string {
  if (row.member) {
    const phone = String(row.member.phone || '').trim()
    const name = String(row.member.name || '').trim()
    if (phone && name) return `${phone}（${name}）`
    return phone || name || `会员#${row.member.id}`
  }
  const mid = Number(row.member_id || 0)
  if (mid > 0) {
    const m = memberMap.value.get(mid)
    if (m) {
      return `${m.phone}${m.name ? `（${m.name}）` : ''}`
    }
    return `会员#${mid}`
  }
  return '-'
}

function operatorLabel(row: StoreAccount): string {
  const operator = row.operator
  if (operator) {
    const nickname = String(operator.nickname || '').trim()
    const username = String(operator.username || '').trim()
    const phone = String(operator.phone || '').trim()
    return nickname || username || phone || `用户#${operator.id}`
  }
  return row.operator_id ? `用户#${row.operator_id}` : '-'
}

const paymentStatusOptions = [
  { label: '已支付', value: 1 },
  { label: '未支付', value: 2 },
]
const paymentStatusFilterOptions = computed(() => [{ label: '全部状态', value: 0 }, ...paymentStatusOptions])

function paymentStatusLabel(v: number | undefined): string {
  return Number(v) === 2 ? '未支付' : '已支付'
}

function paymentStatusDotClass(v: number | undefined): string {
  return Number(v) === 2 ? 'store-account-pay-status__dot--unpaid' : 'store-account-pay-status__dot--paid'
}

function canEditAccount(row: StoreAccount): boolean {
  if (typeof row.can_edit === 'boolean') return row.can_edit
  return isWithinBusinessDays(row.created_at, 1)
}

function canBindConsumables(row: StoreAccount): boolean {
  if (typeof row.can_bind_consumables === 'boolean') return row.can_bind_consumables
  return (row.consumables?.length ?? 0) === 0
}

function businessDateKey(date: Date): string {
  const d = new Date(date)
  if (d.getHours() < 5) d.setDate(d.getDate() - 1)
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

function isWithinBusinessDays(value: string | undefined, days: number): boolean {
  const s = String(value || '').trim()
  if (!s) return false
  const created = new Date(s)
  if (Number.isNaN(created.getTime())) return false
  const current = new Date()
  const createdKey = businessDateKey(created)
  const currentKey = businessDateKey(current)
  if (currentKey < createdKey) return false
  const last = new Date(`${createdKey}T00:00:00`)
  last.setDate(last.getDate() + Math.max(1, days) - 1)
  return currentKey <= businessDateKey(last)
}

function accountRowActions(row: StoreAccount): TableRowAction[] {
  const editable = canEditAccount(row)
  const consumableEditable = canBindConsumables(row)
  const actions: TableRowAction[] = []
  if (Number(row.payment_status || 1) === 2) {
    actions.push({
      label: '支付完成',
      permission: 'store:account:edit',
      onClick: () => void markAccountPaid(row),
    })
  }
  actions.push({ label: '详情', permission: 'store:account:list', onClick: () => openView(row) })
  if (consumableEditable) {
    actions.push({
      label: '绑定消耗品',
      permission: 'store:account:edit',
      onClick: () => openConsumableDlg(row),
    })
  }
  actions.push(
    {
      label: '编辑',
      permission: 'store:account:edit',
      disabled: !editable,
      onClick: () => openEdit(row),
    },
  )
  return actions
}

const createDlg = ref(false)
const customCreateDlg = ref(false)
const saving = ref(false)
interface AccountLine {
  product_path: Array<string | number> | string | number | undefined
  quantity: number
  unit: string
}
interface GiftWineFormState {
  gift_wine_product_path: Array<string | number> | string | number | undefined
  gift_wine_unit: string
  gift_wine_quantity: number
}
interface ConsumableLine {
  kind: 'product' | 'custom'
  consumable_product_id: number | ''
  quantity: number
  name: string
  amount: number
}
const cForm = reactive({
  channel: '',
  order_no: '',
  income_amount: undefined as number | undefined,
  member_id: 0,
  payment_status: 1,
  is_errand_order: 0,
  errand_fee: 0,
  round_amount: 0,
  is_gift_wine: 0,
  gift_wine_product_path: [] as Array<string | number>,
  gift_wine_unit: '',
  gift_wine_quantity: 1,
  lines: [] as AccountLine[],
  remark: '',
})

const consumableProductDlg = ref(false)
const consumableProductFormDlg = ref(false)
const consumableProductSaving = ref(false)
const consumableProductEditId = ref<number | null>(null)
const consumableProductForm = reactive({
  name: '',
  cost_price: 0,
  remark: '',
})

interface CustomAccountLine {
  description: string
  quantity: number
  unit: string
  price: number
  line_remark: string
}
const customForm = reactive({
  channel: '',
  order_no: '',
  income_amount: undefined as number | undefined,
  member_id: 0,
  payment_status: 1,
  is_errand_order: 0,
  errand_fee: 0,
  round_amount: 0,
  is_gift_wine: 0,
  gift_wine_product_path: [] as Array<string | number>,
  gift_wine_unit: '',
  gift_wine_quantity: 1,
  lines: [] as CustomAccountLine[],
  remark: '',
})

function resetCreateForm(channel: string): void {
  cForm.channel = channel
  cForm.order_no = ''
  cForm.income_amount = undefined
  cForm.member_id = 0
  cForm.payment_status = 1
  cForm.is_errand_order = 0
  cForm.errand_fee = 0
  cForm.round_amount = 0
  cForm.is_gift_wine = 0
  cForm.gift_wine_product_path = []
  cForm.gift_wine_unit = ''
  cForm.gift_wine_quantity = 1
  cForm.lines = [makeAccountLine()]
  cForm.remark = ''
}

function makeCustomLine(): CustomAccountLine {
  return {
    description: '',
    quantity: 1,
    unit: '',
    price: 0,
    line_remark: '',
  }
}

function resetCustomForm(channel: string): void {
  customForm.channel = channel
  customForm.order_no = ''
  customForm.income_amount = undefined
  customForm.member_id = 0
  customForm.payment_status = 1
  customForm.is_errand_order = 0
  customForm.errand_fee = 0
  customForm.round_amount = 0
  customForm.is_gift_wine = 0
  customForm.gift_wine_product_path = []
  customForm.gift_wine_unit = ''
  customForm.gift_wine_quantity = 1
  customForm.lines = [makeCustomLine()]
  customForm.remark = ''
}

function onCreateChannelChange(value: string | number | undefined): void {
  resetCreateForm(String(value ?? '').trim())
}

function onCustomChannelChange(value: string | number | undefined): void {
  resetCustomForm(String(value ?? '').trim())
}

function openCreate(): void {
  resetCreateForm('')
  createDlg.value = true
}

function openCustomCreate(): void {
  resetCustomForm('')
  customCreateDlg.value = true
}

function addCustomLine(): void {
  customForm.lines.push(makeCustomLine())
}

function removeCustomLine(idx: number): void {
  customForm.lines = customForm.lines.filter((_, i) => i !== idx)
  if (!customForm.lines.length) customForm.lines.push(makeCustomLine())
}

function makeAccountLine(): AccountLine {
  return {
    product_path: [],
    quantity: 1,
    unit: '',
  }
}

function getProductId(path: Array<string | number> | string | number | undefined): number | null {
  if (Array.isArray(path)) {
    const leaf = path[path.length - 1]
    const id = Number(leaf)
    return Number.isFinite(id) && id > 0 ? id : null
  }
  if (typeof path === 'number' || typeof path === 'string') {
    const id = Number(path)
    return Number.isFinite(id) && id > 0 ? id : null
  }
  return null
}

function productPathById(productId: number | undefined): Array<string | number> {
  const id = Number(productId || 0)
  if (!id) return []
  for (const group of productCascaderOptions.value) {
    const child = group.children?.find((item) => Number(item.value) === id)
    if (child) return [group.value, child.value]
  }
  return [id]
}

function specOptionLabel(s: ProductUnitSpec): string {
  const name = s.unit_name || s.unit_code
  const factor = Number(s.factor_to_base || 0)
  const price = Number(s.sale_price || 0)
  const parts = [name]
  if (factor > 0) parts.push(`换算${factor}`)
  if (price > 0) parts.push(`售价${formatMoney(price)}`)
  return parts.join(' / ')
}

function specOptionValue(s: ProductUnitSpec): string {
  return String(s.unit_name || s.unit_code || '').trim()
}

function lineUnitOptions(line: AccountLine): Array<{ label: string; value: string | number }> {
  const pid = getProductId(line.product_path)
  if (!pid) return []
  const specs = specsByProduct.value.get(pid) ?? []
  if (specs.length > 0) {
    return specs.map((s) => ({ label: specOptionLabel(s), value: specOptionValue(s) }))
  }
  const product = productById.value.get(pid)
  const defaultUnit = product?.unit || unitDict.value[0]?.value || '件'
  return [{ label: defaultUnit, value: defaultUnit }]
}

function giftWineUnitOptions(form: GiftWineFormState): Array<{ label: string; value: string | number }> {
  const pid = getProductId(form.gift_wine_product_path)
  if (!pid) return []
  const specs = specsByProduct.value.get(pid) ?? []
  if (specs.length > 0) {
    return specs.map((s) => {
      const label = `${specOptionLabel(s)} / 成本${formatMoney(s.cost_price)}`
      return { label, value: specOptionValue(s) }
    })
  }
  const product = productById.value.get(pid)
  const defaultUnit = product?.unit || unitDict.value[0]?.value || '件'
  return [{ label: defaultUnit, value: defaultUnit }]
}

function onGiftWineProductChange(form: GiftWineFormState): void {
  const options = giftWineUnitOptions(form)
  form.gift_wine_unit = String(options[0]?.value || '')
}

function onCreateGiftWineProductChange(): void {
  onGiftWineProductChange(cForm)
}

function onCustomGiftWineProductChange(): void {
  onGiftWineProductChange(customForm)
}

function onEditGiftWineProductChange(): void {
  onGiftWineProductChange(eForm)
}

function giftWineCostAmount(form: GiftWineFormState): number {
  const productId = getProductId(form.gift_wine_product_path)
  if (!productId) return 0
  const qty = Number(form.gift_wine_quantity || 0)
  if (qty <= 0) return 0
  return Number((qty * itemCostPrice(productId, form.gift_wine_unit)).toFixed(2))
}

function validateGiftWine(form: GiftWineFormState): boolean {
  if (!getProductId(form.gift_wine_product_path)) {
    toast.warning('请选择赠酒商品')
    return false
  }
  if (!String(form.gift_wine_unit || '').trim()) {
    toast.warning('请选择赠酒规格')
    return false
  }
  if (Number(form.gift_wine_quantity || 0) <= 0) {
    toast.warning('请填写赠酒数量')
    return false
  }
  if (giftWineCostAmount(form) < 0) {
    toast.warning('赠酒成本不能小于0')
    return false
  }
  return true
}

function onCreateProductChange(idx: number): void {
  const line = cForm.lines[idx]
  if (!line) return
  const options = lineUnitOptions(line)
  line.unit = String(options[0]?.value || '')
}

function addCreateLine(): void {
  cForm.lines.push(makeAccountLine())
}

function removeCreateLine(idx: number): void {
  cForm.lines = cForm.lines.filter((_, i) => i !== idx)
  if (!cForm.lines.length) cForm.lines.push(makeAccountLine())
}

async function submitCustomCreate(): Promise<void> {
  if (!customForm.channel.trim()) {
    toast.warning('请选择渠道')
    return
  }
  const items: Array<{
    product_id: number
    product_name: string
    quantity: number
    unit: string
    price: number
    amount: number
    remark: string
    spec: string
  }> = []
  for (const line of customForm.lines) {
    const name = line.description.trim()
    const unit = line.unit.trim()
    if (!name) {
      toast.warning('请填写每条明细的描述')
      return
    }
    if (!unit) {
      toast.warning('请填写每条明细的单位')
      return
    }
    const price = Number(line.price)
    if (!Number.isFinite(price) || price <= 0) {
      toast.warning('自定义明细须填写大于 0 的单价')
      return
    }
    const qty = Number(line.quantity)
    if (!Number.isFinite(qty) || qty <= 0) {
      toast.warning('请填写有效的数量')
      return
    }
    const amount = Number((price * qty).toFixed(2))
    items.push({
      product_id: 0,
      product_name: name,
      quantity: qty,
      unit,
      price,
      amount,
      remark: line.line_remark.trim(),
      spec: '',
    })
  }
  if (!items.length) {
    toast.warning('请至少添加一条明细')
    return
  }
  const incomeAmount = isTakeoutChannel(customForm.channel) ? customForm.income_amount : undefined
  if (incomeAmount !== undefined && (!Number.isFinite(Number(incomeAmount)) || Number(incomeAmount) < 0)) {
    toast.warning('请填写有效的收入金额')
    return
  }
  if (customForm.is_errand_order === 1 && Number(customForm.errand_fee || 0) <= 0) {
    toast.warning('跑腿订单请填写跑腿费用')
    return
  }
  if (Number(customForm.round_amount || 0) < 0) {
    toast.warning('抹零金额不能小于0')
    return
  }
  if (customForm.is_gift_wine === 1 && !validateGiftWine(customForm)) return
  const customGiftWineCostAmount = customForm.is_gift_wine === 1 ? giftWineCostAmount(customForm) : 0
  saving.value = true
  try {
    await createStoreAccount({
      store_id: tenantStoreId.value,
      member_id: customForm.member_id > 0 ? customForm.member_id : undefined,
      payment_status: customForm.payment_status,
      channel: customForm.channel.trim(),
      order_no: customForm.order_no.trim() || undefined,
      income_amount: incomeAmount,
      is_errand_order: customForm.is_errand_order,
      errand_fee: customForm.is_errand_order === 1 ? Number(customForm.errand_fee || 0) : 0,
      round_amount: Number(customForm.round_amount || 0),
      is_gift_wine: customForm.is_gift_wine,
      gift_wine_product_id: customForm.is_gift_wine === 1 ? getProductId(customForm.gift_wine_product_path) : 0,
      gift_wine_unit: customForm.is_gift_wine === 1 ? customForm.gift_wine_unit : '',
      gift_wine_quantity: customForm.is_gift_wine === 1 ? Number(customForm.gift_wine_quantity || 0) : 0,
      gift_wine_cost_amount: customGiftWineCostAmount,
      remark: customForm.remark.trim(),
      other_expense_amount: 0,
      items,
    })
    toast.success('已保存')
    customCreateDlg.value = false
    await reloadAll()
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '失败')
  } finally {
    saving.value = false
  }
}

async function submitCreate(): Promise<void> {
  if (!cForm.channel.trim()) {
    toast.warning('请选择渠道')
    return
  }
  const items = cForm.lines
    .map((line) => ({
      product_id: getProductId(line.product_path),
      quantity: line.quantity,
      unit: line.unit.trim(),
    }))
    .filter((line) => line.product_id && line.quantity > 0 && line.unit)
    .map((line) => ({
      product_id: line.product_id as number,
      quantity: line.quantity,
      unit: line.unit,
      spec: '',
      price: 0,
      amount: 0,
      remark: '',
    }))
  if (!items.length) {
    toast.warning('请至少选择一条有效商品明细')
    return
  }
  const incomeAmount = isTakeoutChannel(cForm.channel) ? cForm.income_amount : undefined
  if (incomeAmount !== undefined && (!Number.isFinite(Number(incomeAmount)) || Number(incomeAmount) < 0)) {
    toast.warning('请填写有效的收入金额')
    return
  }
  if (cForm.is_errand_order === 1 && Number(cForm.errand_fee || 0) <= 0) {
    toast.warning('跑腿订单请填写跑腿费用')
    return
  }
  if (Number(cForm.round_amount || 0) < 0) {
    toast.warning('抹零金额不能小于0')
    return
  }
  if (cForm.is_gift_wine === 1 && !validateGiftWine(cForm)) return
  const createGiftWineCostAmount = cForm.is_gift_wine === 1 ? giftWineCostAmount(cForm) : 0
  saving.value = true
  try {
    await createStoreAccount({
      store_id: tenantStoreId.value,
      member_id: cForm.member_id > 0 ? cForm.member_id : undefined,
      payment_status: cForm.payment_status,
      channel: cForm.channel.trim(),
      order_no: cForm.order_no.trim() || undefined,
      income_amount: incomeAmount,
      is_errand_order: cForm.is_errand_order,
      errand_fee: cForm.is_errand_order === 1 ? Number(cForm.errand_fee || 0) : 0,
      round_amount: Number(cForm.round_amount || 0),
      is_gift_wine: cForm.is_gift_wine,
      gift_wine_product_id: cForm.is_gift_wine === 1 ? getProductId(cForm.gift_wine_product_path) : 0,
      gift_wine_unit: cForm.is_gift_wine === 1 ? cForm.gift_wine_unit : '',
      gift_wine_quantity: cForm.is_gift_wine === 1 ? Number(cForm.gift_wine_quantity || 0) : 0,
      gift_wine_cost_amount: createGiftWineCostAmount,
      remark: cForm.remark.trim(),
      other_expense_amount: 0,
      items,
    })
    toast.success('已保存')
    createDlg.value = false
    await reloadAll()
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '失败')
  } finally {
    saving.value = false
  }
}

const editDlg = ref(false)
const editId = ref(0)
const eForm = reactive({
  member_id: 0,
  payment_status: 1,
  channel: '',
  order_no: '',
  income_amount: undefined as number | undefined,
  is_errand_order: 0,
  errand_fee: 0,
  round_amount: 0,
  is_gift_wine: 0,
  gift_wine_product_path: [] as Array<string | number>,
  gift_wine_unit: '',
  gift_wine_quantity: 1,
  tag_code: '',
  tag_name: '',
  remark: '',
})

function openEdit(row: StoreAccount): void {
  if (!canEditAccount(row)) {
    toast.warning('账单仅支持当前营业日内修改')
    return
  }
  editId.value = row.id
  eForm.member_id = Number(row.member_id || 0)
  eForm.payment_status = Number(row.payment_status || 1)
  eForm.channel = row.channel ?? ''
  eForm.order_no = row.order_no ?? ''
  eForm.income_amount = isTakeoutChannel(eForm.channel) ? Number(row.total_amount || 0) : undefined
  eForm.is_errand_order = Number(row.is_errand_order || 0) === 1 ? 1 : 0
  eForm.errand_fee = Number(row.errand_fee || 0)
  eForm.round_amount = Number(row.round_amount || 0)
  eForm.is_gift_wine = Number(row.is_gift_wine || 0) === 1 ? 1 : 0
  eForm.gift_wine_product_path = productPathById(row.gift_wine_product_id)
  eForm.gift_wine_unit = row.gift_wine_unit ?? ''
  eForm.gift_wine_quantity = Number(row.gift_wine_quantity || 1)
  eForm.tag_code = row.tag_code ?? ''
  eForm.tag_name = row.tag_name ?? ''
  eForm.remark = row.remark ?? ''
  editDlg.value = true
}

async function submitEdit(): Promise<void> {
  const incomeAmount = isTakeoutChannel(eForm.channel) ? eForm.income_amount : undefined
  if (incomeAmount !== undefined && (!Number.isFinite(Number(incomeAmount)) || Number(incomeAmount) < 0)) {
    toast.warning('请填写有效的收入金额')
    return
  }
  if (eForm.is_errand_order === 1 && Number(eForm.errand_fee || 0) <= 0) {
    toast.warning('跑腿订单请填写跑腿费用')
    return
  }
  if (Number(eForm.round_amount || 0) < 0) {
    toast.warning('抹零金额不能小于0')
    return
  }
  if (eForm.is_gift_wine === 1 && !validateGiftWine(eForm)) return
  const editGiftWineCostAmount = eForm.is_gift_wine === 1 ? giftWineCostAmount(eForm) : 0
  saving.value = true
  try {
    await updateStoreAccount(editId.value, {
      member_id: eForm.member_id,
      payment_status: eForm.payment_status,
      channel: eForm.channel.trim(),
      order_no: eForm.order_no.trim() || undefined,
      income_amount: incomeAmount,
      is_errand_order: eForm.is_errand_order,
      errand_fee: eForm.is_errand_order === 1 ? Number(eForm.errand_fee || 0) : 0,
      round_amount: Number(eForm.round_amount || 0),
      is_gift_wine: eForm.is_gift_wine,
      gift_wine_product_id: eForm.is_gift_wine === 1 ? getProductId(eForm.gift_wine_product_path) : 0,
      gift_wine_unit: eForm.is_gift_wine === 1 ? eForm.gift_wine_unit : '',
      gift_wine_quantity: eForm.is_gift_wine === 1 ? Number(eForm.gift_wine_quantity || 0) : 0,
      gift_wine_cost_amount: editGiftWineCostAmount,
      tag_code: eForm.tag_code.trim(),
      tag_name: eForm.tag_name.trim(),
      remark: eForm.remark.trim(),
    })
    toast.success('已保存')
    editDlg.value = false
    await qc.invalidateQueries({ queryKey: ['store-accounts'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '失败')
  } finally {
    saving.value = false
  }
}

async function markAccountPaid(row: StoreAccount): Promise<void> {
  const ok = await confirmDialog({ message: `确认将记账单「${row.account_no || row.id}」标记为已支付？` })
  if (!ok) return
  saving.value = true
  try {
    await updateStoreAccount(row.id, { payment_status: 1 })
    toast.success('已标记为已支付')
    await qc.invalidateQueries({ queryKey: ['store-accounts'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '支付状态修改失败')
  } finally {
    saving.value = false
  }
}

const viewDlg = ref(false)
const viewAccount = ref<StoreAccount | null>(null)
const consumableDlg = ref(false)
const consumableSaving = ref(false)
const consumableTarget = ref<StoreAccount | null>(null)
const consumableLines = ref<ConsumableLine[]>([])

function formatMoney(v: number | string | undefined | null): string {
  const n = Number(v ?? 0)
  return Number.isFinite(n) ? n.toFixed(2) : '0.00'
}

function formatDateTime(v: string): string {
  const s = String(v || '').trim()
  if (!s) return '-'
  const d = new Date(s)
  if (Number.isNaN(d.getTime())) return s.slice(0, 19).replace('T', ' ')
  const pad = (x: number) => String(x).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
}

function accountConsumableAmount(account: StoreAccount): number {
  return (account.consumables ?? []).reduce((sum, c) => sum + Number(c.amount || 0), 0)
}

function itemCostPrice(productId: number, unit?: string): number {
  const specs = specsByProduct.value.get(productId) ?? []
  const normalized = String(unit || '').trim().toLowerCase()
  for (const s of specs) {
    if (!s.is_enabled) continue
    const code = String(s.unit_code || '').trim().toLowerCase()
    const name = String(s.unit_name || '').trim().toLowerCase()
    if (normalized && (normalized === code || normalized === name)) return Number(s.cost_price || 0)
  }
  for (const s of specs) {
    if (!s.is_enabled) continue
    const code = String(s.unit_code || '').trim().toLowerCase()
    const name = String(s.unit_name || '').trim().toLowerCase()
    if (normalized && (code.includes(normalized) || name.includes(normalized) || normalized.includes(code) || normalized.includes(name))) {
      return Number(s.cost_price || 0)
    }
  }
  return 0
}

function accountItemCost(account: StoreAccount): number {
  return (account.items ?? []).reduce((sum, it) => {
    const qty = Number(it.quantity || 0)
    if (qty <= 0) return sum
    return sum + qty * itemCostPrice(it.product_id, it.unit)
  }, 0)
}

function accountNetProfitBreakdown(account: StoreAccount): number {
  return (
    Number(account.total_amount || 0) -
    Number(account.other_expense_amount || 0) -
    Number(account.errand_fee || 0) -
    accountItemCost(account) -
    accountConsumableAmount(account) -
    Number(account.gift_wine_cost_amount || 0) -
    Number(account.round_amount || 0)
  )
}

async function openView(row: StoreAccount): Promise<void> {
  viewAccount.value = null
  viewDlg.value = true
  try {
    const full = await getStoreAccount(row.id)
    viewAccount.value = full
  } catch (e: unknown) {
    viewDlg.value = false
    toast.error(e instanceof Error ? e.message : '加载失败')
  }
}

function makeConsumableLine(): ConsumableLine {
  return { kind: 'product', consumable_product_id: '', quantity: 1, name: '', amount: 0 }
}

function makeCustomConsumableLine(): ConsumableLine {
  return { kind: 'custom', consumable_product_id: '', quantity: 1, name: '', amount: 0 }
}

function addConsumableLine(): void {
  consumableLines.value.push(makeConsumableLine())
}

function addCustomConsumableLine(): void {
  consumableLines.value.push(makeCustomConsumableLine())
}

function removeConsumableLine(idx: number): void {
  consumableLines.value = consumableLines.value.filter((_, i) => i !== idx)
  if (!consumableLines.value.length) consumableLines.value.push(makeConsumableLine())
}

function selectedConsumableProduct(line: ConsumableLine): StoreAccountConsumableProduct | undefined {
  const id = Number(line.consumable_product_id || 0)
  return id > 0 ? consumableProductMap.value.get(id) : undefined
}

function consumableLineAmount(line: ConsumableLine): number {
  if (line.kind === 'custom') return Number(line.amount || 0)
  const product = selectedConsumableProduct(line)
  return Number(product?.cost_price || 0) * Number(line.quantity || 0)
}

const consumableBindTotal = computed(() =>
  consumableLines.value.reduce((sum, line) => sum + consumableLineAmount(line), 0),
)

async function openConsumableDlg(row: StoreAccount): Promise<void> {
  if (!canBindConsumables(row)) {
    toast.warning('该记账单已绑定消耗品，不能重复绑定')
    return
  }
  consumableTarget.value = row
  consumableLines.value = [makeConsumableLine()]
  try {
    const full = await getStoreAccount(row.id)
    if (full.consumables?.length) {
      toast.warning('该记账单已绑定消耗品，不能重复绑定')
      return
    }
    consumableDlg.value = true
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '加载消耗品失败')
  }
}

async function submitConsumables(): Promise<void> {
  if (!consumableTarget.value) return
  const consumables: Array<Record<string, unknown>> = []
  for (const line of consumableLines.value) {
    if (line.kind === 'custom') {
      const name = line.name.trim()
      const amount = Number(line.amount || 0)
      if (!name) {
        toast.warning('请填写自定义消耗品名称')
        return
      }
      if (!Number.isFinite(amount) || amount <= 0) {
        toast.warning(`请填写「${name}」的有效金额`)
        return
      }
      consumables.push({
        product_id: 0,
        product_name: name,
        quantity: 1,
        amount,
      })
      continue
    }
    const productID = Number(line.consumable_product_id || 0)
    const quantity = Number(line.quantity || 0)
    if (productID > 0 && quantity > 0) {
      consumables.push({
        consumable_product_id: productID,
        quantity,
      })
    }
  }
  if (!consumables.length) {
    toast.warning('请至少选择一条有效消耗品明细')
    return
  }
  consumableSaving.value = true
  try {
    await bindStoreAccountConsumables(consumableTarget.value.id, { consumables })
    toast.success('消耗品已绑定')
    consumableDlg.value = false
    await qc.invalidateQueries({ queryKey: ['store-accounts'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '保存失败')
  } finally {
    consumableSaving.value = false
  }
}

function reloadConsumableProducts(): void {
  consumableProductPage.value = 1
  void qc.invalidateQueries({ queryKey: ['store-account-consumable-products'] })
  void qc.invalidateQueries({ queryKey: ['store-account-consumable-products-all'] })
}

function openConsumableProductManage(): void {
  consumableProductDlg.value = true
  reloadConsumableProducts()
}

function openConsumableProductCreate(): void {
  consumableProductEditId.value = null
  consumableProductForm.name = ''
  consumableProductForm.cost_price = 0
  consumableProductForm.remark = ''
  consumableProductFormDlg.value = true
}

function openConsumableProductEdit(row: StoreAccountConsumableProduct): void {
  consumableProductEditId.value = row.id
  consumableProductForm.name = row.name || ''
  consumableProductForm.cost_price = Number(row.cost_price || 0)
  consumableProductForm.remark = row.remark || ''
  consumableProductFormDlg.value = true
}

async function submitConsumableProduct(): Promise<void> {
  const name = consumableProductForm.name.trim()
  if (!name) {
    toast.warning('请填写消耗品名称')
    return
  }
  const body = {
    store_id: tenantStoreId.value,
    name,
    cost_price: Number(consumableProductForm.cost_price || 0),
    remark: consumableProductForm.remark.trim(),
  }
  consumableProductSaving.value = true
  try {
    if (consumableProductEditId.value) {
      await updateStoreAccountConsumableProduct(consumableProductEditId.value, body)
    } else {
      await createStoreAccountConsumableProduct(body)
    }
    toast.success('已保存')
    consumableProductFormDlg.value = false
    await qc.invalidateQueries({ queryKey: ['store-account-consumable-products'] })
    await qc.invalidateQueries({ queryKey: ['store-account-consumable-products-all'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '保存失败')
  } finally {
    consumableProductSaving.value = false
  }
}

async function deleteConsumableProduct(row: StoreAccountConsumableProduct): Promise<void> {
  const ok = await confirmDialog({ message: `删除消耗品「${row.name}」？` })
  if (!ok) return
  try {
    await deleteStoreAccountConsumableProduct(row.id)
    toast.success('已删除')
    await qc.invalidateQueries({ queryKey: ['store-account-consumable-products'] })
    await qc.invalidateQueries({ queryKey: ['store-account-consumable-products-all'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '删除失败')
  }
}

function consumableProductActions(row: StoreAccountConsumableProduct): TableRowAction[] {
  return [
    { label: '编辑', permission: 'store:account:edit', onClick: () => openConsumableProductEdit(row) },
    { label: '删除', permission: 'store:account:edit', danger: true, onClick: () => void deleteConsumableProduct(row) },
  ]
}
</script>

<style scoped>
.store-account-channel {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  min-width: 0;
}

.store-account-channel__icon {
  display: inline-flex;
  width: 22px;
  height: 22px;
  flex: 0 0 22px;
  align-items: center;
  justify-content: center;
  border-radius: 999px;
  color: #fff;
  font-size: 12px;
  font-weight: 800;
  line-height: 1;
  box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.28);
}

.store-account-channel__label {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.store-account-channel__icon--wechat {
  background: linear-gradient(135deg, #19c37d, #10a966);
}

.store-account-channel__icon--meituan {
  background: linear-gradient(135deg, #ffd43b, #f59f00);
  color: #533300;
}

.store-account-channel__icon--eleme {
  background: linear-gradient(135deg, #38bdf8, #2563eb);
}

.store-account-channel__icon--douyin {
  background: linear-gradient(135deg, #111827, #ef476f);
}

.store-account-channel__icon--taobao {
  background: linear-gradient(135deg, #ff7a1a, #f97316);
}

.store-account-channel__icon--redbook {
  background: linear-gradient(135deg, #ff4d6d, #dc2626);
}

.store-account-channel__icon--offline {
  background: linear-gradient(135deg, #64748b, #334155);
}

.store-account-channel__icon--other {
  background: linear-gradient(135deg, #94a3b8, #64748b);
}

.store-account-pay-status {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  white-space: nowrap;
}

.store-account-pay-status__dot {
  width: 7px;
  height: 7px;
  flex: 0 0 7px;
  border-radius: 999px;
  box-shadow: 0 0 0 3px rgba(148, 163, 184, 0.14);
}

.store-account-pay-status__dot--paid {
  background: #22c55e;
  box-shadow: 0 0 0 3px rgba(34, 197, 94, 0.14);
}

.store-account-pay-status__dot--unpaid {
  background: #f59e0b;
  box-shadow: 0 0 0 3px rgba(245, 158, 11, 0.16);
}
</style>
