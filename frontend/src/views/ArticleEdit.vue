<template>
  <div class="edit">
    <h2 class="edit-title">{{ id ? 'Edit Article' : 'Write New Article' }}</h2>
    <el-input
      v-model="form.title"
      placeholder="Article title"
      class="edit-input"
    />
    <el-input
      v-model="tagsText"
      placeholder="Tags, separated by commas"
      class="edit-input"
    />
    <el-input
      v-model="form.content"
      type="textarea"
      :rows="15"
      placeholder="Article content..."
      class="edit-textarea"
    />
    <div class="edit-actions">
      <input
        ref="fileInput"
        type="file"
        accept=".docx,.doc"
        style="display: none"
        @change="onFileChange"
      />
      <el-button class="btn-upload" @click="$refs.fileInput.click()">上传 Word</el-button>
      <div class="edit-actions-right">
        <el-button class="btn-back" @click="$router.back()">返回</el-button>
        <el-button class="btn-save" type="primary" @click="onSave">保存</el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { createArticle, updateArticle, getArticle, uploadFile } from '../api/article'

const route = useRoute()
const router = useRouter()

const form = ref({ title: '', content: '', tags: [], attachments: [] })
const tagsText = ref('')
const id = route.params.id
const fileInput = ref(null)

onMounted(async () => {
  if (id) {
    const res = await getArticle(id)
    form.value = res.data.article
    tagsText.value = (form.value.tags || []).join(',')
  }
})

const onFileChange = async (e) => {
  const file = e.target.files[0]
  if (!file) return
  const res = await uploadFile(file)
  form.value.content = res.data.text
  form.value.attachments = [{
    fileId: res.data.fileId,
    filename: res.data.filename,
    text: res.data.text,
  }]
  ElMessage.success('Upload successful, text extracted')
  e.target.value = ''
}

const onSave = async () => {
  const data = {
    title: form.value.title,
    content: form.value.content,
    tags: tagsText.value.split(',').map((s) => s.trim()).filter(Boolean),
    attachments: form.value.attachments || [],
  }
  if (id) {
    await updateArticle(id, data)
    ElMessage.success('Update successful')
  } else {
    await createArticle(data)
    ElMessage.success('Published successfully')
  }
  router.push('/')
}
</script>

<style scoped>
.edit {
  background: rgba(255, 250, 243, 0.86);
  border: 1px solid var(--paper-edge);
  border-radius: 22px;
  padding: 30px 32px;
  box-shadow: 0 12px 32px rgba(80, 58, 32, 0.06);
}

.edit-title {
  font-family: var(--font-serif);
  font-size: 30px;
  color: var(--leaf-deep);
  margin-bottom: 22px;
  text-align: left;
}

.edit-input,
.edit-textarea {
  margin-bottom: 14px;
}

.edit-input :deep(.el-input__wrapper),
.edit-textarea :deep(.el-textarea__inner) {
  background-color: #fff !important;
  border: 1px solid var(--paper-edge) !important;
  color: var(--ink) !important;
  box-shadow: none !important;
  border-radius: 12px !important;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.edit-input :deep(.el-input__wrapper) {
  min-height: 42px;
}

.edit-input :deep(.el-input__wrapper:hover),
.edit-textarea :deep(.el-textarea__inner:hover) {
  border-color: rgba(63, 107, 79, 0.35) !important;
}

.edit-input :deep(.el-input__wrapper.is-focus),
.edit-textarea :deep(.el-textarea__inner:focus) {
  border-color: rgba(63, 107, 79, 0.55) !important;
  box-shadow: 0 0 0 3px rgba(63, 107, 79, 0.1) !important;
}

.edit-input :deep(.el-input__inner) {
  color: var(--ink) !important;
}

.edit-textarea :deep(.el-textarea__inner) {
  color: var(--ink) !important;
  font-family: var(--font-sans);
  line-height: 1.8;
  padding: 12px 14px;
}

.edit-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-top: 24px;
  padding-top: 20px;
  border-top: 1px solid var(--paper-edge);
}

.edit-actions-right {
  display: flex;
  align-items: center;
  gap: 10px;
}

/* 抵消 Element Plus 相邻按钮默认的外边距，避免间距翻倍 */
.edit-actions-right :deep(.el-button + .el-button) {
  margin-left: 0;
}

.btn-upload,
.btn-back,
.btn-save {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  height: 36px;
  padding: 0 20px;
  font-size: 13px;
  font-weight: 500;
  border-radius: 999px !important;
  transition: color 0.18s, background 0.18s, border-color 0.18s;
}

.btn-upload,
.btn-back {
  background: transparent !important;
  border-color: var(--paper-edge) !important;
  color: var(--muted) !important;
}

.btn-upload:hover,
.btn-back:hover {
  background: rgba(63, 107, 79, 0.06) !important;
  border-color: rgba(63, 107, 79, 0.35) !important;
  color: var(--leaf-deep) !important;
}

.btn-save {
  background: var(--leaf) !important;
  border-color: var(--leaf) !important;
  color: #fffaf3 !important;
  font-weight: 600;
}

.btn-save:hover {
  background: var(--leaf-deep) !important;
  border-color: var(--leaf-deep) !important;
}

@media (max-width: 560px) {
  .edit {
    padding: 22px 18px;
  }

  .edit-title {
    font-size: 26px;
  }

  .edit-actions {
    flex-direction: column;
    align-items: stretch;
  }

  .btn-upload {
    width: 100%;
  }

  .edit-actions-right {
    width: 100%;
  }

  .edit-actions-right :deep(.el-button) {
    flex: 1;
  }
}
</style>
