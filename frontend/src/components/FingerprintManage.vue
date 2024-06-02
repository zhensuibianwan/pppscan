<template>
  <el-row style="margin-top: 2%; margin-left: 3%">
    <el-col :span="6" >
      <el-input v-model="fingerprintKey"  style="width: 430px" size="large"  @keyup.enter="fingerprintDBList()">
        <template #prepend>指纹名称</template>
        <template #append><el-button :icon="Search" style="width: 70px"  @click="fingerprintDBList"/></template>
      </el-input>
    </el-col>
    <el-col :span="13" style="margin-left: 19.72%">
      <el-button dark type="primary" :icon="DocumentAdd" size="large" @click="showAdd=true">新增</el-button>
      <el-button type="success" :icon="Edit" size="large"  @click="fingerprintEdit" >编辑</el-button>
      <el-button type="danger" :icon="Delete" size="large" @click="fingerprintDel" >删除</el-button>
      <el-button type="warning" :icon="Download" size="large"  @click="downLoadFingerprints" >保存</el-button>
      <el-button  :icon="Upload" size="large" @click="fingerprintFileList" >导入</el-button>
    </el-col>
  </el-row>
  <div class="fingerprintTableClass">
    <el-table
        ref="fingerprintTableRef"
        :data="fingerprintTable.data"
        style="width: 100%"
        @selection-change="fingerprintSelectionChange"
        height="100%"
        @sort-change='fingerprintTableSort'
    >
      <el-table-column   type="selection" width="35" />
      <el-table-column   label="指纹名称" width="400">
        <template #default="scope">
          {{ scope.row.name }}
        </template>
      </el-table-column>
      <el-table-column prop="fingerprintScan.length" label="指纹数量" sortable='custom'  width="250" />
      <el-table-column prop="pocsInfo.length" label="关联的POC数量"  sortable='custom'  show-overflow-tooltip/>
    </el-table>
  </div>
  <div style="margin-left: 15%">
    <el-pagination
        v-model:current-page="fingerprintTablePage"
        v-model:page-size="fingerprintTablePageSize"
        :page-sizes="[100, 200, 300, 400,500]"
        :small="false"
        :background="true"
        layout="prev, pager, next, jumper,sizes,total"
        :total="fingerprintTable.dataAll.length"
        @size-change="fingerprintTableHandleSizeChange"
        @current-change="fingerprintTableHandleCurrentChange"
    />
  </div>
  <el-dialog
      v-model="showAdd"
      :before-close="handleClose"
      align-center
      center
      width="90%"
      draggable
      title="管理指纹"
  >
    <el-tabs tab-position="top"  class="paneLabelSize" style="margin-left: 1%">
      <el-tab-pane label="指纹信息" >
        <el-row>
          <el-col :span="12" style="margin-left: 20%;margin-top: 5%">
          <el-form-item label="指纹名称" prop="name">
            <el-input v-model="fingerprintForm.name" size="large" style="width: 200%" />
          </el-form-item>
          </el-col>
        </el-row>
      </el-tab-pane>
      <el-tab-pane label="添加指纹" >
        <el-row style="margin-top: 2%; margin-left: 63%">
          <el-button dark type="primary" :icon="DocumentAdd" size="large" @click="showFingerprintScanAdd=true">新增</el-button>
          <el-button type="success" :icon="Edit" size="large"  @click="fingerprintScanEdit" >编辑</el-button>
          <el-button type="danger" :icon="Delete" size="large" @click="fingerprintScanDel" >删除</el-button>
        </el-row>
        <div class="demo-radius">
          <el-table
              ref="fingerprintScanTableRef"
              :data="fingerprintForm.fingerprintScan"
              style="width: 100%"
              @selection-change="fingerprintScanSelectionChange"
              height="100%"
          >
            <el-table-column   type="selection" width="35" />
            <el-table-column   label="id" width="100">
              <template #default="scope">
                <el-tooltip  placement="top" >
                  <template #content>  {{ scope.row.name }} </template>
                  {{ scope.row.priority }}
                </el-tooltip>
              </template>
            </el-table-column>
            <el-table-column property="path" label="请求路径" width="100" />
            <el-table-column property="request_method" label="请求方式" width="100" />
            <el-table-column property="status_code" label="状态码" width="100" />
            <el-table-column label="请求头" width="100" >
            <template #default="scope">
             <span v-for="(value, key) in scope.row.request_headers" :key="key">
                {{ key }}: {{ value }}
               <br/>
            </span>
            </template>
            </el-table-column>
            <el-table-column property="request_data" label="请求数据"  width="100"/>
            <el-table-column label="检测返回头"  width="100">
              <template #default="scope">
             <span v-for="(value, key) in scope.row.headers" :key="key">
                {{ key }}: {{ value }}
               <br/>
            </span>
              </template>
            </el-table-column>
            <el-table-column property="keyword" label="检测body关键值" width="100" />
            <el-table-column property="favicon_hash" label="faviconHash"  width="100"/>
          </el-table>
        </div>
        <el-dialog
            v-model="showFingerprintScanAdd"
            :before-close="fingerprintScanHandleClose"
            align-center
            center
            width="75%"
            draggable
            title="添加指纹"
        >

          <el-form
              :model="fingerprintScanForm"
              label-width="120px"
              status-icon
              style="margin-left: 5%"
          >
