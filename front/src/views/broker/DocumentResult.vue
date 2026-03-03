<template>
  <div>
    <!-- Header -->
    <div style="display: flex; align-items: center; gap: 12px; margin-bottom: 24px">
      <el-button text @click="router.push('/broker/documents')">
        <el-icon><ArrowLeft /></el-icon> Назад
      </el-button>
      <h2 style="margin: 0">{{ docName }}</h2>
    </div>

    <div v-if="loadingData" style="text-align: center; padding: 60px">
      <el-icon :size="40" class="is-loading"><Loading /></el-icon>
      <p style="color: #909399; margin-top: 16px">Загрузка данных...</p>
    </div>

    <div v-else-if="error" style="text-align: center; padding: 60px">
      <el-result icon="error" title="Ошибка" :sub-title="error">
        <template #extra>
          <el-button type="primary" @click="loadData">Повторить</el-button>
        </template>
      </el-result>
    </div>

    <template v-else-if="aiData">
      <!-- Export toolbar -->
      <el-card style="margin-bottom: 16px">
        <div style="display: flex; align-items: center; justify-content: space-between; flex-wrap: wrap; gap: 12px">
          <span style="color: #606266; font-size: 14px">
            Тип документа: <strong>{{ docTypeLabel }}</strong>
          </span>
          <div style="display: flex; gap: 8px; flex-wrap: wrap">
            <el-button type="warning" @click="printPDF">
              <el-icon><Printer /></el-icon> PDF
            </el-button>
            <el-button type="success" @click="downloadExcel">
              <el-icon><Download /></el-icon> Excel
            </el-button>
            <el-button type="info" @click="downloadXML">
              <el-icon><Document /></el-icon> XML
            </el-button>
            <el-button type="primary" :loading="saving" @click="saveChanges">
              <el-icon><Check /></el-icon> Сохранить изменения
            </el-button>
          </div>
        </div>
      </el-card>

      <!-- Tabs -->
      <el-tabs v-model="activeTab" type="border-card">
        <!-- Tab 1: Main fields -->
        <el-tab-pane label="Основные данные" name="fields">
          <el-table :data="fieldRows" style="width: 100%" border>
            <el-table-column label="Поле" width="220">
              <template #default="{ row }">
                <span style="font-weight: 500">{{ row.label }}</span>
              </template>
            </el-table-column>
            <el-table-column label="Значение">
              <template #default="{ row }">
                <el-input v-model="row.value" size="small" />
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <!-- Tab 2: Items -->
        <el-tab-pane label="Товары" name="items">
          <div v-if="aiData.items && aiData.items.length > 0">
            <el-table :data="aiData.items" style="width: 100%" border>
              <el-table-column
                v-for="col in itemColumns"
                :key="col.key"
                :label="col.label"
                :min-width="col.width"
              >
                <template #default="{ row }">
                  <el-input v-model="row[col.key]" size="small" />
                </template>
              </el-table-column>
            </el-table>
            <div style="margin-top: 16px">
              <el-button type="primary" plain size="small" @click="addItem">
                <el-icon><Plus /></el-icon> Добавить строку
              </el-button>
            </div>
          </div>
          <el-empty v-else description="Нет товарных позиций">
            <el-button type="primary" @click="addItem">Добавить строку</el-button>
          </el-empty>
        </el-tab-pane>
      </el-tabs>
    </template>

    <!-- Hidden print area -->
    <div id="print-area" style="display: none">
      <div class="print-doc">
        <h2 style="text-align: center; margin-bottom: 16px">
          Таможенная декларация — {{ docName }}
        </h2>
        <h3>Основные данные</h3>
        <table class="print-table">
          <tbody>
            <tr v-for="row in fieldRows" :key="row.key">
              <td class="label-cell">{{ row.label }}</td>
              <td>{{ row.value }}</td>
            </tr>
          </tbody>
        </table>
        <template v-if="aiData && aiData.items && aiData.items.length">
          <h3 style="margin-top: 24px">Товары</h3>
          <table class="print-table">
            <thead>
              <tr>
                <th v-for="col in itemColumns" :key="col.key">{{ col.label }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(item, idx) in aiData.items" :key="idx">
                <td v-for="col in itemColumns" :key="col.key">{{ item[col.key] ?? '' }}</td>
              </tr>
            </tbody>
          </table>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useDocumentsStore } from '@/stores/documents'
import { ElMessage } from 'element-plus'
import {
  ArrowLeft, Loading, Printer, Download, Document,
  Check, Plus
} from '@element-plus/icons-vue'
import type { AIData } from '@/types'

const route = useRoute()
const router = useRouter()
const docStore = useDocumentsStore()

const docId = Number(route.params.id)
const loadingData = ref(false)
const saving = ref(false)
const error = ref('')
const activeTab = ref('fields')
const aiData = ref<AIData | null>(null)

// For easy display of doc name
const docName = computed(() => {
  const doc = docStore.documents.find(d => d.id === docId)
  return doc?.original_name || `Документ #${docId}`
})

const docTypeLabel = computed(() => {
  const map: Record<string, string> = {
    customs_declaration: 'Таможенная декларация',
    invoice: 'Инвойс',
    packing_list: 'Упаковочный лист',
  }
  return map[aiData.value?.document_type || ''] || (aiData.value?.document_type || '—')
})

