import { createRouter, createWebHashHistory } from "vue-router";
import PocManage from '../components/PocManage.vue'
import PocScan from "../components/PocScan.vue";
import BatchScan from "../components/BatchScan.vue";
import Exploitation from "../components/Exploitation.vue";
import FingerprintManage from  "../components/FingerprintManage.vue"
import Other from  "../components/Other.vue"
import FingerprintScan from "../components/FingerprintScan.vue"

const router = createRouter({
    history: createWebHashHistory(),
    routes: [
        {
            path: '/',
            name: 'pocManage',
            component: PocManage

        },
        {
            path: '/setting',
            name: 'Setting',
            component: () => import('../components/Setting.vue'),

        },
        {
            path: '/scan',
            name: 'pocScan',
            component: PocScan,
        },
        {
            path: '/batchScan',
            name: 'batchScan',
            component: BatchScan,
        },
        {
            path: '/exploitation',
            name: 'exploitation',
            component: Exploitation,
        },
        {
            path: '/fingerprintManage',
            name: 'fingerprintManage',
            component: FingerprintManage,
        },
        {
            path: '/other',
            name: 'other',
            component: Other,
        },
        {
            path: '/fingerprintScan',
            name: 'fingerprintScan',
            component: FingerprintScan,
        },
    ]
})

export default router
