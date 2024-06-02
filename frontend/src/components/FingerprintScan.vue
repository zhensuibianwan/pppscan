<template>
  <el-row style="margin-top: 2%;margin-left: 3%">
      <el-form-item label="URL" style="font-size: larger;">
        <el-input v-model="searchInput.url"  @keyup.enter="searchFromResult()" @change="checkIsUpdateResultDate" clearable/>
      </el-form-item>
      <el-form-item label="标题" class="paneLabelSize" style="margin-left: 3%">
        <el-input v-model="searchInput.title"  @keyup.enter="searchFromResult()" @change="checkIsUpdateResultDate" clearable/>
      </el-form-item>
    <el-form-item label="指纹" class="paneLabelSize" style="margin-left: 3%">
      <el-input v-model="searchInput.fingerprint"  @keyup.enter="searchFromResult()" @change="checkIsUpdateResultDate" clearable/>
    </el-form-item>
    <el-button   :icon="Search" @click="searchFromResult" type="primary"  style="margin-left: 5%">搜索</el-button>
  </el-row>

  <el-row style="margin-top: 2%;margin-left: 3%">
    <el-button  size="large" :icon="Star" @click="scanSettingShow=true" type="success" >设置</el-button>
    <el-button  size="large" :icon="Download" @click="exportToExcel" style="margin-left: 3%" >导出</el-button>
    <el-button  size="large" :icon="RefreshRight" @click="clearData" type="danger" style="margin-left: 3%" >清空</el-button>
    <el-button  size="large" :icon="Aim" @click="fingerprintScan" type="warning" style="margin-left: 54.5%" v-show="scanShow">扫描</el-button>
    <el-button  size="large" :icon="VideoPause" @click="fingerprintScanPause" type="warning" style="margin-left: 42.5%" v-show="scanPauseShow">暂停</el-button>
    <el-button  size="large" :icon="CaretRight" @click="fingerprintScanContinue" type="warning" style="margin-left: 42.5%" v-show="scanContinueShow" >继续</el-button>
    <el-button  size="large" :icon="Close" @click="fingerprintScanClose" type="danger" style="margin-left: 3%" v-show="scanCloseShow">取消</el-button>
  </el-row>
  <div class="fingerprintScanResultClass">
    <el-table
        ref="fingerprintTableRef"
        :data="result.data"
        style="width: 100%"
        height="100%"
        @sort-change='resultTableSort'
    >
      <el-table-column   label="URL" width="250" fixed >
        <template #default="scope">