<!--            <el-form-item label="优先级">-->
<!--              <el-input v-model="fingerprintScanFrom.priority" size="large" style="width: 85%" @input="priorityToNumber" />-->
<!--            </el-form-item>-->
            <el-form-item label="请求路径" >
              <el-input v-model="fingerprintScanForm.path" size="large" style="width: 85%" />
            </el-form-item>
            <el-form-item label="请求方式">
              <el-select v-model="fingerprintScanForm.request_method" placeholder="please select" size="large" style="width: 85%">
                <el-option label="GET" value="GET" />
                <el-option label="POST" value="POST" />
              </el-select>
<!--              <el-input v-model="fingerprintScanFrom.request_method" size="large" style="width: 85%" />-->
            </el-form-item>
            <el-form-item label="状态码">
              <el-input v-model.number="fingerprintScanForm.status_code" size="large" style="width: 85%"/>
            </el-form-item>
            <el-form-item label="请求头">
              <el-input v-model="request_headersStr" size="large" style="width: 85%" :rows="2" type="textarea" placeholder="一行一个，':'后有空格，例：Location: /login?refer=%2F"/>
            </el-form-item>
            <el-form-item label="请求数据">
              <el-input v-model="fingerprintScanForm.request_data" size="large" style="width: 85%" :rows="2" type="textarea"/>
            </el-form-item>
            <el-form-item label="检测返回头">
              <el-input v-model="headersStr" size="large" :rows="2" type="textarea" style="width: 85%" placeholder="一行一个，':'后有空格，例：Set-Cookie: rememberMe"/>
            </el-form-item>
            <el-form-item label="检测body关键值">
              <el-input v-model="keywordStr" size="large" :rows="2" type="textarea" style="width: 85%" placeholder="一行一个，且的关系"/>
            </el-form-item>
            <el-form-item label="faviconHash">
              <el-input v-model="favicon_hashStr" size="large" :rows="2" type="textarea" style="width: 85%" placeholder="一行一个，或的关系"/>
            </el-form-item>
          </el-form>
          <div style="margin-left: 62%; margin-top: 5%" >
            <el-button  size="large" :icon="CloseBold" @click="fingerprintScanAddCancel" >取消</el-button>
            <el-button size="large" :icon="Select"  type="primary" @click="fingerprintScanAddSave">保存</el-button>
          </div>
        </el-dialog>

      </el-tab-pane>
      <el-tab-pane label="关联POC" >
        <el-row style="margin-top: 2%; margin-left: 74.5%">
          <el-button dark type="primary" :icon="DocumentAdd" size="large" @click="pocsInfoAdd">新增</el-button>
          <el-button type="danger" :icon="Delete" size="large" @click="pocsInfoDel" >删除</el-button>
        </el-row>
        <div class="demo-radius">
          <el-table
              ref="pocsInfoTableRef"
              :data="fingerprintForm.pocsInfo"
              style="width: 100%"
              @selection-change="pocsInfoSelectionChange"
              height="100%"
          >
            <el-table-column   type="selection" width="35" />
            <el-table-column   label="POC名称" width="500">
              <template #default="scope">
                {{ scope.row.name }}
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-tab-pane>
    </el-tabs>
    <div style="margin-left: 75.5%; margin-top: 5%" >
      <el-button  size="large" :icon="CloseBold" @click="fingerprintAddCancel" >取消</el-button>
      <el-button size="large" :icon="Select"  type="primary" @click="fingerprintAddSave">保存</el-button>
    </div>

    <el-dialog
        v-model="pocsInfoAddShow"
        :before-close="pocsInfoAddHandleClose"
        align-center
        center
        width="70%"
        draggable
        title="选择POC"
    >
      <div style="display: flex;margin-left: 3%;" >
        <el-col :span="6" style="">
          <el-input v-model="pocsInfoAddInput"  style="width: 390%; " size="large"  @keyup.enter="pocList()" clearable>
            <template #prepend>poc名称</template>
            <template #append><el-button :icon="Search" style="width: 75px;color: #409eff"  @click="pocList"/></template>
          </el-input>
        </el-col>
      </div>
      <div class="demo-radius3">
        <el-table
            ref="pocsInfoAddRef"
            :data="pocsInfoAddList"
            style="width: 100%"
            @selection-change="pocsInfoAddSelectionChange"
            max-height="400px"
        >
          <el-table-column   type="selection" fixed  width="35" />
          <el-table-column   label="poc名称" width="430">
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
      <div style="margin-left: 68%;margin-top: 2%">
        <el-button  size="large" :icon="CloseBold" @click="pocsInfoAddCancel">取消</el-button>
        <el-button size="large" :icon="Select"  type="primary" @click="pocsInfoAddSave">确认</el-button>
      </div>
    </el-dialog>
  </el-dialog>

  <el-dialog
      v-model="fingerprintUploadShow"
      align-center
      center
      width="70%"
      draggable
      title="选择指纹"
      :before-close="fingerprintUploadHandleClose"
  >
    <div style="display: flex;margin-left: 3%;" >
      <el-col :span="6" style="">
        <el-input v-model="inputUpload"  style="width: 300%; " size="large"  @keyup.enter="FingerprintFileList()" clearable>
          <template #prepend>文件名称</template>
          <template #append><el-button :icon="Search" style="width: 75px;color: #409eff"  @click="fingerprintFileList"/></template>
        </el-input>
      </el-col>
      <el-button  style="margin-left: 58%" type="danger" :icon="Delete" size="large" @click="RemoveFingerprintFile" >删除</el-button>
    </div>
    <div class="demo-radius3">
      <el-table
          ref="UploadFingerprintListRef"
          :data="file.fingerprintData"
          style="width: 100%"
          @selection-change="UploadFingerprintSelectionChange"
          max-height="400px"
      >
        <el-table-column   type="selection" fixed  width="35" />
        <el-table-column   label="指纹名称" width="580">
          <template #default="scope">
            {{ scope.row.name }}
          </template>
        </el-table-column>
      </el-table>
    </div>
    <!--    <el-button size="large" :icon="Select"  type="primary" @click="" style="margin-top: 2%; margin-left: 86%">确认</el-button>-->
    <div style="margin-left: 68%;margin-top: 2%">
      <el-button  size="large" :icon="CloseBold" @click="fingerprintUploadClose">取消</el-button>
      <el-button size="large" :icon="Select"  type="primary" @click="UploadFingerprintToDB">确认</el-button>
    </div>
  </el-dialog>

