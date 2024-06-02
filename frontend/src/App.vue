<script lang="ts" setup>
import HelloWorld from './components/HelloWorld.vue'
import {BackgroundSetting} from "../wailsjs/go/main/App";

let backgroundImage
async function backgroundImagUrl() {
  backgroundImage = await BackgroundSetting()

  let bodyElement = document.body;

// 修改元素的样式属性
  bodyElement.style.backgroundImage = "url(" + backgroundImage+ ")";
  bodyElement.style.backgroundSize = "cover";
  bodyElement.style.backgroundRepeat = "no-repeat";
  bodyElement.style.backgroundPosition = "left center";
  bodyElement.style.opacity = "0.8";
}

backgroundImagUrl()



</script>

<template>
  <HelloWorld/>
</template>

<style>

body {
  background-image: url("background-image.png");
  background-size: cover;
  background-repeat: no-repeat;
  background-position: left center;
  opacity: 0.8;
}

.labelSize > .el-form-item__label{
  font-size: large;
  color: #38baff
}

</style>
