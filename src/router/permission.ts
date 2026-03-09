import type { Router } from 'vue-router'
import { createDiscreteApi } from 'naive-ui'
import { useAuthStore } from '@/store/'

const naiveApi = createDiscreteApi(['loadingBar'])

export function setupPageGuard(router: Router) {
  router.beforeEach(async (to, from, next) => {
    // const authStore = useAuthStoreWithout()
    const authStore = useAuthStore()
    naiveApi.loadingBar.start()

    // 非管理员路由拦截
    if (authStore.userInfo?.role !== 1 && to.path.includes('admin'))
      next({ name: '404' })

    else
      next()
  })

  router.afterEach((to, from) => {
    naiveApi.loadingBar.finish()
  })
}
