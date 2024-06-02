<template >
  <div class="topUI" >
    <el-col :span="6">
      <el-input v-model="input"  style="width: 430px" size="large"  @keyup.enter="pocList()">
        <template #prepend>漏洞名称</template>
        <template #append><el-button :icon="Search" style="width: 70px"  @click="pocList"/></template>
      </el-input>
    </el-col>
    <div class="flex" style="margin-left: 21.5%">
      <el-button dark type="primary" :icon="DocumentAdd" size="large" @click="showAdd = true">新增</el-button>
      <el-button type="success" :icon="Edit" size="large"  @click="editPoc" >编辑</el-button>
      <el-button type="danger" :icon="Delete" size="large" @click="pocDel" >删除</el-button>
      <el-button type="warning" :icon="Download" size="large"  @click="downLoadPocs" >保存</el-button>
      <el-button  :icon="Upload" size="large" @click="pocList2" >导入</el-button>
    </div>
  </div>

  <div class="demo-radius">
    <el-table
        ref="multipleTableRef"
        :data="table.tableData"
        style="width: 100%"
        @selection-change="handleSelectionChange"
        height="100%"
    >
      <el-table-column   type="selection" width="35" />
      <el-table-column   label="漏洞名称" width="400">
        <template #default="scope">
        <el-tooltip  placement="top" >
          <template #content>  {{ scope.row.description }} </template>
          {{ scope.row.name }}
        </el-tooltip>
        </template>
      </el-table-column>
      <el-table-column property="hunter" label="鹰图语法" width="260" />
      <el-table-column property="fofa" label="FOFA语法" show-overflow-tooltip />
    </el-table>
  </div>
  <div style="margin-left: 15%">
    <el-pagination
        v-model:current-page="pocManageTablePage"
        v-model:page-size="pocManageTablePageSize"
        :page-sizes="[100, 200, 300, 400,500]"
        :small="false"
        :background="true"
        layout="prev, pager, next, jumper,sizes,total"
        :total="table.tableDataAll.length"
        @size-change="pocManageTableHandleSizeChange"
        @current-change="pocManageTableHandleCurrentChange"
    />
  </div>


<!------------------------------------  AddPoc弹窗页面----------------------------------------->
  <el-dialog
      v-model="showAdd"
      :before-close="handleClose"
      align-center
      center
      width="90%"
      draggable
      title="添加POC"
  >
      <el-steps :active="addActive" finish-status="success" simple class="step" >
        <el-step title="漏洞描述" />
        <el-step title="Poc编写" />
        <el-step :title="pocUIValue[p.id-1].title" v-for="p in pocUI" :key="p.id" />
      </el-steps>
      <div class="addPoc">
        <div class="form" v-show="addStep1.show">
          <el-form
              :model="form"
              label-width="120px"
              ref="addPocRuleFormRef"
              status-icon
              :rules="addPocRules"
          >
            <el-form-item label="漏洞名称" prop="name">
              <el-input v-model="form.name" size="large" style="width: 85%" />
            </el-form-item>
            <el-form-item label="鹰图语法">
              <el-input v-model="form.hunter" size="large" style="width: 85%"/>
            </el-form-item>
            <el-form-item label="FOFA语法">
              <el-input v-model="form.fofa" size="large" style="width: 85%"/>
            </el-form-item>
            <el-form-item label="CMS" prop="cms">
              <el-autocomplete
                  v-model="form.cms"
                  :fetch-suggestions="cmsSearch"
                  clearable
                  size="large"
                  style="width: 85%"
              />
