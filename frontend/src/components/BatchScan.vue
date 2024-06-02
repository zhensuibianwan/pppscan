<template >
  <el-form    label-position="top" class="background">
    <el-form-item style="margin-left: 3%; margin-top: 3%">

<!--      <el-upload-->
<!--          :auto-upload="false"-->
<!--          :on-change="uploadChange"-->
<!--          :show-file-list="false"-->
<!--          style="margin-left: 3%"-->
<!--      >-->
<!--        <el-button type="primary" size="large" :icon="Upload">导入</el-button>-->
<!--      </el-upload>-->
      <el-button  size="large" :icon="Setting" @click="showBatchScanSetting=true" type="success" >设置</el-button>
      <el-button  size="large" :icon="Download" @click="exportToExcel" style="margin-left: 3%" >导出</el-button>
      <el-button  size="large" :icon="RefreshRight" @click="clearInput" type="danger" style="margin-left: 3%" >清空</el-button>
      <el-button  size="large" :icon="Aim" @click="batchScan" type="warning" style="margin-left: 57%" v-show="scanShow">扫描</el-button>
      <el-button  size="large" :icon="VideoPause" @click="pauseScan" type="warning" style="margin-left: 46.5%" v-show="scanPauseShow">暂停</el-button>
      <el-button  size="large" :icon="CaretRight" @click="continueScan" type="warning" style="margin-left: 46.5%" v-show="scanContinueShow" >继续</el-button>
      <el-button  size="large" :icon="Close" @click="closeScan" type="danger" style="margin-left: 2%" v-show="scanCloseShow">取消</el-button>
    </el-form-item>
<!--    <el-form-item label="URL"  class="labelSize">-->
<!--      <el-input v-model="BatchScanForm.urls" style="width: 75%" :rows="5" type="textarea"/>-->
<!--      <el-button  size="large" :icon="Aim" @click="batchScan" type="warning" style="margin-left: 5%" v-show="scanShow">扫描</el-button>-->
<!--      <dev>-->
<!--        <el-button  size="large" :icon="VideoPause" @click="pauseScan" type="warning" style="margin-left: 5%" v-show="scanPauseShow">暂停</el-button>-->
<!--        <el-button  size="large" :icon="CaretRight" @click="continueScan" type="warning" style="margin-left: 5%" v-show="scanContinueShow" >继续</el-button>-->
<!--        <el-button  size="large" :icon="Close" @click="closeScan" type="danger" style="margin-left: 5%; margin-top: 10%" v-show="scanCloseShow">取消</el-button>-->
<!--      </dev>-->
<!--    </el-form-item>-->
<!--    <el-form-item label="扫描结果"  class="labelSize"  style="margin-top: 5%">-->
<!--      <el-input v-model="BatchScanForm.result" style="width: 75%;border: 2px solid var(&#45;&#45;el-border-color);border-radius: 4px;" :rows="10" type="textarea"/>-->
<!--    </el-form-item>-->
  </el-form>

  <div class="batchScanResultClass">
    <el-table
        ref="multipleTableRef"
        :data="result.data"
        style="width: 100%"
        @selection-change="handleResultSelectionChange"
        height="100%"
    >
      <el-table-column   type="selection" width="35" />
      <el-table-column   label="URL" width="250">
        <template #default="scope">
          <span @dblclick="openURL(scope.row.url)">{{ scope.row.url }}</span>
        </template>
      </el-table-column>
      <el-table-column property="pocName" label="漏洞名称" width="400" />
      <el-table-column property="print" label="结果输出" width="330"/>
    </el-table>
    <el-progress :text-inside="true" :stroke-width="30" :percentage="progress"  :color="colors"/>
  </div>

  <div style="margin-left: 10%; margin-top: 3%">
    <el-pagination
        v-model:current-page="resultTablePage"
        v-model:page-size="resultTablePageSize"
        :page-sizes="[10, 20, 50, 100]"
        :small="false"
        :background="true"
        layout="prev, pager, next, jumper,sizes,total"
        :total="result.dataAll.length"
        @size-change="resultTableHandleSizeChange"
        @current-change="resultTableHandleCurrentChange"
    />
  </div>
  <el-dialog
      v-model="showBatchScanSetting"
      align-center
      center
      width="85%"
      draggable
      title="扫描设置"
  >

    <el-row>
      <el-form-item label="URL"  class="labelSize" style="margin-left: 5%">
        <el-input v-model="BatchScanForm.urls" style="width: 400px" :rows="5" type="textarea"/>
      </el-form-item>
      <el-upload
          :auto-upload="false"
          :on-change="uploadChange"
          :show-file-list="false"
          style="margin-left: 3%"
      >
        <el-button type="primary" size="large" :icon="Upload">导入</el-button>
      </el-upload>
      <el-button  size="large" :icon="Star" @click="showAdd=true" type="success" style="margin-left: 3%">选择poc</el-button>
    </el-row>

    <div style="margin-left: 63.5%;margin-top: 2%">
      <el-button  size="large" :icon="CloseBold" @click="BatchScanSettingCancelORSure">取消</el-button>
      <el-button size="large" :icon="Select"  type="primary" @click="BatchScanSettingCancelORSure">确认</el-button>
    </div>
  </el-dialog>

  <el-dialog
      v-model="showAdd"
      align-center
      center
      width="80%"
      draggable
      title="选择POC"
  >
    <el-col :span="6" style="margin-left: 3%;">
      <el-input v-model="input"  style="width: 378%; " size="large"  @keyup.enter="pocList()" clearable>
        <template #prepend>漏洞名称</template>
        <template #append><el-button :icon="Search" style="width: 75px;color: #409eff"  @click="pocList"/></template>
      </el-input>
    </el-col>
    <div class="demo-radius2" >
      <el-table
          ref="batchScanPocListRef"
          :data="table.tableData"
          style="width: 100%"
          @selection-change="BatchHandleSelectionChange"
          height="100%"
      >
        <el-table-column   type="selection" width="35" />
        <el-table-column   label="漏洞名称" width="430">
          <template #default="scope">
            <el-tooltip  placement="top" >
              <template #content>  {{ scope.row.description }} </template>
              {{ scope.row.name }}
            </el-tooltip>
          </template>
        </el-table-column>
      </el-table>
    </div>