</template>

<script lang="ts" setup>

import {Delete, DocumentAdd, Edit, Search, CloseBold, Select,Download,Upload} from '@element-plus/icons-vue'
import {reactive, ref} from "vue";
import {
  AddFingerprint,
  DelFingerprint, DelFingerprintFile,
  FingerprintFileList,
  FingerprintList, FingerprintLocalList, FingerprintsDownload, FingerprintsUpload, PocList, PocList2, PocsUpload
} from "../../wailsjs/go/main/App";
import {publicCode} from "../../wailsjs/go/models";
import Fingerprint = publicCode.Fingerprint;
import {ElMessage, ElMessageBox, ElTable} from "element-plus";

const fingerprintKey = ref('')

const fingerprintTable = reactive({
  data: [] as Array<Fingerprint>,
  dataAll: [] as Array<Fingerprint>,
})

async function fingerprintDBList() {
  fingerprintTable.dataAll = await FingerprintList(fingerprintKey.value)
  fingerprintTable.data = JSON.parse(JSON.stringify(fingerprintTable.dataAll))
  fingerprintTable.data = fingerprintTable.data.splice((fingerprintTablePage.value - 1) * fingerprintTablePageSize.value, fingerprintTablePageSize.value)
}

async function fingerprintLocalList() {
  fingerprintTable.dataAll = await FingerprintLocalList(fingerprintKey.value)
  fingerprintTable.data = JSON.parse(JSON.stringify(fingerprintTable.dataAll))
  fingerprintTable.data = fingerprintTable.data.splice((fingerprintTablePage.value - 1) * fingerprintTablePageSize.value, fingerprintTablePageSize.value)
}

