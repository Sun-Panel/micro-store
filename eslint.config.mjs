import antfu from '@antfu/eslint-config'

export default antfu({
  rules: {
    // 将未使用变量规则调整为警告级别
    'unused-imports/no-unused-vars': 'warn',
    // 如果使用了 typescript 的未使用变量规则也需要调整
    // '@typescript-eslint/no-unused-vars': 'warn',
    'no-console': 'warn', // 取消打印标红提醒
    // 'antfu/curly': ['error', 'multi-line'], // 这样配置确保在所有的if和else语句中使用大括号-宽松模式
  },
  ignores: ['docker-compose', '**/docker-compose/**', 'kubernetes', '**/kubernetes/**', 'service', 'service/**'],
})
