import { createPinia } from 'pinia';
import { createApp } from 'vue';

import App from './App.vue';
import i18n from './i18n';
import router from './router';

import 'bootstrap';
import "bootstrap-icons/font/bootstrap-icons.scss";
import "bootstrap/scss/bootstrap.scss";
import './assets/main.js';
import './assets/main.scss';

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(i18n)

app.mount('#app')