<!--              <el-input v-model="form.cms" size="large" style="width: 85%"/>-->
            </el-form-item>
            <el-form-item label="漏洞说明">
              <el-input v-model="form.description" size="large" :rows="3" type="textarea" style="width: 85%"/>
            </el-form-item>
            <div style="margin-left: 60%; margin-bottom: 2%">
              <el-button  size="large" :icon="CloseBold" @click="addPocCancel">取&emsp;消</el-button>
              <el-button size="large" :icon="ArrowRight"  type="primary" @click="nextStep">下一步</el-button>
            </div>
          </el-form>
        </div>
        <div v-show="addStep2.show">
          <el-form  label-width="70px"  :label-position="'left'" style="margin-top: 3%" size="large">
          <el-form-item label="请求次数" prop="region" style=" margin-left: 2%" size="large">
            <el-select v-model="form.optionValue" clearable placeholder="Select" @change="pocsUIChange"  style="width: 40%;" size="large">
              <el-option
                  v-for="item in options"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
              />
            </el-select>
            <div style="margin-left: 15%">
              <el-button  size="large" :icon="ArrowLeft" @click="lastStep" >上一步</el-button>
              <el-button size="large" :icon="ArrowRight"  type="primary" @click="nextStep" v-bind:disabled="checkOption.value">下一步</el-button>
            </div>
          </el-form-item>
          </el-form>
        </div>

        <dev style="display: flex;" v-for="p in pocUI" :key="p.id" v-show="pocUIValue[p.id-1].show">

          <dev style="width: 50%">
            <el-input
                v-model="form.request[p.id-1].pocString"
                :rows="20"
                type="textarea"
                placeholder="请求包"
                style="margin-left: 3%;margin-top: 2%;margin-bottom: 2%"
                id="request"
                @blur="blurEvent"
                @focus="cursorPosition_Refresh1(p.id)"
                @change="cursorPosition_Refresh1(p.id)"
            />
          </dev>

          <dev style="margin-left: 3%;width: 45%;margin-top: 2%">
            <el-form-item label="操作栏">
              <el-button  size="large" :icon="Setting" @click="addPoc_Setting(p.id)" style="width: 31%;">设置</el-button>
              <el-button  size="large" :icon="RefreshRight" @click="addPoc_Reset(p.id)" style="width: 31%; margin-left: 3%" type="danger">重置</el-button>
              <el-button  size="large" :icon="Plus" @click="appendValue(p.id)" style="width: 31%; margin-left: 3%" type="primary" >插入</el-button>
            </el-form-item>
            <el-form-item label="添加值"  style="margin-top: 10%">
              <el-select v-model="pocUIValue[p.id-1].addValue" placeholder="选择要添加的值" style="width: 100%" size="large" clearable>
                <el-option
                    v-for="item in addData.value"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                />
              </el-select>
<!--              <el-button  size="large" :icon="Plus" @click="appendValue(p.id)" style="width: 26%; margin-left: 3%" type="primary" >Add</el-button>-->
            </el-form-item>

            <el-form-item label="状态码" style="margin-top: 10%">
              <el-input v-model="form.request[p.id-1].status" size="large" style="width: 100%" clearable placeholder="返回包的状态码"/>
            </el-form-item>

<!--            <el-form-item label="&nbsp;input&nbsp;&nbsp;" style="margin-top: 10%">-->
<!--              <el-input v-model="form.request[p.id-1].Input" size="large" style="width: 100%" clearable placeholder="请输入用于检测的无害化的值"/>-->
<!--            </el-form-item>-->


            <el-form-item label="判断值" style="margin-top: 10%">
              <el-input v-model="form.request[p.id-1].check" size="large" style="width: 100%" clearable   @change="cursorPosition_Refresh3(p.id)" @focus="cursorPosition_Refresh3(p.id)" @blur="blurEvent" placeholder="同时只支持&&、||其中一种运算"/>
            </el-form-item>
            <el-form-item label="输出值" style="margin-top: 10%">
              <el-input v-model="form.request[p.id-1].print" size="large" style="width: 100%" clearable  @change="cursorPosition_Refresh2(p.id)" @focus="cursorPosition_Refresh2(p.id)" @blur="blurEvent" placeholder="漏洞存在的程序输出" type="textarea"/>
            </el-form-item>

            <div style="margin-left: 38%; margin-top: 6%" v-show="pocUIValue[p.id-1].buttonShow1">
              <el-button  size="large" :icon="ArrowLeft" @click="lastStep" >上一步</el-button>
              <el-button size="large" :icon="ArrowRight"  type="primary" @click="nextStep">下一步</el-button>
            </div>

            <div style="margin-left: 38%; margin-top: 6%" v-show="pocUIValue[p.id-1].buttonShow2">
              <el-button  size="large" :icon="ArrowLeft" @click="lastStep" >上一步</el-button>
              <el-button size="large" :icon="Select"  type="primary" @click="addPoc_Save">完&nbsp;&nbsp;&nbsp;&nbsp;成</el-button>
            </div>
          </dev>
          <el-dialog
              v-model="pocUIValue[p.id-1].settingShow"
              align-center
              center
              width="80%"
              draggable
              title="设置返回值"
          >
            <div>
              <div class="flex" style="margin-left: 60%;">
                <el-button color="#1a7bb9" dark type="primary" :icon="DocumentAdd" size="large" @click="addPoc_SettingAdd(p.id)">新增</el-button>
                <el-button type="success" :icon="Edit" size="large" @click="needDataEdit(p.id)">编辑</el-button>
                <el-button type="danger" :icon="Delete" size="large"  @click="needDataDel(p.id)">删除</el-button>
              </div>
              <el-table
                  ref="needDataTableRef"
                  :data="form.needData"
                  style="width: 100%"
                  @selection-change="needDataChange"
                  height="100%"
              >
                <el-table-column   type="selection" width="35" />
                <el-table-column   label="name" width="350">
                  <template #default="scope">
                      {{ scope.row.label }}
                  </template>
                </el-table-column>
                <el-table-column property="value" label="value" width="350" show-overflow-tooltip/>
              </el-table>
            </div>

            <div style="margin-left: 38%; margin-top: 5%" >
              <el-button  size="large" :icon="CloseBold" @click="addPoc_SettingCancel(p.id)" >取消</el-button>
              <el-button size="large" :icon="Select"  type="primary" @click="addPoc_SettingCancel(p.id)">完成</el-button>
            </div>
          </el-dialog>

          <el-dialog
              v-model="pocUIValue[p.id-1].settingAddShow"
              align-center
              center
              width="60%"
              draggable
              title="设置返回值"
          >
            <el-button  size="large" :icon="Refresh" @click="symbolReplace(p.id)" style="width: 15%; margin-left: 80%" type="info" >替换</el-button>
            <el-form-item label="返回值"  style="margin-top: 2%">
              <el-select v-model="pocUIValue[p.id-1].settingName" placeholder="需要的返回值" style="width: 45%" size="large" clearable  :disabled="pocUIValue[p.id-1].NeedSelect">
                <el-option
                    v-for="item in pocUIValue[p.id-1].NeedKey"
                    :key="item.value"
                    :label="item.label+'.'+addData.addNum"
                    :value="item.value+'.'+addData.addNum"
                    :disabled="item.disabled"
                />
              </el-select>
              <el-input v-model="pocUIValue[p.id-1].settingValue" size="large" style="width: 45%; margin-left: 5%" @blur="blurEvent" clearable type="textarea"/>
            </el-form-item>

            <div style="margin-left: 60%; margin-top: 5%" >
              <el-button  size="large" :icon="CloseBold" @click="addPoc_SettingAddCancel(p.id)" >取消</el-button>
              <el-button size="large" :icon="Select"  type="primary" @click="addPoc_SettingAddSave(p.id)">保存</el-button>
            </div>

          </el-dialog>

        </dev>

      </div>

  </el-dialog>
