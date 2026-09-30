<template>
  <div v-if="article" class="detail">
    <h2 class="title">{{ article.title }}</h2>
    <div class="meta">
      <span>{{ formatTime(article.createdAt) }}</span>
      <span class="meta-dot">·</span>
      <span>Free Templates by My Little Yard</span>
    </div>
    <div class="tags">
      <el-tag v-for="t in article.tags" :key="t" size="small" class="tag-dark">{{ t }}</el-tag>
    </div>
    <div class="content" v-html="renderContent"></div>

    <div v-if="article.attachments && article.attachments.length" class="attachments">
      <div class="attachments-title">Attachments</div>
      <a
        v-for="a in article.attachments"
        :key="a.fileId"
        :href="`http://localhost:8080/api/files/${a.fileId}/download`"
        class="attachment-link"
      >
        📎 {{ a.filename }}
      </a>
    </div>

    <div class="detail-footer">
      <a class="action-link" @click="$router.push(`/edit/${article.id}`)">Edit this article</a>
      <span class="meta-dot">|</span>
      <a class="action-link" @click="$router.push('/')">Back to list</a>
    </div>
  </div>
  <div v-else class="empty-state">
    <p>文章不存在</p>
  </div>
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
  background-color: #141c33;
  border: 1px solid #1e293b;
  border-radius: 4px;
  padding: 24px;
}

.title {
  font-size: 24px;
  color: #3a82ff;
  margin-bottom: 12px;
}

.meta {
  font-size: 12px;
  color: #64748b;
  margin-bottom: 12px;
}

.meta-dot {
  margin: 0 6px;
  color: #334155;
}

.tags {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
  margin-bottom: 20px;
}

.tag-dark {
  background-color: #1e293b !important;
  border-color: #2a3654 !important;
  color: #94a3b8 !important;
}

.content {
  line-height: 1.8;
  color: #cbd5e1;
  font-size: 15px;
  margin-bottom: 24px;
}

.content :deep(img) {
  max-width: 100%;
  border-radius: 4px;
  margin: 12px 0;
  display: block;
}

.attachments {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #24304a;
}

.attachments-title {
  font-weight: bold;
  color: #3a82ff;
  font-size: 14px;
  margin-bottom: 8px;
}

.attachment-link {
  display: block;
  color: #3a82ff;
  text-decoration: none;
  margin-bottom: 6px;
  font-size: 14px;
}

.attachment-link:hover {
  color: #6aa1ff;
}

.detail-footer {
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px solid #24304a;
  font-size: 13px;
}

.action-link {
  color: #3a82ff;
  cursor: pointer;
  text-decoration: underline;
}

.action-link:hover {
  color: #6aa1ff;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #64748b;
  font-size: 16px;
}
</style>
