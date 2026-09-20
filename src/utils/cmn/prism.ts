import Prism from 'prismjs'
// 编辑器 codesample 插件产生的是 Prism 格式的标记(language-markup 等)，
// 官方文档要求前台使用 Prism.js 渲染：https://www.tiny.cloud/docs/tinymce/latest/codesample/
// Prism 核心已包含 markup、css、clike、javascript
import 'prismjs/components/prism-bash'
import 'prismjs/components/prism-c'
import 'prismjs/components/prism-cpp'
import 'prismjs/components/prism-csharp'
import 'prismjs/components/prism-go'
import 'prismjs/components/prism-java'
import 'prismjs/components/prism-json'
import 'prismjs/components/prism-markdown'
import 'prismjs/components/prism-markup-templating'
import 'prismjs/components/prism-php'
import 'prismjs/components/prism-python'
import 'prismjs/components/prism-ruby'
import 'prismjs/components/prism-sql'
import 'prismjs/components/prism-typescript'
import 'prismjs/components/prism-yaml'

// 编辑器代码块语言下拉框：与上方已注册的语言保持一致
export const CODESAMPLE_LANGUAGES = [
  { text: 'HTML/XML', value: 'markup' },
  { text: 'CSS', value: 'css' },
  { text: 'JavaScript', value: 'javascript' },
  { text: 'TypeScript', value: 'typescript' },
  { text: 'JSON', value: 'json' },
  { text: 'YAML', value: 'yaml' },
  { text: 'Bash/Shell', value: 'bash' },
  { text: 'Go', value: 'go' },
  { text: 'Python', value: 'python' },
  { text: 'Java', value: 'java' },
  { text: 'C', value: 'c' },
  { text: 'C#', value: 'csharp' },
  { text: 'C++', value: 'cpp' },
  { text: 'PHP', value: 'php' },
  { text: 'Ruby', value: 'ruby' },
  { text: 'SQL', value: 'sql' },
  { text: 'Markdown', value: 'markdown' },
]

export default Prism
