<template>
  <el-tabs tab-position="left" style="height: 92%"  class="paneLabelSize">
    <el-tab-pane label="扫描设置" >
      <div style="display: flex;margin-left: 3%">
        <el-form :model="SettingForm.proxy"   label-position="left" style="margin-top: 7.5%;width: 40%;margin-left: 3%" label-width="35%">
          <el-form-item label="代理状态"  class="labelSize" >
            <el-switch v-model="SettingForm.proxy.enable" />
          </el-form-item>
          <el-form-item label="代理模式"  class="labelSize" style="margin-top: 10%">
            <el-select v-model="SettingForm.proxy.mode" placeholder="please select" style="width: 50%;">
              <el-option label="HTTP" value="HTTP" />
              <el-option label="HTTPS" value="HTTPS" />
              <el-option label="SOCKS5" value="SOCKS5" />
              <el-option label="SOCKS4" value="SOCKS4" />
            </el-select>
          </el-form-item>
          <el-form-item label="代理地址"  class="labelSize" style="margin-top: 10%">
            <el-input v-model="SettingForm.proxy.host" style="width: 100%"/>
          </el-form-item>
          <el-form-item label="代理端口"  class="labelSize" style="margin-top: 10%">
            <el-input v-model="SettingForm.proxy.port" style="width: 100%"/>
          </el-form-item>
          <el-form-item label="代理用户名"  class="labelSize" style="margin-top: 10%">
            <el-input v-model="SettingForm.proxy.username" style="width: 100%"/>
          </el-form-item>
          <el-form-item label="代理密码"  class="labelSize" style="margin-top: 10%">
            <el-input v-model="SettingForm.proxy.password" style="width: 100%"/>
          </el-form-item>
        </el-form>
        <el-form :model="SettingForm.scan"   label-position="left" style="margin-top: 10%;width: 40%;margin-left: 7%" label-width="32%">
          <el-form-item label="扫描参数"  class="labelSize">
          </el-form-item>
          <el-form-item label="扫描超时(s)"  class="labelSize" style="margin-top: 10%">
            <el-input v-model.number="SettingForm.scan.timeout" style="width: 100%"/>
          </el-form-item>
          <el-form-item label="扫描线程"  class="labelSize" style="margin-top: 10%">
            <el-input v-model.number="SettingForm.scan.threadNum" style="width: 100%" />
          </el-form-item>
          <div style="margin-left: 35%;margin-top: 10%">
            <el-button  size="large" :icon="CloseBold" @click="scanSettingCancel">取&emsp;消</el-button>
            <el-button size="large" :icon="Select"  type="primary" @click="settingSave(0)">保&emsp;存</el-button>
          </div>
        </el-form>
      </div>
    </el-tab-pane>
    <el-tab-pane label="数据库配置">
      <div style="display: flex;margin-left: 3%">
        <el-form :model="SettingForm.db"   label-position="left" style="margin-top: 10%;width: 40%;margin-left: 3%">
          <el-form-item label="是否使用公共数据库"  class="labelSize"  >
            <el-switch v-model="SettingForm.db.mode" />
            <el-button  size="large" :icon="Refresh" @click="InitializeDB" type="danger" style="margin-bottom: 2%;margin-left: 20%">初始化</el-button>
          </el-form-item>
          <el-form-item label="公共数据库地址"  class="labelSize"  style="margin-top: 10%">
            <el-input v-model="SettingForm.db.host" style="width: 100%"/>
          </el-form-item>
          <el-form-item label="公共数据库端口"  class="labelSize"  style="margin-top: 10%">
            <el-input v-model="SettingForm.db.port" style="width: 100%"/>
          </el-form-item>
          <el-form-item label="公共数据库用户"  class="labelSize"  style="margin-top: 10%">
            <el-input v-model="SettingForm.db.username" style="width: 100%"/>
          </el-form-item>
          <el-form-item label="公共数据库密码"  class="labelSize"  style="margin-top: 10%">
            <el-input v-model="SettingForm.db.password" style="width: 100%" type="password" show-password/>
          </el-form-item>
<!--          <div style="margin-left: 35%;margin-top: 10%">-->
<!--            <el-button  size="large" :icon="CloseBold" @click="scanSettingCancel">取&emsp;消</el-button>-->
<!--            <el-button size="large" :icon="Select"  type="primary" @click="settingSave(1)">保&emsp;存</el-button>-->
<!--          </div>-->
        </el-form>
        <el-form :model="SettingForm.db"   label-position="left" style="margin-top: 10%;width: 40%;margin-left: 7%">
<!--          <el-button  size="large" :icon="Refresh" @click="InitializeDB" type="danger" style="margin-bottom: 2%;margin-right: 70%">初始化</el-button>-->
<!--          <el-form-item label="公共数据库地址"  class="labelSize" style="margin-top: 5%">-->
<!--            <el-input v-model="SettingForm.db.hostPublic" style="width: 100%"/>-->
<!--          </el-form-item>-->
<!--          <el-form-item label="公共数据库端口"  class="labelSize" style="margin-top: 10%">-->
<!--            <el-input v-model="SettingForm.db.portPublic" style="width: 100%"/>-->
<!--          </el-form-item>-->
<!--          <el-form-item label="公共数据库用户"  class="labelSize" style="margin-top: 10%">-->
<!--            <el-input v-model="SettingForm.db.usernamePublic" style="width: 100%"/>-->
<!--          </el-form-item>-->
<!--          <el-form-item label="公共数据库密码"  class="labelSize" style="margin-top: 10%">-->
<!--            <el-input v-model="SettingForm.db.passwordPublic" style="width: 100%" type="password" show-password/>-->
<!--          </el-form-item>-->
          <div style="margin-left: 15%;margin-top: 80%">
            <el-button  size="large" :icon="CloseBold" @click="scanSettingCancel">取&emsp;消</el-button>
            <el-button size="large" :icon="Select"  type="primary" @click="settingSave(1)">保&emsp;存</el-button>
          </div>
        </el-form>
      </div>
    </el-tab-pane>
    <el-tab-pane label="其他设置">
      <el-form :model="SettingForm.other"   label-position="left" style="margin-top: 10%;width: 40%;margin-left: 6%">
        <el-form-item label="是否使用本地预加载"  class="labelSize" style="margin-top: 10%">
          <el-switch v-model="SettingForm.other.localLoading" />
        </el-form-item>
