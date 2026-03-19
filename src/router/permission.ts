import type { Router } from 'vue-router'
import { createDiscreteApi } from 'naive-ui'
import { useAuthStore } from '@/store/'
import { ROLE_USER } from '@/utils/role'

const naiveApi = createDiscreteApi(['loadingBar'])

export function setupPageGuard(router: Router) {
  router.beforeEach(async (to, from, next) => {
    // const authStore = useAuthStoreWithout()
    const authStore = useAuthStore()
    naiveApi.loadingBar.start()

    // 仅普通用户角色时拦截（不能访问admin路由）
    if (authStore.userInfo?.role === ROLE_USER && to.path.includes('admin')) {
      next({ name: '404' })
      return
    }

    next()
  })

  router.afterEach((to, from) => {
    naiveApi.loadingBar.finish()
  })
}
