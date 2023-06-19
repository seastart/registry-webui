<script setup lang="ts">
import { ref } from 'vue';
import { RouterLink, RouterView } from 'vue-router';
import { Api } from './api';
const api = new Api();
const title = ref('docker registry web ui');
api.getAppInfo().then((appTitle:string) => {
    title.value = appTitle;
    document.title = appTitle;
});
</script>

<template>
    <div class="dropdown colormoder">
        <button class="btn btn-link nav-link dropdown-toggle opacity-50" id="bd-theme" type="button"
            data-bs-toggle="dropdown" aria-expanded="false">
            <i class="bi me-1 theme-icon-active"></i>
            <span class="d-lg-none" id="bd-theme-text">{{$t("message.toggleMode")}}</span>
        </button>
        <ul class="dropdown-menu">
            <li>
                <button type="button" class="dropdown-item d-flex align-items-center" data-bs-theme-value="light"
                    aria-pressed="false">
                    <i class="bi bi-sun-fill me-2 opacity-50"></i>
                    {{$t("message.lightMode")}}
                </button>
            </li>
            <li>
                <button type="button" class="dropdown-item d-flex align-items-center active" data-bs-theme-value="dark"
                    aria-pressed="true">
                    <i class="bi bi-moon-stars-fill me-2 opacity-50"></i>
                    {{$t("message.darkMode")}}
                </button>
            </li>
            <li>
                <button type="button" class="dropdown-item d-flex align-items-center" data-bs-theme-value="auto"
                    aria-pressed="false">
                    <i class="bi bi-circle-half me-2 opacity-50"></i>
                    {{$t("message.autoMode")}}
                </button>
            </li>
        </ul>
    </div>
    <header class="py-2 border-bottom">
        <RouterLink to="/" class="link-body-emphasis link-underline link-underline-opacity-0"><img alt="Seastart" class="logo" src="@/assets/logo.svg"/>{{ title }}
        </RouterLink>
    </header>

    <RouterView />
</template>