<!------------------------------------  AddPoc弹窗页面----------------------------------------->



  <!------------------------------------  导入弹窗页面----------------------------------------->
  <el-dialog
      v-model="UploadShow"
      align-center
      center
      width="70%"
      draggable
      title="选择POC"
      :before-close="UploadHandleClose"
  >
    <div style="display: flex;margin-left: 3%;" >
    <el-col :span="6" style="">
      <el-input v-model="inputUpload"  style="width: 300%; " size="large"  @keyup.enter="pocList2()" clearable>
        <template #prepend>文件名称</template>
        <template #append><el-button :icon="Search" style="width: 75px;color: #409eff"  @click="pocList2"/></template>
      </el-input>
    </el-col>
    <el-button  style="margin-left: 58%" type="danger" :icon="Delete" size="large" @click="RemovePocFile" >删除</el-button>
    </div>
    <div class="demo-radius3">
      <el-table
          ref="UploadPocListRef"
          :data="file.pocData"
          style="width: 100%"
          @selection-change="UploadPocSelectionChange"
          max-height="400px"
      >
        <el-table-column   type="selection" fixed  width="35" />
        <el-table-column   label="poc名称" width="600">
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
      <el-button  size="large" :icon="CloseBold" @click="UploadClose">取消</el-button>
      <el-button size="large" :icon="Select"  type="primary" @click="UploadPocToDB">确认</el-button>
    </div>
  </el-dialog>


</template>

<script lang="ts" setup>
import {Delete, DocumentAdd, Edit, Search, CloseBold, ArrowRight, ArrowLeft, Select,Plus, Setting,RefreshRight,Refresh,Download,Upload} from '@element-plus/icons-vue'
import {onMounted, reactive, ref} from 'vue'
import {ElMessage, ElMessageBox} from 'element-plus'
import type {FormInstance, FormRules} from 'element-plus'
import {ElTable} from 'element-plus'
import {
  PocList,
  AddPoc,
  DelPoc,
  PocsDownload,
  PocList2,
  PocsUpload,
  DelPocFile, LocalList, FingerprintGetName, FingerprintList
} from '../../wailsjs/go/main/App'
import {publicCode} from "../../wailsjs/go/models";

const handleClose = (done: () => void) => {
  ElMessageBox.confirm(
      '确定关闭窗口？',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
      )
      .then(() => {
        //location.reload()
        multipleTableRef.value!.clearSelection()
        AddOREditCancel()
        done()
      })
      .catch(() => {
        // catch error
      })
}


// 获取光标
let selectionStart
let selectionEnd
function blurEvent(e) {
  selectionStart = e.currentTarget.selectionStart;
  selectionEnd = e.currentTarget.selectionEnd;
}



//--------  Main start  ---------------------

const input = ref('')
const inputUpload = ref('')
const UploadShow = ref(false)

const multipleTableRef = ref<InstanceType<typeof ElTable>>()
const UploadPocListRef = ref<InstanceType<typeof ElTable>>()
const pocSelection = ref<publicCode.Poc[]>([])
const pocFileSelection = ref<publicCode.Poc[]>([])

