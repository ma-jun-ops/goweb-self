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
  background: #3f6b4f !important;
  border-color: #3f6b4f !important;
  box-shadow: 0 10px 24px rgba(63, 107, 79, 0.28);
  z-index: 999;
}

.chat-window {
  position: fixed;
  right: 24px;
  bottom: 24px;
  width: 380px;
  height: 560px;
  background: #fffaf3;
  border-radius: 18px;
  box-shadow: 0 16px 40px rgba(80, 58, 32, 0.16);
  overflow: hidden;
  z-index: 999;
  display: flex;
  flex-direction: column;
  border: 1px solid rgba(70, 58, 42, 0.1);
}

.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #3f6b4f;
  color: #fffaf3;
  font-weight: 600;
}

.close-btn {
  color: #fff;
  font-size: 18px;
}

@media (max-width: 480px) {
  .chat-btn {
    right: 16px;
    bottom: 16px;
    width: 50px;
    height: 50px;
  }

  .chat-window {
    right: 12px;
    left: 12px;
    bottom: 12px;
    width: auto;
    height: 70vh;
  }
}
</style>