function fingerprintList(){
  fingerprintLocalList()
  fingerprintDBList()
}

fingerprintList()

const fingerprintTablePage = ref(1)
const fingerprintTablePageSize = ref(100)

function fingerprintTableHandleSizeChange(){
  fingerprintTable.data = JSON.parse(JSON.stringify(fingerprintTable.dataAll))
  fingerprintTable.data = fingerprintTable.data.splice((fingerprintTablePage.value - 1) * fingerprintTablePageSize.value, fingerprintTablePageSize.value)
}

function fingerprintTableHandleCurrentChange(){
  fingerprintTable.data = JSON.parse(JSON.stringify(fingerprintTable.dataAll))
  fingerprintTable.data = fingerprintTable.data.splice((fingerprintTablePage.value - 1) * fingerprintTablePageSize.value, fingerprintTablePageSize.value)
}

const fingerprintTableSort = (column) => {
  if (column.prop) { //该列有绑定prop字段走这个分支
    fingerprintTable.data = JSON.parse(JSON.stringify(fingerprintTable.dataAll))
    if (column.order == 'ascending') {//当用户点击的是升序按钮，即ascending时
      if  (column.prop == 'fingerprintScan.length') {
        fingerprintTable.data = fingerprintTable.data.sort(function (a, b){return a.fingerprintScan.length - b.fingerprintScan.length})
      }else if (column.prop == 'pocsInfo.length') {
        fingerprintTable.data = fingerprintTable.data.sort(function (a, b){return a.pocsInfo.length - b.pocsInfo.length})
      }
      fingerprintTable.data = fingerprintTable.data.splice((fingerprintTablePage.value - 1) * fingerprintTablePageSize.value, fingerprintTablePageSize.value)
    } else if (column.order == 'descending') {
      //当用户点击的是升序按钮，即descending时
      if  (column.prop == 'fingerprintScan.length') {

        fingerprintTable.data = fingerprintTable.data.sort(function (a, b){return b.fingerprintScan.length - a.fingerprintScan.length})
      }else if (column.prop == 'pocsInfo.length'){
        fingerprintTable.data = fingerprintTable.data.sort(function (a, b){return b.pocsInfo.length - a.pocsInfo.length})
      }
      fingerprintTable.data = fingerprintTable.data.splice((fingerprintTablePage.value - 1) * fingerprintTablePageSize.value, fingerprintTablePageSize.value)
    }
  }
}




const fingerprintForm = reactive({
  uuid: '',
  name: '',
  fingerprintScan: [] as publicCode.FingerprintScanData[],
  pocsInfo: [] as publicCode.PocsInfoData[],
} as publicCode.Fingerprint)


const fingerprintScanForm = reactive({
  path: '',
  request_method: '',
  request_headers: {} as {[key: string]: string},
  request_data: '',
  status_code: 0,
  headers: {} as {[key: string]: string},
  keyword: [] as string[],
  favicon_hash: [] as string[],
  priority: 0,
  name: '',
} as publicCode.FingerprintScanData)

const request_headersStr = ref('')
const headersStr = ref('')
const keywordStr = ref('')
const favicon_hashStr = ref('')

const pocsInfoForm = reactive({
  uuid: '',
  name: ''
})

const fingerprintTableRef = ref<InstanceType<typeof ElTable>>()
const fingerprintSelection = ref<publicCode.Fingerprint[]>([])

