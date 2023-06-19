<script setup lang="ts">
import { Api } from '@/api';
import "@/assets/repo.scss";
import Alert from '@/components/Alert.vue';
import Layers from '@/components/Layers.vue';
import type { Image, Repo } from '@/types';
import { formatDate, formatSize } from '@/utils';
import { Transition, nextTick, onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';

const props = defineProps<{
    reponame: string[],
}>();

const isLoading = ref(false);
const loadingError = ref(null);
const repo = ref<Repo | null>(null);
const route = useRoute();
const host = location.host;

// 声明一个 ref 来存放该元素的引用
// 必须和模板里的 ref 同名
const trLayers = ref<HTMLElement | null>(null)
const selectImage = ref<Image | null>(null)

function clickImageRow(image: Image, index: number, event: Event) {
    if (selectImage.value == image) {
        selectImage.value = null;
    } else {
        selectImage.value = image;
        nextTick(() => {
            (event.currentTarget as HTMLElement).insertAdjacentElement('afterend', trLayers.value!);
        });
    }
}

// 获取repo详情
function getRepoDetail() {
    const api = new Api();
    api.getRepoDetail(props.reponame.join('/'), true).then((crepo: Repo) => {
        isLoading.value = false;
        repo.value = crepo;
    }).catch((e) => {
        isLoading.value = false;
        loadingError.value = e;
        console.log(e);
    })
}


onMounted(() => {
    getRepoDetail();
});
</script>

<template>
    <Alert :msg="loadingError" type="danger" v-if="loadingError" />
    <div v-if="isLoading" class="loading">Loading...</div>
    <template v-if="repo">
        <div class="repo-head">
            <h4 class="name mt-2">{{ repo.name }}</h4>
            <div class="desc fw-light">{{ repo.desc }}</div>
        </div>
        <div class="tag my-2 p-2" v-for="tag in repo.tags" :key="tag.name">
            <div class="name fw-bold">{{ tag.name }}</div>
            <div class="row my-2 gx-0">
                <div class="col-lg">{{ formatDate(tag.created) }}<br> {{ tag.change_log }}</div>
                <div class="col-lg px-lg-2 border border-secondary-subtle rounded user-select-all">docker pull {{ host }}/{{ repo.name }}:{{ tag.name }}</div>
            </div>
            <div class="table-responsive">
                <table class="table table-hover">
                    <thead>
                        <tr>
                            <th>DIGEST</th>
                            <th>OS/ARCH</th>
                            <th>{{$t("message.size")}}</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="(image, index) in tag.images" :key="image.digest"
                            @click="clickImageRow(image, index, $event)" class="image-tr">
                            <td :title="image.digest">{{ image.digest.replace("sha256:", "").substring(0, 12) }}</td>
                            <td>{{ image.os }}/{{ image.arch }}</td>
                            <td>{{ formatSize(image.size) }}</td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
        <Transition name="collapse">
            <tr ref="trLayers" v-if="selectImage">
                <td colspan="3" class="bg-secondary-subtle">
                    <Layers :layers="selectImage?.layers"></Layers>
                </td>
            </tr>
        </Transition>
    </template></template>
