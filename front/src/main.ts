import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import App from './App.vue'
import router from './router'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(ElementPlus)

for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.mount('#app')

// Global green theme override for Element Plus
const style = document.createElement('style')
style.textContent = `
  :root {
    --el-color-primary: #26A65B;
    --el-color-primary-light-1: #3cb46c;
    --el-color-primary-light-2: #51bf7c;
    --el-color-primary-light-3: #6bca8d;
    --el-color-primary-light-4: #80d49d;
    --el-color-primary-light-5: #93d9ad;
    --el-color-primary-light-6: #a8e3be;
    --el-color-primary-light-7: #beedce;
    --el-color-primary-light-8: #d3f2de;
    --el-color-primary-light-9: #e9f9ef;
    --el-color-primary-dark-2: #1e8549;
  }
`
document.head.appendChild(style)