<!--          <a :href="scope.row.url" target="_blank">{{ scope.row.url }}</a>-->
          <span @dblclick="openURL(scope.row.url)">{{ scope.row.url }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="title" label="标题"  width="300" />
      <el-table-column prop="fingerprint" label="指纹"  sortable='custom'  width="200">
        <template #default="scope">
             <span v-for="key in scope.row.fingerprint" :key="key">
                {{ key }}
               <br>
            </span>
        </template>
      </el-table-column>
      <el-table-column label="存在的漏洞"  sortable='custom'  width="300" prop="vulnerability">
        <template #default="scope">
             <span v-for="key in scope.row.vulnerability" :key="key">
                {{ key }}
               <br>
            </span>
        </template>
      </el-table-column>
    </el-table>
    <el-progress :text-inside="true" :stroke-width="30" :percentage="progress"  :color="colors" style="width: 100%; "/>
  </div>
  <div style="margin-left: 15%;margin-top: 3%">
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
      v-model="scanSettingShow"
      align-center
      center
      width="85%"
      draggable
      title="指纹扫描设置"
  >

    <el-row>
      <el-form-item label="URL"  class="labelSize" style="margin-left: 5%">
        <el-input v-model="url" style="width: 400px" :rows="5" type="textarea"/>
      </el-form-item>
      <el-upload
          :auto-upload="false"
          :on-change="uploadChange"
          :show-file-list="false"
          style="margin-left: 3%"
      >
        <el-button type="primary" size="large" :icon="Upload">导入</el-button>
      </el-upload>
    </el-row>

    <el-form-item label="联动POC" style="margin-left: 5%">
      <el-switch v-model="linkPoc" />
    </el-form-item>
    <el-form-item label="输出没有指纹的网站" style="margin-left: 5%">
      <el-switch v-model="isEasyOutput" />
    </el-form-item>
    <div style="margin-left: 68%;margin-top: 2%">
      <el-button  size="large" :icon="CloseBold" @click="scanSettingShowCancel">取消</el-button>
      <el-button size="large" :icon="Select"  type="primary" @click="scanSettingShowSure">确认</el-button>
    </div>
  </el-dialog>

</template>

<script lang="ts" setup>
import {Search, Aim,Upload,RefreshRight,Star,Download,Close,CaretRight,Select,CloseBold,VideoPause} from '@element-plus/icons-vue';
import {computed, onMounted, reactive, ref, watch} from "vue";
import {publicCode} from "../../wailsjs/go/models";
import FingerprintScanResult = publicCode.FingerprintScanResult;
import {EventsOn} from "../../wailsjs/runtime";
import {
  DefaultBrowserOpenUrl,
  FingerprintList,
  FSScan,
  FSScanClose,
  FSScanContinue,
  FSScanPause
} from "../../wailsjs/go/main/App";
import {ElMessage, ElMessageBox} from "element-plus";

const searchInput = reactive({
  url: '',
  title: '',
  fingerprint: '',
})

const result = reactive({
  data: [] as Array<FingerprintScanResult>,
  dataAll: [] as Array<FingerprintScanResult>,
})

const url = ref('')
const linkPoc = ref(true)
const isEasyOutput = ref(true)

function scanSettingShowCancel(){
  scanSettingShow.value = false
}

function scanSettingShowSure(){
  scanSettingShow.value = false
}


const colors = [
  { color: '#f56c6c', percentage: 20 },
  { color: '#5cb87a', percentage: 40 },
  { color: '#e6a23c', percentage: 60 },
  { color: '#67c23a', percentage: 80 },
  { color: '#1989fa', percentage: 100 },
]

const progress = ref(0)

const scanSettingShow = ref(false)

const scanShow = ref(true)
const scanPauseShow = ref(false)
const scanContinueShow = ref(false)
const scanCloseShow = ref(false)



const isUpdateResultDate = ref(true)

EventsOn('FingerprintScan', function (msg){
  result.dataAll.push(msg)
  if(isUpdateResultDate.value){
    result.data = JSON.parse(JSON.stringify(paginatedData.value))
  }
})

function checkIsUpdateResultDate(){
  isUpdateResultDate.value = !(searchInput.url != '' || searchInput.title != '' || searchInput.fingerprint != '');
}

function searchFromResult(){
  if(searchInput.url != '' || searchInput.title != '' || searchInput.fingerprint != ''){
    result.data = [] as Array<FingerprintScanResult>
    result.dataAll.forEach((item) =>{
      if(item.fingerprint == null){
        item.fingerprint = [] as Array<string>
      }
      if(item.url.includes(searchInput.url) && item.title.includes(searchInput.title) && item.fingerprint.toString().includes(searchInput.fingerprint)){
        result.data.push(item)
      }
    })
  }
}


EventsOn('FingerprintScanProgress', function (msg){
  progress.value = msg
  if (progress.value == 100) {
    scanShow.value = true
    scanPauseShow.value = false
    scanContinueShow.value = false
    scanCloseShow.value = false
  }
})

function fingerprintScan(){
  if(url.value != ''){
    progress.value = 0
    result.dataAll = [] as Array<FingerprintScanResult>
    result.data = [] as Array<FingerprintScanResult>
    FSScan(url.value,linkPoc.value,isEasyOutput.value)
    scanShow.value = false
    scanPauseShow.value = true
    scanCloseShow.value = true
  }else {
    ElMessage.error('url不能为空！')
  }

}

function fingerprintScanPause(){
  FSScanPause()
  scanPauseShow.value = false
  scanContinueShow.value = true
}

function fingerprintScanContinue(){
  FSScanContinue()
  scanPauseShow.value = true
  scanContinueShow.value = false
}

function fingerprintScanClose(){
  FSScanClose()
  scanShow.value = true
  scanPauseShow.value = false
  scanContinueShow.value = false
  scanCloseShow.value = false

}


const resultTablePage = ref(1)
const resultTablePageSize = ref(10)

const paginatedData = computed(() => {
  const startIndex = (resultTablePage.value - 1) * resultTablePageSize.value
  const endIndex = startIndex + resultTablePageSize.value
  return result.dataAll.slice(startIndex, endIndex)
})

// // 监听 resultTablePage 和 resultTablePageSize 的变化，更新分页后的数据
// watch([resultTablePage, resultTablePageSize], () => {
//   if(isUpdateResultDate.value) {
//     result.data = JSON.parse(JSON.stringify(paginatedData.value))
//   }
// })

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

// const test = reactive([
//     {
//   "fingerprint": ['apache-tomcat', '亿赛通-电子文档安全管理系统'],
//   "title": "111电子文档安全管理系统",
//   "url": "https://42.237.15.249:8443",
//   "vulnerability": ["亿赛通-电子文档安全管理系统UploadFileFromClientServiceForClient 任意文件上传", "亿赛通-电子文档安全管理系统createSmartSecmysql 信息泄露","亿赛通-电子文档安全管理系统solr 命令执行"]
//   },
//   {
//     "fingerprint": ['apache-tomcat', '亿赛通-电子文档安全管理系统','亿赛通-电子文档安全管理系统'],
//     "title": "电子文档222安全管理系统",
//     "url": "https://14.237.15.249:8443",
//     "vulnerability": ["亿赛通-电子文档安全管理系统UploadFileFromClientServiceForClient 任意文件上传", "亿赛通-电子文档安全管理系统createSmartSecmysql 信息泄露"]
//   },
//
// ])
//
// onMounted(() =>{
//   result.dataAll = test
//   result.data = test
// })

function exportToExcel() {
  let csv = 'URL,标题,指纹,存在的漏洞\n'
  result.dataAll.forEach((item) =>{
    let fingerprintStr = ''
    let vulnerabilityStr = ''
    if(item.fingerprint != null){
      for (let i=0;i<item.fingerprint.length; i++){
        if (i == item.fingerprint.length-1){
          fingerprintStr = fingerprintStr + `${item.fingerprint[i]}`
        }else {
          fingerprintStr = fingerprintStr + `${item.fingerprint[i]}\n`
        }

      }
      fingerprintStr = '"' + fingerprintStr + '"'
    }else {
      fingerprintStr = ""
    }

    if(item.vulnerability != null){
      for (let i=0;i<item.vulnerability.length; i++){
        if (i == item.vulnerability.length-1){
          vulnerabilityStr = vulnerabilityStr + `${item.vulnerability[i]}`
        }else {
          vulnerabilityStr = vulnerabilityStr + `${item.vulnerability[i]}\n`
        }
      }

      vulnerabilityStr = '"' + vulnerabilityStr + '"'
    }else {
      vulnerabilityStr = ""
    }

    csv = csv + item.url + "," + item.title + "," + fingerprintStr  + "," +vulnerabilityStr + "\n"
  })
  let csvContent = 'data:text/csv;charset=utf-8,' + encodeURIComponent(csv);
  let link = document.createElement('a');
  link.setAttribute('href', csvContent);
  let currentTime = new Date().toLocaleString('sv', {year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute:'2-digit'}).replace(/[-\s:]/g, '');
  link.setAttribute('download', currentTime + '_result' +  '.csv');
  // link.setAttribute('download', 'filename.csv'); // 可以自定义下载的文件名
  document.body.appendChild(link);
  link.click();

}


function clearData(){
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
        result.dataAll = [] as Array<FingerprintScanResult>
        result.data = [] as Array<FingerprintScanResult>
        url.value = ''
        searchInput.url = ''
        searchInput.title = ''
        progress.value = 0
      })
}