<!--    <el-button size="large" :icon="Select"  type="primary" @click="" style="margin-top: 2%; margin-left: 86%">确认</el-button>-->
    <div style="margin-left: 72%;margin-top: 2%">
      <el-button  size="large" :icon="CloseBold" @click="cancelSelectPoc">取消</el-button>
      <el-button size="large" :icon="Select"  type="primary" @click="sureSelectPoc">确认</el-button>
    </div>
  </el-dialog>

</template>

<script lang="ts" setup>

import {reactive, ref,computed} from "vue";
import {publicCode} from "../../wailsjs/go/models";
import PocScanResult = publicCode.PocScanResult;
import {Search, Aim,Upload,RefreshRight,Star,Download,Close,CaretRight,Select,CloseBold,VideoPause,Setting} from '@element-plus/icons-vue';
import {
  DefaultBrowserOpenUrl,
  LocalList,
  PocList,
  ScanBatch,
  ScanClose,
  ScanContinue,
  ScanPause
} from "../../wailsjs/go/main/App";
import {ElMessage, ElMessageBox, ElTable} from "element-plus";
import {EventsOn} from "../../wailsjs/runtime";

const BatchScanForm = reactive({
  urls: "",
  pocData: [] as Array<publicCode.Poc>,
})

const showAdd = ref(false)
const showBatchScanSetting = ref(false)
const input = ref('')

const colors = [
  { color: '#f56c6c', percentage: 20 },
  { color: '#5cb87a', percentage: 40 },
  { color: '#e6a23c', percentage: 60 },
  { color: '#67c23a', percentage: 80 },
  { color: '#1989fa', percentage: 100 },
]

const scanShow = ref(true)
const scanPauseShow = ref(false)
const scanContinueShow = ref(false)
const scanCloseShow = ref(false)


const pocSelection = ref<publicCode.Poc[]>([])

const batchScanPocListRef = ref<InstanceType<typeof ElTable>>()

function sureSelectPoc() {
  showAdd.value = false
}

function cancelSelectPoc() {
  batchScanPocListRef.value!.clearSelection()
  showAdd.value = false
}

const BatchHandleSelectionChange = (val: publicCode.Poc[]) => {

  pocSelection.value = val
}


const table = reactive({
  tableData: [] as Array<publicCode.Poc>,
})

async function pocList() {
  table.tableData = await PocList(input.value)
}

async function pocListFile() {
  table.tableData = await LocalList(input.value)
}


function loadingList(){
  pocListFile()
  pocList()
}


loadingList()




