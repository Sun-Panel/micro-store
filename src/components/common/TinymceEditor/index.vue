<script setup lang="ts">
import { computed, ref, withDefaults } from 'vue'
import Editor from '@tinymce/tinymce-vue'
import type { RawEditorOptions } from 'tinymce'
import { useOsTheme } from 'naive-ui'
import { useAppStore, useAuthStore } from '@/store'
import { apiRespErrMsg } from '@/utils/cmn'
import Prism, { CODESAMPLE_LANGUAGES } from '@/utils/cmn/prism'

// 编辑器内置的 Prism 语言很少(无 Go/TypeScript 等)，
// 配合 codesample_global_prismjs 让它复用上面注册好的这一份
Object.assign(window, { Prism })

const props = withDefaults(defineProps<{
  content?: string
  initConfig?: RawEditorOptions
}>(), {
})

const emit = defineEmits<{
  (e: 'update:content', v: string): void
}>()

const authStore = useAuthStore()
const appStore = useAppStore()
const osTheme = useOsTheme()

// 与站点主题保持一致
const isDark = computed(() => {
  if (appStore.theme === 'auto')
    return osTheme.value === 'dark'
  return appStore.theme === 'dark'
})

// 图片上传自定义逻辑
/** 文档：https://www.tiny.cloud/docs/tinymce/latest/upload-images/ */
type UploadFn = RawEditorOptions['images_upload_handler']
const handleFileUpload: UploadFn = async (blobInfo: any): Promise<string> => {
  const formData = new FormData()
  formData.append('imgfile', blobInfo.blob())

  // 添加自定义的头部参数
  const customHeaders = {
    token: authStore.token || '',
  }

  try {
    const response = await fetch('/api/file/uploadImg', {
      method: 'POST',
      body: formData,
      headers: new Headers(customHeaders),
    })
    const res = await response.json()
    if (response.ok) {
      let url = ''
      if (res.code === 0) {
        url = res.data.imageUrl
      }
      else {
        apiRespErrMsg(res)
        throw new Error(res.msg || 'Upload failed')
      }
      // 假设上传成功后返回文件的URL
      return url
    }
    else {
      throw new Error(res.message || 'Upload failed')
    }
  }
  catch (error) {
    console.error('Upload error:', error)
    throw error
  }
}

const contentRef = computed({
  get: () => props.content || '',
  set(value: string | undefined) {
    emit('update:content', value || '')
  },
})

// 初始化配置
const initConfig = ref<RawEditorOptions>({
  ...props.initConfig,
  // 如果不设置license_key为gpl，console会一直提示一个模式问题
  license_key: 'gpl',
  // 移除tinymce右上角升级提示
  promotion: false,
  // 移除tinymce右下角品牌提示
  branding: false,
  // 默认高度，调用方可通过 initConfig 覆盖
  height: props.initConfig?.height ?? 500,
  // 跟随站点主题
  skin: props.initConfig?.skin ?? (isDark.value ? 'oxide-dark' : 'oxide'),
  content_css: props.initConfig?.content_css ?? (isDark.value ? 'dark' : 'default'),
  // 设置语言
  language: props.initConfig?.language || appStore.language.replace('-', '_'),
  // 配置插件列表
  plugins: props.initConfig?.plugins || 'preview importcss searchreplace autolink save directionality code visualblocks visualchars fullscreen image link media codesample table charmap pagebreak nonbreaking anchor insertdatetime advlist lists wordcount help charmap quickbars emoticons accordion',
  // 配置工具栏
  toolbar: props.initConfig?.toolbar || [
    'undo redo | blocks fontfamily fontsize | link image | table media | code fullscreen preview',
    'bold italic underline strikethrough | lineheight align numlist bullist | forecolor backcolor removeformat | charmap emoticons codesample',
  ],
  valid_children: '+div[style]',
  // 使用页面内的 Prism，从而支持内置版本不具备的语言
  codesample_global_prismjs: true,
  codesample_languages: CODESAMPLE_LANGUAGES,
  // 保持服务端返回的 URL 原样：默认 convert_urls + relative_urls 会把 /uploads/xxx
  // 改写成相对后台编辑页的路径，导致前台 /hPage/xxx 预览时图片无法解析
  convert_urls: false,
  relative_urls: false,
  remove_script_host: false,
  // 内容默认样式
  content_style:
    'body { font-family:Helvetica,Arial,sans-serif; font-size:14px }',
  // 配置图片上传功能
  images_upload_handler: handleFileUpload,
})
</script>

<template>
  <Editor
    v-model="contentRef"
    :init="initConfig"
    tinymce-script-src="/tinymce/tinymce.min.js"
  />
</template>
