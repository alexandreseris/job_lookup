import { createApp } from 'vue'
import App from './App.vue'
import './style.css';
import { createVuetify } from 'vuetify'
import { createPinia } from 'pinia'


const vuetify = createVuetify()
const pinia = createPinia()

createApp(App).use(vuetify).use(pinia).mount('#app')