const handleSelectionChange = (val: publicCode.Poc[]) => {

  pocSelection.value = val
}



const UploadPocSelectionChange = (val: publicCode.Poc[]) => {

  pocFileSelection.value = val
}



function UploadClose() {
  UploadShow.value = false
  file.pocData = null
  file.pocDataAll = null
  UploadPocListRef.value!.clearSelection()
}

const UploadHandleClose = (done: () => void) => {
  done()
  file.pocData = null
  file.pocDataAll = null
  UploadPocListRef.value!.clearSelection()
}

function UploadPocToDB() {

  ElMessageBox.confirm(
      '确定导入选中的poc到数据库吗？',
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
  )
      .then(async () => {
        let msg = await PocsUpload(pocFileSelection.value)
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
        UploadClose()
        await pocList()
      })

      .catch(() => {})
}


function RemovePocFile() {

  ElMessageBox.confirm(
      '确定删除选中的本地poc文件吗？',
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
  )
      .then(async () => {
        let msg = await DelPocFile(pocFileSelection.value)
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
        UploadShow.value = false
        UploadPocListRef.value!.clearSelection()
        await pocList2()
      })

      .catch(() => {})
}




function pocDel(){
  if (pocSelection.value.length > 0){
    ElMessageBox.confirm(
        '确定删除选中的吗？',
        '提示',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
        }
    )
        .then(async () => {
          let msg = await DelPoc(pocSelection.value)
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
          await pocList()
        })

        .catch(() => {})
  }

}


const table = reactive({
  tableData: [] as Array<publicCode.Poc>,
  tableDataAll: [] as Array<publicCode.Poc>,
})


const file = reactive({
  pocDataAll: [] as Array<publicCode.Poc>,
  pocData: [] as Array<publicCode.Poc>,
})



const pocManageTablePage = ref(1)
const pocManageTablePageSize = ref(100)

function pocManageTableHandleSizeChange(){
  table.tableData = JSON.parse(JSON.stringify(table.tableDataAll))
  table.tableData = table.tableData.splice((pocManageTablePage.value - 1) * pocManageTablePageSize.value, pocManageTablePageSize.value)
}

function pocManageTableHandleCurrentChange(){
  table.tableData = JSON.parse(JSON.stringify(table.tableDataAll))
  table.tableData = table.tableData.splice((pocManageTablePage.value - 1) * pocManageTablePageSize.value, pocManageTablePageSize.value)
}




async function pocList() {
  // table.tableData = await PocList(input.value);
  table.tableDataAll = await PocList(input.value)
  table.tableData = JSON.parse(JSON.stringify(table.tableDataAll))
  table.tableData = table.tableData.splice((pocManageTablePage.value - 1) * pocManageTablePageSize.value, pocManageTablePageSize.value)

}

async function pocList2() {
  UploadShow.value = true
  file.pocDataAll = await PocList2(inputUpload.value)
  file.pocData = file.pocDataAll
  if (file.pocDataAll != null && table.tableDataAll != null){
    for(let i=0;i<file.pocDataAll.length;i++){
      for(let j=0;j<table.tableDataAll.length;j++){
        if (file.pocDataAll[i].uuid == table.tableDataAll[j].uuid){
          file.pocData.splice(i,1)
        }
      }
    }
  }

}

async function pocListFile() {
    //table.tableData = await LocalList(input.value)
  table.tableDataAll = await LocalList(input.value)
  table.tableData = JSON.parse(JSON.stringify(table.tableDataAll))
  table.tableData = table.tableData.splice((pocManageTablePage.value - 1) * pocManageTablePageSize.value, pocManageTablePageSize.value)
}


function loadingList(){
  pocListFile()
  pocList()
}

loadingList()


function downLoadPocs(){
  if (pocSelection.value.length > 0){
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
          let msg = await PocsDownload(pocSelection.value)
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
          multipleTableRef.value!.clearSelection()
        })

        .catch(() => {})
  }
}



//-----------------------  Main END  ---------------------------



//-----------------------  addPoc START  ---------------------------
const addActive = ref(0)
const addStep1 = reactive({
  show: true
})

const addStep2 = reactive({
  show: false
})


const addPocRuleFormRef = ref<FormInstance>()

const checkName = (rule: any, value: any, callback: any) => {
  if (value === '') {
    callback(new Error('请输入漏洞名称！'))
  } else {
    if (form.name !== '') {
      if (!addPocRuleFormRef.value) return
      addPocRuleFormRef.value.validateField('name', () => null)
    }
    callback()
  }
}

const checkCMS = (rule: any, value: any, callback: any) => {
  if (value === '') {
    callback(new Error('请输入漏洞名称！'))
  } else {
    if (form.cms !== '') {
      if (!addPocRuleFormRef.value) return
      addPocRuleFormRef.value.validateField('cms', () => null)
    }
    callback()
  }
}