const fingerprintSelectionChange = (val: publicCode.Fingerprint[]) => {

  fingerprintSelection.value = val
}

const showAdd = ref(false)

const handleClose = (done: () => void) => {
  done()
  location.reload()
}

function fingerprintAddCancel(){
  showAdd.value = false
  location.reload()
}

let oldFingerprintName = ''

function fingerprintAddSave(){
  ElMessageBox.confirm(
      '确定保存吗？',
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
  )
      .then(async () => {
        if (fingerprintForm.name != ''){
          if (oldFingerprintName != fingerprintForm.name && oldFingerprintName != ''){
            if(fingerprintForm.fingerprintScan != null) {
              for (let i = 0; i < fingerprintForm.fingerprintScan.length; i++){
                fingerprintForm.fingerprintScan[i].name = fingerprintForm.name
              }
            }
          }
          let msg = await AddFingerprint(fingerprintForm)
          console.log(msg)
          showAdd.value = false
          if (msg[0]==="1"){
            ElMessage({
              type: 'success',
              message: msg[1],
            })
          }else {
            ElMessage({
              type: 'error',
              message: msg[1],
            })
          }
          setTimeout(function() {
            location.reload();
          }, 800);
        }else {
          ElMessage.error("指纹名称不能为空！")
        }

      })
      .catch(() => {
      })

}

function fingerprintEdit(){
  if(fingerprintSelection.value.length == 1) {
    showAdd.value = true
    oldFingerprintName = fingerprintSelection.value[0].name
    fingerprintForm.name = fingerprintSelection.value[0].name
    fingerprintForm.uuid = fingerprintSelection.value[0].uuid
    fingerprintForm.fingerprintScan = fingerprintSelection.value[0].fingerprintScan
    fingerprintForm.pocsInfo = fingerprintSelection.value[0].pocsInfo
  }else if (fingerprintSelection.value.length > 1){
    ElMessage.error("只能同时编辑一个！")
  }

}



function fingerprintDel(){
  ElMessageBox.confirm(
      '确定删除吗？',
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
  )
      .then(async () => {
        if (fingerprintSelection.value.length >= 1) {
          let msg = await DelFingerprint(fingerprintSelection.value)
          console.log(msg)
          if (msg[0] === '1') {
            ElMessage({
              type: 'success',
              message: msg[1],
            })
          } else {
            ElMessage({
              type: 'error',
              message: msg[1],
            })
          }
          await fingerprintDBList()

        }
      })
      .catch(() => {
      })
}

function downLoadFingerprints(){
  if (fingerprintSelection.value.length > 0){
    ElMessageBox.confirm(
        '确定保存选中到本地吗？',
        '提示',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
        }
    )
        .then(async () => {
          let msg = await FingerprintsDownload(fingerprintSelection.value)
          console.log(msg)
          if (msg[0]==='1'){
            ElMessage({
              type: 'success',
              message: msg[1],
            })
          }else {
            ElMessage({
              type: 'error',
              message: msg[1],
            })
          }
          fingerprintTableRef.value!.clearSelection()
        })

        .catch(() => {})
  }
}

const fingerprintUploadShow = ref(false)
const UploadFingerprintListRef = ref<InstanceType<typeof ElTable>>()
const inputUpload = ref('')
const fingerprintFileSelection = ref<publicCode.Fingerprint[]>([])
const file = reactive({
  fingerprintDataAll: [] as Array<publicCode.Fingerprint>,
  fingerprintData: [] as Array<publicCode.Fingerprint>,
})
const UploadFingerprintSelectionChange = (val: publicCode.Fingerprint[]) => {

  fingerprintFileSelection.value = val
}

async function fingerprintFileList(){
  fingerprintUploadShow.value = true

  file.fingerprintData = await FingerprintFileList(inputUpload.value)
}

function RemoveFingerprintFile(){
  ElMessageBox.confirm(
      '确定删除选中的本地指纹JSON文件吗？',
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
  )
      .then(async () => {
        let msg = await DelFingerprintFile(fingerprintFileSelection.value)
        console.log(msg)
        if (msg[0]==='1'){
          ElMessage({
            type: 'success',
            message: msg[1],
          })
        }else {
          ElMessage({
            type: 'error',
            message: msg[1],
          })
        }
        fingerprintUploadShow.value = false
        UploadFingerprintListRef.value!.clearSelection()
        await fingerprintFileList()
      })

      .catch(() => {})
}

