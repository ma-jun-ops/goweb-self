<template>
  <div>
    <div v-if="showChat" class="chat-window">
      <div class="chat-header">
        <span>博客助手</span>
        <el-button text @click="showChat = false" class="close-btn">✕</el-button>
      </div>
      <iframe
        src="http://localhost/agent/D6ezLa0EoN1UZPrI"
        style="width: 100%; height: 100%; border: none"
        allow="microphone;clipboard-write"
      ></iframe>
    </div>
    <el-button
      v-else
      class="chat-btn"
      type="primary"
      circle
      size="large"
      @click="showChat = true"
    >
      💬
    </el-button>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const showChat = ref(false)
</script>

<style scoped>
.chat-btn {
  position: fixed;
  right: 24px;
  bottom: 24px;
  width: 56px;
  height: 56px;
  font-size: 24px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 999;
}

.chat-window {
  position: fixed;
  right: 24px;
  bottom: 24px;
  width: 380px;
  height: 560px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
  overflow: hidden;
  z-index: 999;
  display: flex;
  flex-direction: column;
}

.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #409eff;
  color: #fff;
  font-weight: bold;
}

.close-btn {
  color: #fff;
  font-size: 18px;
}
</style>