interface  CMS {
  label: string
  value: string
}

const cms = ref<CMS[]> ([])
const cmsSearch= (queryString: string, cb: any) => {
  const results = queryString
      ? cms.value.filter(createFilter(queryString))
      : cms.value
  // call callback function to return suggestions
  cb(results)
}

const createFilter = (queryString: string) => {
  return (restaurant: CMS) => {
    return (
        restaurant.label.toLowerCase().indexOf(queryString.toLowerCase()) === 0
    )
  }
}

onMounted(async () => {
  let cmsName = await FingerprintGetName("")
  for(let i=0; i<cmsName.length; i++){
    cms.value.push({label: cmsName[i], value: cmsName[i]})
  }

})


const addPocRules = reactive<FormRules<typeof form>>({
  name: [{validator: checkName, trigger: 'blur'}],
  cms: [{validator: checkCMS, trigger: 'blur'}],
})

const form = reactive({
  uuid: "",
  name: "",
  hunter: "",
  fofa: "",
  cms: "",
  description: "",
  optionValue: null,
  needData: [
  ],
  request: [{
    pocString: "",
    status: "",
    check: "",
    print: "",
  },
    {
      pocString: "",
      status: "",
      check: "",
      print: "",
    },
    {
      pocString: "",
      status: "",
      check: "",
      print: "",
    },
    {
      pocString: "",
      status: "",
      check: "",
      print: "",
    },
    {
      pocString: "",
      status: "",
      check: "",
      print: "",
    },
    ]
} as publicCode.Poc)



const showAdd = ref(false)
const options = [
  {
    value: 0,
    label: '0'
  },
  {
    value: 1,
    label: '1',
  },
  {
    value: 2,
    label: '2',
  },
  {
    value: 3,
    label: '3',
  },
  {
    value: 4,
    label: '4',
  },
  {
    value: 5,
    label: '5',
  },
]





const nextStep = () => {

  if (form.name==""||form.cms==""){
    ElMessage({
      message: '漏洞名称和CMS必填！',
      type: 'error',
    })

  }else {
    addActive.value++

    if (addActive.value === 0) {
      addStep1.show = true
      addStep2.show = false
    }

    if (addActive.value === 1) {
      addStep1.show = false
      addStep2.show = true
    }

    if  (addActive.value >= 2) {
      for(let i=0; i<5; i++){
        if(i === addActive.value-2){
          pocUIValue[i].show =true
          if (i != times.value-1) pocUIValue[i].buttonShow1 = true
        }else {
          pocUIValue[i].show =false
          pocUIValue[i].buttonShow1 = false
        }
      }
      addStep2.show = false
    }
  }
}

const lastStep = () => {
  addActive.value--

  if (addActive.value === 0) {
    addStep1.show = true
    addStep2.show = false

  }

  if (addActive.value === 1) {
    addStep1.show = false
    addStep2.show = true
    pocUIValue[0].show =false
    pocUIValue[0].buttonShow1 =false
  }

  if  (addActive.value >= 2) {
    for(let i=0; i<5; i++){
      if(i === addActive.value-2){
        pocUIValue[i].show =true
        if (i != times.value-1) pocUIValue[i].buttonShow1 = true
      }else {
        pocUIValue[i].show =false
        pocUIValue[i].buttonShow1 = false
      }
    }
    addStep2.show = false
  }

}


const addPocCancel = () => {
  showAdd.value = false
  location.reload();
}


const pocUI = ref([]);
const times = ref(0)
const checkOption = reactive({
  value: true
})





const pocsUIChange = () => {

  if (times.value != 0) pocUIValue[times.value-1].buttonShow2 = false

  if (times.value < form.optionValue) {
    for (let i = times.value; i< form.optionValue; i++ ){
      const p = {id: i+1, name: "request"+i+1};
      pocUI.value.push(p);
    }
  }
  if (times.value > form.optionValue) {
    for (let i = times.value; i>=form.optionValue; i-- ){
      pocUI.value.splice(i,1);
      if (i != form.optionValue){
        form.request[i-1].pocString = ""
        form.request[i-1].status = ""
        form.request[i-1].check = ""
        form.request[i-1].print = ""
      }
    }
  }
  times.value = form.optionValue
  if(times.value != 0){
    checkOption.value = false
  }else {
    checkOption.value = true
  }

  pocUIValue[times.value-1].buttonShow1 = false
  pocUIValue[times.value-1].buttonShow2 = true

}