function UploadFingerprintToDB(){
  ElMessageBox.confirm(
      '确定导入选中的指纹到数据库吗？',
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
  )
      .then(async () => {
        let msg = await FingerprintsUpload(fingerprintFileSelection.value)
        console.log(msg)
        if (msg[0]==='1'){
          ElMessage({
            type: 'success',
            message: msg[1],
          })
        }else {
          ElMessage({
            type: 'error',
            message: msg[1],
          })
        }
        fingerprintUploadClose()
        await fingerprintList()
      })

      .catch(() => {})
}


function fingerprintUploadClose() {
  fingerprintUploadShow.value = false
  file.fingerprintData = null
  file.fingerprintDataAll = null
  UploadFingerprintListRef.value!.clearSelection()
}

const fingerprintUploadHandleClose = (done: () => void) => {
  done()
  file.fingerprintData = null
  file.fingerprintDataAll = null
  UploadFingerprintListRef.value!.clearSelection()
}


const fingerprintScanTableRef = ref<InstanceType<typeof ElTable>>()
const fingerprintScanSelection = ref<publicCode.FingerprintScanData[]>([])

const fingerprintScanSelectionChange = (val: publicCode.FingerprintScanData[]) => {

  fingerprintScanSelection.value = val

}

const showFingerprintScanAdd = ref(false)

const fingerprintScanHandleClose = (done: () => void) => {
  done()
  fingerprintScanForm.name = ''
  fingerprintScanForm.path = ''
  fingerprintScanForm.request_method = ''
  fingerprintScanForm.request_headers= {}as {[key: string]: string}
  fingerprintScanForm.request_data= ''
  fingerprintScanForm.status_code= 0
  fingerprintScanForm.headers =  {} as {[key: string]: string}
  fingerprintScanForm.keyword = [] as string[]
  fingerprintScanForm.favicon_hash = [] as string[]
  fingerprintScanForm.priority = 0
  request_headersStr.value = ''
  headersStr.value = ''
  keywordStr.value = ''
  favicon_hashStr.value = ''
}

function fingerprintScanAddCancel(){
  showFingerprintScanAdd.value = false
  fingerprintScanForm.name = ''
  fingerprintScanForm.path = ''
  fingerprintScanForm.request_method = ''
  fingerprintScanForm.request_headers= {}as {[key: string]: string}
  fingerprintScanForm.request_data= ''
  fingerprintScanForm.status_code= 0
  fingerprintScanForm.headers =  {} as {[key: string]: string}
  fingerprintScanForm.keyword = [] as string[]
  fingerprintScanForm.favicon_hash = [] as string[]
  fingerprintScanForm.priority = 0
  request_headersStr.value = ''
  headersStr.value = ''
  keywordStr.value = ''
  favicon_hashStr.value = ''
}