function batchScan(){
  if (pocSelection.value.length > 0 && BatchScanForm.urls != ""){
    //BatchScanForm.result = ""
    ScanBatch( BatchScanForm.urls,pocSelection.value)
    scanShow.value = false
    scanPauseShow.value = true
    scanCloseShow.value = true
  }else {
    ElMessage.error("请选择POC")
  }
}


function pauseScan(){

  scanPauseShow.value = false
  scanContinueShow.value = true
  ScanPause()

}

function  continueScan(){
  scanContinueShow.value = false
  scanPauseShow.value = true
  ScanContinue()
}


function closeScan(){
  scanContinueShow.value = false
  scanCloseShow.value = false
  scanPauseShow.value = false
  scanShow.value = true
  ScanClose()
}


const file = ref('')


const uploadChange = (file) => {
  if (file.raw) {
    const reader = new FileReader()
    reader.onload = () => {
      BatchScanForm.urls = reader.result.toString()
    }
    reader.readAsText(file.raw)
  }
}


function clearInput(){
  ElMessageBox.confirm(
      '确定清空吗？',
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
  )
      .then(() => {
        BatchScanForm.urls = ""
        result.data = [] as Array<PocScanResult>
        result.dataAll = []as Array<PocScanResult>
        cancelSelectPoc()
        //BatchScanForm.result = ""
        ElMessage({
          type: 'success',
          message: '清空成功！',
        })
        // setTimeout(function() {
        //   location.reload();
        // }, 800);
      })
      .catch(() => {
      })
}

const result = reactive({
   data: [] as Array<PocScanResult>,
  dataAll: [] as Array<PocScanResult>,
})


EventsOn('BatchScan', function (msg){
  result.dataAll.push(msg)
  result.data = JSON.parse(JSON.stringify(paginatedData.value))

})


const progress = ref(0)

EventsOn('BatchScanProgress', function (msg){
  progress.value = msg
  if (progress.value == 100) {
    scanShow.value = true
    scanPauseShow.value = false
    scanContinueShow.value = false
    scanCloseShow.value = false
  }
})



function exportToExcel() {
  let csv = 'URL,漏洞名称,结果输出\n'
  result.dataAll.forEach((item) =>{
    csv = csv + item.url + "," + `${item.pocName}` + "," + `${item.print}` + "\n"
  })
  let csvContent = 'data:text/txt;charset=utf-8,' + encodeURIComponent(csv);
  let link = document.createElement('a');
  link.setAttribute('href', csvContent);
  let currentTime = new Date().toLocaleString('sv', {year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute:'2-digit'}).replace(/[-\s:]/g, '');
  link.setAttribute('download', currentTime + '_result' +  '.csv');
  // link.setAttribute('download', 'filename.csv'); // 可以自定义下载的文件名
  document.body.appendChild(link);
  link.click();

}


function BatchScanSettingCancelORSure(){
  showBatchScanSetting.value = false
}

const resultTablePage = ref(1)
const resultTablePageSize = ref(10)

const paginatedData = computed(() => {
  const startIndex = (resultTablePage.value - 1) * resultTablePageSize.value
  const endIndex = startIndex + resultTablePageSize.value
  return result.dataAll.slice(startIndex, endIndex)
})

// 分页处理函数，切换分页时更新 result.data
function resultTableHandleCurrentChange() {
  // resultTablePage.value = 1  // 切换分页时重置当前页码为 1
  result.data = JSON.parse(JSON.stringify(paginatedData.value))
}

// 分页处理函数，修改每页显示数量时更新 result.data
function resultTableHandleSizeChange(){
  // resultTablePage.value = 1  // 修改每页显示数量时重置当前页码为 1
  result.data = JSON.parse(JSON.stringify(paginatedData.value))
}

function  openURL(url) {
  DefaultBrowserOpenUrl(url);
}

</script>

<style>


.demo-radius2 {
  height: 300PX;
  width: 94%;
  border: 2px solid var(--el-border-color);
  border-radius: 4px;
  margin-top: 2%;
  margin-left: 3%;
}

.background{
  background-size: cover;
  background-repeat: no-repeat;
  background-position: left center;
  opacity: 0.8;
}

.batchScanResultClass {
  height: 65%;
  width: 94%;
  border: 2px solid var(--el-border-color);
  border-radius: 4px;
  margin-top: 2%;
  margin-left: 3%;
  background-size: cover;
  background-repeat: no-repeat;
  background-position: left center;
  opacity: 0.8;
}

</style>