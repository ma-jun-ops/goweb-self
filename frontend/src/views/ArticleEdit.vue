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
      <el-button class="btn-upload" @click="$refs.fileInput.click()">Upload Word</el-button>
      <el-button class="btn-save" type="primary" @click="onSave">Save</el-button>
      <el-button class="btn-back" @click="$router.back()">Back</el-button>
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
  background-color: #141c33;
  border: 1px solid #1e293b;
  border-radius: 4px;
  padding: 24px;
}

.edit-title {
  font-size: 22px;
  color: #3a82ff;
  margin-bottom: 20px;
}

.edit-input,
.edit-textarea {
  margin-bottom: 12px;
}

.edit-input :deep(.el-input__wrapper),
.edit-textarea :deep(.el-textarea__inner) {
  background-color: #0c1120 !important;
  border: 1px solid #24304a !important;
  color: #cbd5e1 !important;
  box-shadow: none !important;
}

.edit-input :deep(.el-input__inner) {
  color: #cbd5e1 !important;
}

.edit-textarea :deep(.el-textarea__inner) {
  color: #cbd5e1 !important;
  font-family: 'Georgia', serif;
  line-height: 1.6;
}

.edit-actions {
  display: flex;
  gap: 12px;
  margin-top: 16px;
}

.btn-upload {
  background-color: #1e293b !important;
  border-color: #2a3654 !important;
  color: #cbd5e1 !important;
}

.btn-upload:hover {
  background-color: #24304a !important;
  border-color: #334155 !important;
}

.btn-save {
  background-color: #3a82ff !important;
  border-color: #3a82ff !important;
  color: #0c1120 !important;
  font-weight: bold;
}

.btn-save:hover {
  background-color: #6aa1ff !important;
  border-color: #6aa1ff !important;
}

.btn-back {
  background-color: #1e293b !important;
  border-color: #2a3654 !important;
  color: #cbd5e1 !important;
}

.btn-back:hover {
  background-color: #24304a !important;
  border-color: #334155 !important;
}
</style>