const FIELD_LABELS: Record<string, string> = {
  declaration_number: 'Номер декларации',
  date: 'Дата',
  sender: 'Отправитель',
  receiver: 'Получатель',
  country_origin: 'Страна происхождения',
  country_dest: 'Страна назначения',
  currency: 'Валюта',
  total_value: 'Общая стоимость',
  customs_value: 'Таможенная стоимость',
}

const ITEM_COLUMNS = [
  { key: 'number', label: '№', width: '60' },
  { key: 'hs_code', label: 'Код ТН ВЭД', width: '130' },
  { key: 'description', label: 'Описание', width: '200' },
  { key: 'quantity', label: 'Кол-во', width: '90' },
  { key: 'unit', label: 'Ед.изм.', width: '90' },
  { key: 'weight_net', label: 'Вес нетто', width: '110' },
  { key: 'weight_gross', label: 'Вес брутто', width: '110' },
  { key: 'value', label: 'Стоимость', width: '110' },
  { key: 'duty_rate', label: 'Пошлина', width: '100' },
  { key: 'vat_rate', label: 'НДС', width: '90' },
]

// Build editable field rows from aiData.fields
interface FieldRow { key: string; label: string; value: string }
const fieldRows = ref<FieldRow[]>([])

const itemColumns = computed(() => {
  if (!aiData.value?.items?.length) return ITEM_COLUMNS
  // Detect columns present in data
  const keys = new Set<string>()
  aiData.value.items.forEach(item => Object.keys(item).forEach(k => keys.add(k)))
  const ordered = ITEM_COLUMNS.filter(c => keys.has(c.key))
  const extra = [...keys].filter(k => !ITEM_COLUMNS.find(c => c.key === k))
  return [...ordered, ...extra.map(k => ({ key: k, label: k, width: '120' }))]
})

async function loadData() {
  loadingData.value = true
  error.value = ''
  try {
    // Ensure documents list is loaded for doc name
    if (!docStore.documents.length) {
      await docStore.fetchDocuments().catch(() => {})
    }
    aiData.value = await docStore.fetchAIData(docId)
    buildFieldRows()
  } catch (err: any) {
    const msg = err.response?.data?.error || ''
    const errorMap: Record<string, string> = {
      'AI data is not ready yet': 'Данные AI ещё не готовы',
      'document not found': 'Документ не найден',
      'access denied': 'Нет доступа',
    }
    error.value = errorMap[msg] || 'Не удалось загрузить данные'
  } finally {
    loadingData.value = false
  }
}

function buildFieldRows() {
  if (!aiData.value) return
  const fields = aiData.value.fields || {}
  // Show known fields in order, then the rest
  const ordered = Object.keys(FIELD_LABELS).filter(k => k in fields)
  const extra = Object.keys(fields).filter(k => !(k in FIELD_LABELS))
  fieldRows.value = [...ordered, ...extra].map(key => ({
    key,
    label: FIELD_LABELS[key] || key,
    value: String(fields[key] ?? ''),
  }))
}

function collectAIData(): AIData {
  const fields: Record<string, string | number> = {}
  fieldRows.value.forEach(row => {
    fields[row.key] = row.value
  })
  return {
    document_type: aiData.value?.document_type || '',
    fields,
    items: aiData.value?.items || [],
  }
}

async function saveChanges() {
  saving.value = true
  try {
    await docStore.updateAIData(docId, collectAIData())
    // Sync local aiData
    aiData.value = collectAIData()
    ElMessage.success('Изменения сохранены')
  } catch {
    ElMessage.error('Ошибка сохранения')
  } finally {
    saving.value = false
  }
}

async function downloadExcel() {
  try {
    const doc = docStore.documents.find(d => d.id === docId)
    await docStore.downloadExcel(docId, doc?.original_name || 'document.pdf')
    ElMessage.success('Excel скачивается')
  } catch {
    ElMessage.error('Ошибка скачивания Excel')
  }
}

async function downloadXML() {
  try {
    const doc = docStore.documents.find(d => d.id === docId)
    await docStore.downloadXML(docId, doc?.original_name || 'document.pdf')
    ElMessage.success('XML скачивается')
  } catch {
    ElMessage.error('Ошибка скачивания XML')
  }
}

function printPDF() {
  const printArea = document.getElementById('print-area')
  if (!printArea) return
  printArea.style.display = 'block'
  window.print()
  printArea.style.display = 'none'
}

function addItem() {
  if (!aiData.value) return
  const newItem: Record<string, string | number> = {}
  itemColumns.value.forEach(col => { newItem[col.key] = '' })
  if (!aiData.value.items) aiData.value.items = []
  aiData.value.items.push(newItem)
}

onMounted(loadData)
</script>

<style scoped>
#print-area {
  position: fixed;
  top: 0; left: 0;
  width: 100%;
  background: #fff;
  z-index: 9999;
  padding: 24px;
  box-sizing: border-box;
}

@media print {
  body > * { display: none !important; }
  #print-area { display: block !important; position: static; }

  .print-table {
    width: 100%;
    border-collapse: collapse;
    font-size: 12px;
    margin-bottom: 12px;
  }
  .print-table th,
  .print-table td {
    border: 1px solid #000;
    padding: 4px 8px;
    text-align: left;
  }
  .print-table th { background: #dce6f1; font-weight: bold; }
  .label-cell { font-weight: 500; width: 200px; }
}
</style>