<!--        <el-form-item label="是否只加载本地"  class="labelSize" style="margin-top: 10%">-->
<!--          <el-switch v-model="SettingForm.other.onlyLocalLoading" />-->
<!--        </el-form-item>-->
        <el-form-item label="设置背景图片"  class="labelSize" style="margin-top: 10%">
        <el-upload
            class="avatar-uploader"
            :auto-upload="false"
            :show-file-list="false"
            @change="backGroundChange"
        >

          <img v-if="image" :src="image" class="avatar" />
          <el-icon v-else class="avatar-uploader-icon"><Plus /></el-icon>
        </el-upload>
        </el-form-item>
<!--        <el-form-item label="设置背景图片"  class="labelSize"  style="margin-top: 10%">-->
<!--          <el-input v-model="SettingForm.other.backgroundImage" style="width: 100%" placeholder="输入远程图片地址"/>-->
<!--        </el-form-item>-->
      </el-form>
        <div style="margin-left: 35%;margin-top: 10%">
          <el-button  size="large" :icon="CloseBold" @click="scanSettingCancel">取&emsp;消</el-button>
          <el-button size="large" :icon="Select"  type="primary" @click="settingSave(0)">保&emsp;存</el-button>
        </div>

    </el-tab-pane>
  </el-tabs>
</template>

<script lang="ts" setup>

import {reactive, ref} from 'vue'
import {CloseBold, Select,Refresh,Plus} from '@element-plus/icons-vue'
import {ElMessage, ElMessageBox} from "element-plus";
import {
  BackgroundImageSave,
  BackgroundSetting,
  DBInitialize,
  SaveSetting,
  ShowSetting
} from "../../wailsjs/go/main/App";
import {publicCode} from "../../wailsjs/go/models";
import ScanSetting = publicCode.ScanSetting;
import ProxySetting = publicCode.ProxySetting;
import DBSetting = publicCode.DBSetting;
import AllSetting = publicCode.AllSetting;
import OtherSetting = publicCode.OtherSetting;
import Background = publicCode.Background;


const image = ref('')

const backgroundSetting = reactive({
  image: ''
} as Background)

async function backgroundImagUrl() {
  backgroundSetting.image = await BackgroundSetting()
}
backgroundImagUrl()

const SettingForm = reactive({
  dbSetting: {
    mode: false,
    host: '',
    port: '',
    username: '',
    password: '',
    // hostPublic: '',
    // portPublic: '',
    // usernamePublic: '',
    // passwordPublic: '',
  } as DBSetting,
  proxySetting: {
    host: '',
    port: '',
    username: '',
    password: '',
    enable: false,
    mode: '',
  } as ProxySetting,
  scanSetting: {
    timeout: 5,
    threadNum: 20,
  } as ScanSetting,
  otherSetting: {
    localLoading: false,
    // onlyLocalLoading: false,
  } as OtherSetting
} as AllSetting)


function scanSettingCancel() {
  location.reload();
}

async function settingSave(dbChange) {
  await BackgroundImageSave(backgroundSetting)
  let msg = await SaveSetting(SettingForm,dbChange)
  if (msg[0] === "1") {
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
  setTimeout(function () {
    location.reload();
  }, 800);
}


async function settingShow() {
  let  allSetting
  allSetting = await ShowSetting()
  SettingForm.db = allSetting.db
  SettingForm.proxy = allSetting.proxy
  SettingForm.scan = allSetting.scan
  SettingForm.other = allSetting.other
}

settingShow()


function InitializeDB(){
  ElMessageBox.confirm(
      '确定初始化吗？',
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
  )
      .then(async () => {
        let msg = await DBInitialize()
        console.log(msg)
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
      })
      .catch(() => {
      })
}

// const timeoutConvertToNumber = (value) =>{
//   SettingForm.scan.timeout = Number(value);
// }
//
// const threadConvertToNumber = (value) =>{
//   SettingForm.scan.threadNum = Number(value);
// }

const backGroundChange = (file) => {
  if (file.raw) {
    const reader = new FileReader()
    reader.onload = () => {
      backgroundSetting.image = reader.result.toString()
      image.value = backgroundSetting.image
    }
    reader.readAsDataURL(file.raw);
  }
}


</script>

<style scoped>
.avatar-uploader .avatar {
  width: 178px;
  height: 178px;
  display: block;
}
</style>


<style>
.paneLabelSize .el-tabs__item{
  font-size: large;
  color: #38baff;
}

.avatar-uploader .el-upload {
  border: 1px dashed var(--el-border-color);
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: var(--el-transition-duration-fast);
}

.avatar-uploader .el-upload:hover {
  border-color: var(--el-color-primary);
}

.el-icon.avatar-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 178px;
  height: 178px;
  text-align: center;
}
</style>