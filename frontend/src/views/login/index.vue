<script setup>
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { NForm, NFormItem, NInput, NButton, NCheckbox, NCard, NSpace } from 'naive-ui'
import { useUserStore } from '@/stores/user'
import { useMessage } from '@/composables/useMessage'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const message = useMessage()

const form = reactive({
  username: '',
  password: '',
  remember: false
})

const loading = ref(false)
const errors = ref({})

async function handleLogin() {
  errors.value = {}

  if (!form.username) {
    errors.value.username = '请输入用户名'
    return
  }
  if (!form.password) {
    errors.value.password = '请输入密码'
    return
  }

  loading.value = true
  try {
    await userStore.login({ username: form.username, password: form.password })
    message.success('登录成功')
    const redirect = route.query.redirect || '/dashboard'
    router.push(redirect)
  } catch (error) {
    message.error(error.message || '登录失败')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-100">
    <n-card class="w-96">
      <div class="text-center mb-6">
        <h1 class="text-2xl font-bold text-primary-600">Push Notify</h1>
        <p class="text-gray-500 mt-2">代码推送通知系统</p>
      </div>
      <n-form>
        <n-form-item label="用户名" :validation-status="errors.username ? 'error' : ''" :feedback="errors.username">
          <n-input v-model:value="form.username" placeholder="请输入用户名" @keyup.enter="handleLogin" />
        </n-form-item>
        <n-form-item label="密码" :validation-status="errors.password ? 'error' : ''" :feedback="errors.password">
          <n-input v-model:value="form.password" type="password" placeholder="请输入密码" show-password-on="click" @keyup.enter="handleLogin" />
        </n-form-item>
        <n-form-item>
          <n-checkbox v-model:checked="form.remember">记住我</n-checkbox>
        </n-form-item>
        <n-form-item>
          <n-button type="primary" block :loading="loading" @click="handleLogin">登录</n-button>
        </n-form-item>
      </n-form>
    </n-card>
  </div>
</template>
