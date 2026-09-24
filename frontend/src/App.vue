<template>
  <el-container class="app-container">
    <el-header class="header">
      <h1 class="logo" @click="$router.push('/')">我的博客</h1>
      <el-button type="primary" @click="$router.push('/edit')">写文章</el-button>
    </el-header>
    <el-main class="main">
      <router-view />
    </el-main>
    <ChatBot />
  </el-container>
</template>

<script setup>
import ChatBot from './components/ChatBot.vue'
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html,
body,
#app {
  height: 100%;
  background-color: #ffffff;
}

.app-container {
  min-height: 100%;
  background-color: #ffffff;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid #eee;
  background-color: #ffffff;
}

.logo {
  cursor: pointer;
  font-size: 20px;
}

.main {
  background-color: #ffffff;
}
</style>
