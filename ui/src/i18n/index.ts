import { createI18n } from 'vue-i18n'
const messages = {
    zh: {
        message: {
            lightMode: "日间模式",
            darkMode: "夜间模式",
            autoMode: "自动",
            toggleMode: "显示模式",
            repo: "仓库",
            desc: "简介",
            lastUpdated: "上次更新",
            size: "大小",
            loading: "加载中...",
            loadMore: "加载更多...",
            noMore: "没有更多了~",
        }
    },
    en: {
        message: {
            lightMode: "Light",
            darkMode: "Dark",
            autoMode: "Auto",
            toggleMode: "Toggle Theme",
            repo: "Repo",
            desc: "Desc",
            lastUpdated: "Last Update",
            size: "Size",
            loading: "loading...",
            loadMore: "load more...",
            noMore: "no more~",
        }
    }
}
const language = (navigator.language || 'en').toLocaleLowerCase() 
const i18n = createI18n({
    locale: language.split('-')[0] || 'en',
    fallbackLocale: 'en',
    messages: messages,
})
export default i18n