function fingerprintScanAddSave() {
  ElMessageBox.confirm(
      '确定保存吗？',
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
  )
      .then(async () => {

        if (fingerprintScanForm.priority == 0){
          let fingerprintScanNum = fingerprintForm.fingerprintScan.length
          fingerprintScanForm.priority = fingerprintScanNum + 1

          fingerprintScanForm.path= fingerprintScanForm.path.replace(" ", "")

        if(request_headersStr.value != '' ){
          let request_headersAry1 =  request_headersStr.value.split('\n')
          request_headersAry1.forEach((ary) => {
            let request_headersAry2 = ary.split(": ",2)
            fingerprintScanForm.request_headers[request_headersAry2[0]] = request_headersAry2[1]
          })
        }
        if(headersStr.value != '' ){
          let headersAry1 =  headersStr.value.split('\n')
          headersAry1.forEach((ary) => {
            let headersAry2 = ary.split(": ",2)
            fingerprintScanForm.headers[headersAry2[0]] = headersAry2[1]
          })
        }

        if(keywordStr.value != '' ){
          let keywordAry =  keywordStr.value.split('\n')
          for (let k=0;k<keywordAry.length;k++){
            fingerprintScanForm.keyword.push(keywordAry[k])
          }
        }

        if(favicon_hashStr.value != '' ){
          let favicon_hashAry =  favicon_hashStr.value.split('\n')
          for (let f=0;f<favicon_hashAry.length;f++){
            fingerprintScanForm.favicon_hash.push(favicon_hashAry[f])
          }
        }
          fingerprintScanForm.name = fingerprintForm.name
        let fingerprintScanObject = JSON.parse(JSON.stringify(fingerprintScanForm))
        fingerprintForm.fingerprintScan.push(fingerprintScanObject)
        showFingerprintScanAdd.value = false
        fingerprintScanAddCancel()
        }else {
          for(let i=0;i<fingerprintForm.fingerprintScan.length;i++){
            if(fingerprintForm.fingerprintScan[i].priority == fingerprintScanForm.priority){
              fingerprintScanForm.path = fingerprintScanForm.path.replace(" ", "")
              fingerprintForm.fingerprintScan[i].path = fingerprintScanForm.path
              fingerprintForm.fingerprintScan[i].request_data = fingerprintScanForm.request_data
              fingerprintForm.fingerprintScan[i].request_method = fingerprintScanForm.request_method
              fingerprintForm.fingerprintScan[i].status_code = fingerprintScanForm.status_code

              if(request_headersStr.value != '' ){
                let request_headersAry1 =  request_headersStr.value.split('\n')
                request_headersAry1.forEach((ary) => {
                  if(ary != ''){
                    let request_headersAry2 = ary.split(": ",2)
                    fingerprintScanForm.request_headers[request_headersAry2[0]] = request_headersAry2[1]
                  }
                })
              }
              fingerprintForm.fingerprintScan[i].request_headers = fingerprintScanForm.request_headers

              if(headersStr.value != '' ){
                let headersAry1 =  headersStr.value.split('\n')
                headersAry1.forEach((ary) => {
                  if (ary != ''){
                    let headersAry2 = ary.split(": ",2)
                    fingerprintScanForm.headers[headersAry2[0]] = headersAry2[1]
                  }
                })
              }
              fingerprintForm.fingerprintScan[i].headers = fingerprintScanForm.headers

              if(keywordStr.value != '' ){
                fingerprintScanForm.keyword = [] as string[]
                let keywordAry =  keywordStr.value.split('\n')
                for (let k=0;k<keywordAry.length;k++){
                  if(keywordAry[k] != ''){
                    fingerprintScanForm.keyword.push(keywordAry[k])
                  }
                }
              }
              fingerprintForm.fingerprintScan[i].keyword = fingerprintScanForm.keyword

              if(favicon_hashStr.value != '' ){
                fingerprintScanForm.favicon_hash = [] as string[]
                let favicon_hashAry =  favicon_hashStr.value.split('\n')
                for (let f=0;f<favicon_hashAry.length;f++){
                  if(favicon_hashAry[f] != ''){
                    fingerprintScanForm.favicon_hash.push(favicon_hashAry[f])
                  }
                }
              }
              fingerprintForm.fingerprintScan[i].favicon_hash = fingerprintScanForm.favicon_hash

              fingerprintScanAddCancel()
            }

          }
        }

      })
      .catch(() => {
      })

}


function fingerprintScanEdit() {
  if(fingerprintScanSelection.value.length == 1){
    showFingerprintScanAdd.value = true
    fingerprintScanForm.priority = fingerprintScanSelection.value[0].priority
    fingerprintScanForm.path = fingerprintScanSelection.value[0].path
    fingerprintScanForm.name = fingerprintScanSelection.value[0].name
    fingerprintScanForm.status_code =  fingerprintScanSelection.value[0].status_code
    fingerprintScanForm.request_method = fingerprintScanSelection.value[0].request_method
    fingerprintScanForm.request_headers = fingerprintScanSelection.value[0].request_headers
    fingerprintScanForm.request_data = fingerprintScanSelection.value[0].request_data
    fingerprintScanForm.headers = fingerprintScanSelection.value[0].headers
    fingerprintScanForm.keyword = fingerprintScanSelection.value[0].keyword
    fingerprintScanForm.favicon_hash = fingerprintScanSelection.value[0].favicon_hash
    for (let stringsKey in fingerprintScanForm.request_headers) {
        request_headersStr.value = request_headersStr.value + stringsKey + ": " + fingerprintScanForm.request_headers[stringsKey] + "\n"
    }
    for (let stringsKey in fingerprintScanForm.headers) {
        headersStr.value = headersStr.value + stringsKey + ": " + fingerprintScanForm.headers[stringsKey] + "\n"
    }
    for (let i=0; i<fingerprintScanForm.keyword.length; i++ ){
      keywordStr.value = keywordStr.value + fingerprintScanForm.keyword[i] + "\n"
    }
    for (let i=0; i<fingerprintScanForm.favicon_hash.length; i++ ){
      favicon_hashStr.value = favicon_hashStr.value + fingerprintScanForm.favicon_hash[i] + "\n"
    }
  }else if(fingerprintScanSelection.value.length >1){
    ElMessage.error("只能同时编辑一个！")
  }
}