const addData = reactive({
  value: [
    {
      label: "request.url.0",
      value: "request.url.0"
    },
    {
      label: "编码(目前支持在POST的body中)",
      value: "{{hex_decode(\"\")}}"
    },
    {
      label: "Content-Length最大值(目前支持&&运算)",
      value: "{{Content-LengthMax(\"\")}}"
    },
    {
      label: "Content-Length最小值(目前支持&&运算)",
      value: "{{Content-LengthMin(\"\")}}"
    },
    {
      label: "响应时间(目前支持&&运算)",
      value: "{{Time(\"\")}}"
    },

  ],
  addNum: 1
})

const pocUIValue = reactive([
  {
    title: "Request1",
    show: false,
    buttonShow1: false,
    buttonShow2: false,
    addValue: "",
    settingShow: false,
    settingName: "",
    settingValue: "",
    settingAddShow: false,
    NeedSelect: false,
    NeedKey: [
      {
        label: "input",
        value: "input",
        disabled: false
      },
      { label: "request1.body",
        value: "request1.body"
      },
      { label: "request1.header",
        value: "request1.header"
      },
      ],
  },
  {
    title: "Request2",
    show: false,
    buttonShow1: false,
    buttonShow2: false,
    addValue: "",
    settingShow: false,
    settingName: "",
    settingValue: "",
    settingAddShow: false,
    NeedSelect: false,
    NeedKey: [
      {
        label: "input",
        value: "input",
        disabled: false
      },
      { label: "request1.body",
        value: "request1.body"
      },
      { label: "request1.header",
        value: "request1.header"
      },
      { label: "request2.body",
        value: "request2.body"
      },
      { label: "request2.header",
        value: "request2.header"
      },
    ],
  },
  {
    title: "Request3",
    show: false,
    buttonShow1: false,
    buttonShow2: false,
    addValue: "",
    settingShow: false,
    settingName: "",
    settingValue: "",
    settingAddShow: false,
    NeedSelect: false,
    NeedKey: [
      {
        label: "input",
        value: "input",
        disabled: false
      },
      { label: "request1.body",
        value: "request1.body"
      },
      { label: "request1.header",
        value: "request1.header"
      },
      { label: "request2.body",
        value: "request2.body"
      },
      { label: "request2.header",
        value: "request2.header"
      },
      { label: "request3.body",
        value: "request3.body"
      },
      { label: "request3.header",
        value: "request3.header"
      },
    ],
  },
  {
    title: "Request4",
    show: false,
    buttonShow1: false,
    buttonShow2: false,
    addValue: "",
    settingShow: false,
    settingName: "",
    settingValue: "",
    settingAddShow: false,
    NeedSelect: false,
    NeedKey: [
      {
        label: "input",
        value: "input",
        disabled: false
      },
      { label: "request1.body",
        value: "request1.body"
      },
      { label: "request1.header",
        value: "request1.header"
      },
      { label: "request2.body",
        value: "request2.body"
      },
      { label: "request2.header",
        value: "request2.header"
      },
      { label: "request3.body",
        value: "request3.body"
      },
      { label: "request3.header",
        value: "request3.header"
      },
      { label: "request4.body",
        value: "request4.body"
      },
      { label: "request4.header",
        value: "request4.header"
      },
    ],
  },
  {
    title: "Request5",
    show: false,
    buttonShow1: false,
    buttonShow2: false,
    addValue: "",
    settingShow: false,
    settingName: "",
    settingValue: "",
    settingAddShow: false,
    NeedSelect: false,
    NeedKey: [
      {
        label: "input",
        value: "input",
        disabled: false
      },
      { label: "request1.body",
        value: "request1.body"
      },
      { label: "request1.header",
        value: "request1.header"
      },
      { label: "request2.body",
        value: "request2.body"
      },
      { label: "request2.header",
        value: "request2.header"
      },
      { label: "request3.body",
        value: "request3.body"
      },
      { label: "request3.header",
        value: "request3.header"
      },
      { label: "request4.body",
        value: "request4.body"
      },
      { label: "request4.header",
        value: "request4.header"
      },
      { label: "request5.body",
        value: "request5.body"
      },
      { label: "request5.header",
        value: "request5.header"
      },
    ],
  },

])


let str
function cursorPosition_Refresh1(id){

  str =  form.request[id-1].pocString
}

function cursorPosition_Refresh2(id){

  str = form.request[id-1].print
}

function cursorPosition_Refresh3(id){

  str = form.request[id-1].check
}


