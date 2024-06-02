<template>
  <el-menu
      class="el-menu-demo"
      mode="horizontal"
      router
  >
    <el-menu-item index="/setting"  ><el-icon size="56" ><Setting /></el-icon></el-menu-item>
    <el-menu-item v-for="(item,i) in data.navList" :key="i" :index="item.name" >
      {{ item.navItem }}
    </el-menu-item>
  </el-menu>
  <div class="h-6" />

  <router-view v-slot="{ Component }">
    <keep-alive exclude="FingerprintManage,PocManage">
      <component :is="Component"  />
    </keep-alive>
  </router-view>
</template>

<script lang="ts" setup>

const  data = {
      navList: [
        {name: '/', navItem: 'POC管理'},
        {name: '/fingerprintManage', navItem: '指纹管理'},
        {name: '/scan', navItem: '漏洞扫描'},
        {name: '/exploitation', navItem: '漏洞利用'},
        {name: '/batchScan', navItem: '多线程扫描'},
        {name: '/fingerprintScan', navItem: '指纹识别'},
        {name: '/other', navItem: '申明'},
      ]

    }

</script>