function fingerprintScanDel() {
  ElMessageBox.confirm(
      '确定删除吗？',
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
  )
      .then(() => {
        if(fingerprintScanSelection.value.length >= 1) {
          fingerprintScanSelection.value.forEach((data)=>{
            for (let i=0;i<fingerprintForm.fingerprintScan.length;i++){
              if(fingerprintForm.fingerprintScan[i].priority == data.priority){
                fingerprintForm.fingerprintScan.splice(i,1)
              }
            }
          })
          ElMessage.success("刪除成功！")
        }
      })
      .catch(() => {
      })

}




const pocsInfoTableRef = ref<InstanceType<typeof ElTable>>()
const pocsInfoSelection = ref<publicCode.PocsInfoData[]>([])

const pocsInfoSelectionChange = (val: publicCode.PocsInfoData[]) => {

  pocsInfoSelection.value = val

}

async function pocsInfoAdd() {
  let list = await PocList(pocsInfoAddInput.value)
  pocsInfoAddList.value = list.filter(data => data.cms.includes(fingerprintForm.name))
  pocsInfoAddShow.value = true
}

function pocsInfoDel() {
  ElMessageBox.confirm(
      '确定删除吗？',
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
  )
      .then(() => {
        if(pocsInfoSelection.value.length >= 1) {
          pocsInfoSelection.value.forEach((data)=>{
            for (let i=0;i<fingerprintForm.pocsInfo.length;i++){
              if(fingerprintForm.pocsInfo[i].uuid== data.uuid){
                fingerprintForm.pocsInfo.splice(i,1)
              }
            }
          })
          ElMessage.success("刪除成功！")
        }
      })
      .catch(() => {
      })

}

const  pocsInfoAddShow = ref(false)
const  pocsInfoAddList = ref<publicCode.Poc[]>([])
const  pocsInfoAddInput = ref('')

const  pocsInfoAddRef =  ref<InstanceType<typeof ElTable>>()
const  pocsInfoAddSelection = ref<publicCode.Poc[]>([])
const pocsInfoAddSelectionChange = (val: publicCode.Poc[]) => {

  pocsInfoAddSelection.value = val

}

async function pocList() {
  pocsInfoAddList.value = await PocList(pocsInfoAddInput.value)
}



function pocsInfoAddSave(){
  if(pocsInfoAddSelection.value.length >= 1) {
    pocsInfoAddSelection.value.forEach((data) => {
      pocsInfoForm.name = data.name
      pocsInfoForm.uuid = data.uuid
      fingerprintForm.pocsInfo.push(JSON.parse(JSON.stringify(pocsInfoForm)))
    })
    pocsInfoAddShow.value = false
    pocsInfoAddRef.value!.clearSelection()
    pocsInfoAddList.value = [] as publicCode.Poc[]
    ElMessage.success("添加成功！")
  }
}


function pocsInfoAddCancel(){
  pocsInfoAddShow.value = false
  pocsInfoAddRef.value!.clearSelection()
  pocsInfoAddList.value = [] as publicCode.Poc[]
}

const pocsInfoAddHandleClose = (done: () => void) => {
  done()
  pocsInfoAddRef.value!.clearSelection()
  pocsInfoAddList.value = [] as publicCode.Poc[]
}
</script>


<style scoped>

.fingerprintTableClass {
  height: 70%;
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