function appendValue(id) {
  let index= selectionStart
  if (str ===  form.request[id-1].pocString && str != "" && pocUIValue[id-1].addValue != "")  {
    if (pocUIValue[id-1].addValue.includes("request"+id)){
      ElMessage({
        message: '当前返回包的值只允许插入在输出值！',
        type: 'warning',
      })
    }else {
      if(pocUIValue[id-1].addValue.includes("hex_decode")) {
        form.request[id-1].pocString = str.slice(0, index) + pocUIValue[id-1].addValue + str.slice(index);
      }else {
        form.request[id-1].pocString = str.slice(0, index) + "~" + pocUIValue[id-1].addValue + "~" + str.slice(index);
        if(pocUIValue[id-1].addValue.includes("url")) {
          form.needData = form.needData.filter(data => data.label != "request.url.0")
          form.needData.push({label:pocUIValue[id-1].addValue, value:""})
        }
      }
    }
  }
  if (str ===  form.request[id-1].print &&  pocUIValue[id-1].addValue != "" && str != "")  {
    form.request[id-1].print = str.slice(0, index) + "~" + pocUIValue[id-1].addValue + "~" + str.slice(index);
    if(pocUIValue[id-1].addValue.includes("url")) {
      form.needData = form.needData.filter(data => data.label != "request.url.0")
      form.needData.push({label:pocUIValue[id-1].addValue, value:""})
    }
  }
  if (str ===  form.request[id-1].check &&  pocUIValue[id-1].addValue != "")  {
    if(pocUIValue[id-1].addValue.includes("input")){
      form.request[id-1].check = str.slice(0, index) + "~" + pocUIValue[id-1].addValue + "~" + str.slice(index);
    }else if (pocUIValue[id-1].addValue.includes("Content-Length")){
      form.request[id-1].check = str.slice(0, index) + pocUIValue[id-1].addValue + str.slice(index);
    } else if (pocUIValue[id-1].addValue.includes("Time")){
      form.request[id-1].check = str.slice(0, index) + pocUIValue[id-1].addValue + str.slice(index);
    } else {
      ElMessage({
        message: '判断值这只允许插入input、Content-Length和Time！',
        type: 'warning',
      })
    }

  }

}

function addPoc_Reset(id){

  ElMessageBox.confirm(
      '确定重置吗？',
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
  )
      .then(() => {
        form.request[id-1].pocString = form.request[id-1].pocString.replace(/~.*~/,"")
        form.request[id-1].print = form.request[id-1].print.replace(/~.*~/,"")
        ElMessage({
          type: 'success',
          message: '重置成功！',
        })
      })
      .catch(() => {
      })
}

function addPoc_Setting(id) {
  pocUIValue[id-1].settingShow = true
}

function addPoc_SettingCancel(id) {
  pocUIValue[id-1].settingShow = false
}

function addPoc_SettingAdd(id) {
  pocUIValue[id-1].settingAddShow = true
}

function addPoc_SettingAddCancel(id) {
  pocUIValue[id-1].settingAddShow = false
}


function symbolReplace(id) {
  let start = selectionStart
  let end = selectionEnd
  let S = pocUIValue[id-1].settingValue
  pocUIValue[id-1].settingValue = S.slice(0, start) + "~" + S.slice(end);

}

function addPoc_SettingAddSave(id) {

  pocUIValue[id-1].NeedSelect = false
  pocUIValue[id-1].settingAddShow = false
  if (pocUIValue[id-1].settingName.includes("input")){
    form.needData = form.needData.filter(data => data.label != "input")
    form.needData.push({label:"input", value:pocUIValue[id-1].settingValue})
    addData.value =  addData.value.filter(data => data.label != "input")
    addData.value.push({label:"input", value:"input"})
    for (let i = 0; i < form.optionValue; i++){
      pocUIValue[i].NeedKey[0].disabled=true
    }
  }else {
    form.needData = form.needData.filter(data => data.label != pocUIValue[id-1].settingName)
    form.needData.push({label:pocUIValue[id-1].settingName, value:pocUIValue[id-1].settingValue})
    addData.value.push({label:pocUIValue[id-1].settingName, value:pocUIValue[id-1].settingName})
    addData.addNum++
  }
  pocUIValue[id-1].settingName = null
  pocUIValue[id-1].settingValue = null

}


interface Need{
  label:string,
  value:string
}

const needDataTableRef = ref<InstanceType<typeof ElTable>>()
const needDataSelection = ref<Need[]>([])

const needDataChange = (val: Need[]) => {
  needDataSelection.value = val
}

function  needDataDel(id) {

  needDataSelection.value.forEach((data)=>{
    for (let i=0; i<form.needData.length; i++){
      if (needDataSelection.value[0].label.includes("url")) {
        form.needData.splice(i,1)
        form.request[id-1].pocString = form.request[id-1].pocString.replace("~"+data.label+"~","")
        form.request[id-1].print = form.request[id-1].print.replace("~"+data.label+"~","")
      }else {
        if(form.needData[i].label === data.label){
          form.needData.splice(i,1)
          addData.value.splice(i,1)
          form.request[id-1].pocString = form.request[id-1].pocString.replace("~"+data.label+"~","")
          form.request[id-1].print = form.request[id-1].print.replace("~"+data.label+"~","")
          if(data.label.includes("input")){
            for (let i = 0; i < form.optionValue; i++){
              pocUIValue[i].NeedKey[0].disabled=false
              form.request[id-1].pocString = form.request[id-1].pocString.replace(/~input~/,"")
              form.request[id-1].print = form.request[id-1].print.replace(/~input~/,"")
            }
          }
        }
      }
    }
  } )
}


