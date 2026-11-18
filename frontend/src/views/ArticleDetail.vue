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
        :href="`/api/files/${a.fileId}/download`"
        class="attachment-link"
      >
        📎 {{ a.filename }}
      </a>
    </div>

    <div class="detail-footer">
      <button class="pill-btn" type="button" @click="$router.push('/articles')">返回列表</button>
      <button class="pill-btn primary" type="button" @click="$router.push(`/edit/${article.id}`)">编辑文章</button>
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
      html += `<img src="/api/files/${parts[i]}/image" style="max-width: 100%" />`
    }
  }
  return html
})

const escapeHtml = (s) => s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
const formatTime = (s) => new Date(s).toLocaleString()
</script>

<style scoped>
.detail {
  background: rgba(255, 250, 243, 0.86);
  border: 1px solid var(--paper-edge);
  border-radius: 22px;
  padding: 32px 34px;
  box-shadow: 0 12px 32px rgba(80, 58, 32, 0.06);
}

.title {
  font-family: var(--font-serif);
  font-size: 34px;
  line-height: 1.3;
  color: var(--ink);
  margin-bottom: 12px;
}

.meta {
  font-size: 13px;
  color: var(--faint);
  margin-bottom: 14px;
}

.meta-dot {
  margin: 0 6px;
  color: var(--faint);
}

.tags {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
  margin-bottom: 24px;
}

.tag-dark {
  background: rgba(63, 107, 79, 0.1) !important;
  border-color: transparent !important;
  color: var(--leaf-deep) !important;
}

.content {
  line-height: 1.9;
  color: var(--ink);
  font-size: 16px;
  margin-bottom: 24px;
  word-break: break-word;
}

.content :deep(img) {
  max-width: 100%;
  border-radius: 12px;
  margin: 12px 0;
  display: block;
}

.attachments {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid var(--paper-edge);
}

.attachments-title {
  font-weight: 600;
  color: var(--leaf-deep);
  font-size: 14px;
  margin-bottom: 8px;
}

.attachment-link {
  display: block;
  color: var(--leaf);
  text-decoration: none;
  margin-bottom: 6px;
  font-size: 14px;
}

.attachment-link:hover {
  color: var(--leaf-deep);
}

.detail-footer {
  margin-top: 20px;
  padding-top: 18px;
  border-top: 1px solid var(--paper-edge);
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 10px;
}

.pill-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  height: 36px;
  padding: 0 18px;
  border: 1px solid var(--paper-edge);
  background: transparent;
  color: var(--muted);
  border-radius: 999px;
  font-size: 13px;
  font-weight: 500;
  font-family: inherit;
  cursor: pointer;
  transition: color 0.18s, background 0.18s, border-color 0.18s;
}

.pill-btn:hover {
  color: var(--leaf-deep);
  border-color: rgba(63, 107, 79, 0.35);
  background: rgba(63, 107, 79, 0.06);
}

.pill-btn.primary {
  background: var(--leaf);
  border-color: var(--leaf);
  color: #fffaf3;
  font-weight: 600;
}

.pill-btn.primary:hover {
  background: var(--leaf-deep);
  border-color: var(--leaf-deep);
  color: #fffaf3;
}

.empty-state {
  text-align: center;
  padding: 48px 20px;
  color: var(--muted);
}

@media (max-width: 560px) {
  .detail {
    padding: 22px 18px;
  }

  .title {
    font-size: 28px;
  }

  .detail-footer {
    flex-direction: column;
    align-items: stretch;
  }

  .pill-btn {
    width: 100%;
  }
}
</style>