const resultTableSort = (column) => {
  if (column.prop) { //该列有绑定prop字段走这个分支
    if (column.order === 'ascending') {//当用户点击的是升序按钮，即ascending时
      if  (column.prop === 'fingerprint') {
        result.dataAll = result.dataAll.sort(function (a, b){
          const lengthA = a.fingerprint ? a.fingerprint.length : 0;
          const lengthB = b.fingerprint ? b.fingerprint.length : 0;
          if (lengthA < lengthB) {
            return -1;
          } else if (lengthA > lengthB) {
            return 1;
          } else {
            return 0;
          }
        })
      }else if (column.prop == 'vulnerability') {
        result.dataAll = result.dataAll.sort(function (a, b){
          const lengthA = a.vulnerability ? a.vulnerability.length : 0;
          const lengthB = b.vulnerability ? b.vulnerability.length : 0;
          if (lengthA < lengthB) {
            return -1;
          } else if (lengthA > lengthB) {
            return 1;
          } else {
            return 0;
          }
        })
      }
    } else if (column.order === 'descending') {
      if  (column.prop === 'fingerprint') {
        result.dataAll = result.dataAll.sort(function (a, b){
          const lengthA = a.fingerprint ? a.fingerprint.length : 0;
          const lengthB = b.fingerprint ? b.fingerprint.length : 0;
          if (lengthA < lengthB) {
            return 1;
          } else if (lengthA > lengthB) {
            return -1;
          } else {
            return 0;
          }
        })
      }else if (column.prop === 'vulnerability'){
        result.dataAll = result.dataAll.sort(function (a, b){
          const lengthA = a.vulnerability ? a.vulnerability.length : 0;
          const lengthB = b.vulnerability ? b.vulnerability.length : 0;
          if (lengthA < lengthB) {
            return 1;
          } else if (lengthA > lengthB) {
            return -1;
          } else {
            return 0;
          }
        })
      }
    }

    resultTablePage.value = 2
    resultTablePage.value = 1

  }
}


function  openURL(url) {
  DefaultBrowserOpenUrl(url);
}


const uploadChange = (file) => {
  if (file.raw) {
    const reader = new FileReader()
    reader.onload = () => {
      url.value = reader.result.toString()
    }
    reader.readAsText(file.raw)
  }
}


</script>

<style scoped>

.fingerprintScanResultClass {
  height: 60%;
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