function needDataEdit(id) {

  if (needDataSelection.value.length===1){
    for (let i=0; i<form.needData.length; i++) {
      if (needDataSelection.value[0].label.includes("url")) {
        ElMessage({
          message: '不允许编辑！',
          type: 'warning',
        })
      }else {
        if(form.needData[i].label === needDataSelection.value[0].label){
          pocUIValue[id-1].settingAddShow = true
          pocUIValue[id-1].NeedSelect = true
          pocUIValue[id-1].settingName = form.needData[i].label
          pocUIValue[id-1].settingValue = form.needData[i].value
        }
      }
    }
  }else if(needDataSelection.value.length>1) {
    ElMessage({
      message: '只允许同时编辑一个！',
      type: 'warning',
    })
  }

}


function  addPoc_Save(){

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
        let msg = await AddPoc(form)
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

      })
      .catch(() => {
      })

}

//--------  addPoc end  ---------------------

//-----------------------  editPoc START  ---------------------------

function editPoc(){
  if (pocSelection.value.length === 1){
    form.uuid = pocSelection.value[0].uuid
    form.name = pocSelection.value[0].name
    form.cms = pocSelection.value[0].cms
    form.hunter = pocSelection.value[0].hunter
    form.fofa = pocSelection.value[0].fofa
    form.description = pocSelection.value[0].description
    form.optionValue = pocSelection.value[0].optionValue
    form.needData = pocSelection.value[0].needData
    for (let i = 0; i < pocSelection.value[0].request.length; i++){
      form.request[i] = JSON.parse(JSON.stringify(pocSelection.value[0].request[i]))
    }

    if (form.needData != null ){
      form.needData.forEach((need)=>{
        if (!need.label.includes("url")){
          let label = need.label
          addData.value.push({label:label, value:label})
          addData.addNum++
        }

      })
    }

    showAdd.value = true
    pocsUIChange()

  }else if (pocSelection.value.length > 1) {
    ElMessage({
      type: 'error',
      message: '只允许编辑一个',
    })
  }
}


//-----------------------  editPOC END  ---------------------------


function AddOREditCancel() {
  form.uuid = ""
  form.name = ""
  form.cms = ""
  form.hunter = ""
  form.fofa = ""
  form.description = ""
  form.optionValue = 0
  form.needData = []
  for (let i = 0; i < 5; i++){
    form.request[i].print = ""
    form.request[i].status = ""
    form.request[i].check = ""
    form.request[i].pocString = ""
  }
  addData.value = [
    {
      label: "request.url.0",
      value: "request.url.0"
    },
    {
      label: "编码(目前支持在POST的body中)",
      value: "{{hex_decode(\"\")}}"
    },
    {
      label: "Content-Length最大值(目前支持&&运算)",
      value: "{{Content-LengthMax(\"\")}}"
    },
    {
      label: "Content-Length最小值(目前支持&&运算)",
      value: "{{Content-LengthMin(\"\")}}"
    },
    {
      label: "响应时间(目前支持&&运算)",
      value: "{{Time(\"\")}}"
    },
  ]
  addData.addNum = 1

  addActive.value = 0
  addStep1.show = true
  addStep2.show = false
  checkOption.value = true
  pocUI.value = []
  times.value = 0
  for(let i = 0; i < 5; i++){
    pocUIValue[i].show = false
    pocUIValue[i].buttonShow1 = false
    pocUIValue[i].buttonShow2 = false
    pocUIValue[i].addValue = ""
    pocUIValue[i].settingShow = false
    pocUIValue[i].settingName = ""
    pocUIValue[i].settingValue = ""
    pocUIValue[i].settingAddShow = false
    pocUIValue[i].NeedSelect = false
  }
  showAdd.value = false
}







</script>

<style>

/*-----------------------  MAIN START  ---------------------------*/

.demo-radius {
  height: 72%;
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

.demo-radius3 {
  width: 94%;
  border: 2px solid var(--el-border-color);
  border-radius: 4px;
  margin-top: 2%;
  margin-left: 3%;
}


.topUI {
  margin-top: 2%;
  margin-bottom: 0;
  margin-left: 3%;
  display: flex;
}


/*-----------------------  MAIN END  ---------------------------*/



/*-----------------------  AddPOC START  ---------------------------*/

.addPoc {
  height: 75%;
  width: 94%;
  border: 2px solid var(--el-border-color);
  border-radius: 4px;
  margin-left: 3%;
}

.step{
  margin-top: 0;
  margin-left: 3%;
  margin-right: 2.5%;
}

.form{
  margin-top: 5%;
  margin-left: 5%;

}

/*-----------------------  AddPOC END  ---------------------------*/


</style>