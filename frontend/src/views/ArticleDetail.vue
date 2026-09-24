<template>
  <div v-if="article" class="detail">
    <h2 class="title">{{ article.title }}</h2>
    <div class="meta">
      <el-tag v-for="t in article.tags" :key="t" style="margin-right: 4px">{{ t }}</el-tag>
      <span class="time">{{ formatTime(article.createdAt) }}</span>
    </div>
    <div class="content" v-html="renderContent"></div>

    <div v-if="article.attachments && article.attachments.length" class="attachments">
      <div class="attachments-title">附件</div>
      <a
        v-for="a in article.attachments"
        :key="a.fileId"
        :href="`http://localhost:8080/api/files/${a.fileId}/download`"
        class="attachment-link"
      >
        📎 {{ a.filename }}
      </a>
    </div>
  </div>
  <el-empty v-else description="文章不存在" />
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { getArticle } from '../api/article'

const route = useRoute()
const article = ref(null)

onMounted(async () => {
  const res = await getArticle(route.params.id)
  article.value = res.data.article
})

const renderContent = computed(() => {
  if (!article.value) return ''
  const content = article.value.content || ''
  const parts = content.split(/\[IMG:([0-9a-f]{24})\]/)
  let html = ''
  for (let i = 0; i < parts.length; i++) {
    if (i % 2 === 0) {
      html += escapeHtml(parts[i]).replace(/\n/g, '<br>')
    } else {
      html += `<img src="http://localhost:8080/api/files/${parts[i]}/image" style="max-width: 100%" />`
    }
  }
  return html
})

const escapeHtml = (s) => s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')

const formatTime = (s) => new Date(s).toLocaleString()
</script>

<style scoped>
.detail {
  background-color: #fff;
}

.title {
  margin-bottom: 12px;
}

.meta {
  margin-bottom: 16px;
}

.time {
  margin-left: 8px;
  color: #999;
  font-size: 14px;
}

.content {
  line-height: 1.8;
  margin-bottom: 24px;
}

.content :deep(img) {
  max-width: 100%;
  border-radius: 4px;
  margin: 8px 0;
  display: block;
}

.attachments {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #eee;
}

.attachments-title {
  font-weight: bold;
  margin-bottom: 8px;
}

.attachment-link {
  display: block;
  color: #409eff;
  text-decoration: none;
  margin-bottom: 6px;
}
</style>
