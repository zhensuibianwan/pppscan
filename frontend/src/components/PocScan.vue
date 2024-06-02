<template>

  <div class="topUI">
    <el-col :span="6">
      <el-input v-model="input"  style="width: 191%" size="large"  @keyup.enter="pocList()" clearable>
        <template #prepend>漏洞名称</template>
        <template #append><el-button :icon="Search" style="width: 75px;color: #409eff"  @click="pocList"/></template>
      </el-input>
    </el-col>
    <el-col :span="6" style="margin-left: 22%" >
      <el-input v-model="url"  style="width: 191%;margin-left: 10%;"  size="large" @keyup.enter="pocScan()" clearable>
        <template #prepend>漏洞地址</template>
        <template #append><el-button :icon="Aim" style="width: 75px;color: #ff4d51" size="large"  @click="pocScan"/></template>
      </el-input>
    </el-col>

  </div>
  <div style="height: 82%; display: flex;" class="background">
    <div class="demo-radiusPocScan" style="width: 46%;margin-left: 3%;margin-top: 2%" >
      <el-table
          ref="multipleTableRef"
          :data="table.tableData"
          style="width: 100%"
          @selection-change="handleSelectionChange"
          height="100%"
      >
        <el-table-column   type="selection" width="35" />
        <el-table-column   label="漏洞名称" width="500">
          <template #default="scope">
            <el-tooltip  placement="top" >
              <template #content>  {{ scope.row.description }} </template>
              {{ scope.row.name }}
            </el-tooltip>
          </template>
        </el-table-column>
      </el-table>
    </div>
    <el-input
        v-model="result"
        :rows="28"
        type="textarea"
        style="width: 46%;margin-left: 2%;margin-top: 2%"
        id="request"
        clearable
        size="large"
    />
<!--    <div class="demo-radius" style="width: 46%;margin-left: 2%;margin-top: 2%">-->

<!--    </div>-->
  </div>

</template>

<script lang="ts" setup>
import {Search, Aim} from '@element-plus/icons-vue';
import {reactive, ref} from 'vue'
import {publicCode} from "../../wailsjs/go/models";
import {LocalList, PocList, PocScan} from "../../wailsjs/go/main/App";
import {EventsOn} from "../../wailsjs/runtime";

const input = ref('')
const url = ref('')
const result = ref('')

const pocSelection = ref<publicCode.Poc[]>([])

const handleSelectionChange = (val: publicCode.Poc[]) => {

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

let time = 0

function pocScan(){
  if (pocSelection.value.length > 0){
    if (result.value === ""){
      time = 0
    }
    if(time != 0){
      result.value = result.value + "\n"
    }
    PocScan(pocSelection.value, url.value)
    time = 1
  }
}


EventsOn('PocScan', function (msg){
  result.value = result.value + msg
})


</script>

<style scoped>

.demo-radiusPocScan {

  border: 2px solid var(--el-border-color);
  border-radius: 4px;
}

.topUI {
  margin-top: 2%;
  margin-bottom: 0;
  margin-left: 3%;
  display: flex;
}


.background{
  background-size: cover;
  background-repeat: no-repeat;
  background-position: left center;
  opacity: 0.8;
}

</style>