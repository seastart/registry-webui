<script setup lang="ts">
import { Api } from '@/api';
import Alert from '@/components/Alert.vue';
import type { Repo } from '@/types';
import { formatDate } from "@/utils";
import { onMounted, ref } from 'vue';

const api = new Api();
// 声明一个 ref 来存放该元素的引用
// 必须和模板里的 ref 同名
const sentinel = ref<HTMLElement|null>(null)

const allrepoes = ref<Repo[]>([]);
const page = ref(1);
const isLoading = ref(false);
const loadingError = ref(null);
const hasMore = ref(true);

function fetchPageRepoes() {
    isLoading.value = true;
    api.getRepoes("", page.value).then(([repoes, more]) => {
        allrepoes.value.push(...repoes);
        isLoading.value = false;
        hasMore.value = more;
        page.value++;
    }).catch((e) => {
        isLoading.value = false;
        loadingError.value = e;
        console.log(e);
    })
}

onMounted(() => {
    // console.log("homeview onMounted");
    fetchPageRepoes();
    const observer = new IntersectionObserver(([entry]) => {
        if (entry.isIntersecting && !isLoading.value && hasMore.value) {
            fetchPageRepoes();
        }
    }, {
        root: null,
        rootMargin: "0px",
        threshold: 1
    });
    observer.observe(sentinel.value as HTMLElement);
})
</script>

<template>
    <Alert :msg="loadingError" type="danger" v-if="loadingError" />
    <div v-if="isLoading" class="loading">Loading...</div>
    <div class="table-responsive" v-if="!isLoading && !loadingError">
        <table class="table table-hover">
            <thead>
                <tr>
                    <th>{{$t("message.repo")}}</th>
                    <th>Tags</th>
                    <th>{{$t("message.desc")}}</th>
                    <th>{{$t("message.lastUpdated")}}</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="repo in allrepoes" :key="repo.name">
                    <td><RouterLink :to="{name:'repo', params: {reponame: repo.name.split('/')}}">{{ repo.name }}</RouterLink></td>
                    <td>
                        <template v-for="(tag,index) in repo.tags" :key="tag.name">
                            <span v-if="index<5" class="badge text-bg-primary me-1">{{ tag.name }}</span>
                        </template>
                        <span v-if="repo.tags.length>5">...</span>
                    </td>
                    <td>{{ repo.desc }}</td>
                    <td>{{ formatDate(repo.last_update) }}</td>
                </tr>
            </tbody>
        </table>
        <div v-if="isLoading" class="text-center">{{$t("message.loading")}}</div>
        <div v-else-if="hasMore" ref="sentinel" class="text-center">{{$t("message.loadMore")}}</div>
        <div v-else class="text-center">{{$t("message.noMore")}}</div>
    </div>
</template>
