import { RouteRecordRaw } from 'vue-router';

export default [
    {
        path: '/login',
        name: 'login',
        component: () => import('../views/login.vue')
    },
    {
        path: '/',
        name: 'index',
        component: () => import('../views/index.vue'),
        props: route => {
            const { page, domain, ip } = route.query;
            let queryPage = Number(page);
            if (isNaN(queryPage) || queryPage < 1) {
                queryPage = 1;
            }

            return {
                page: queryPage,
                domain,
                ip
            }
        }
    }
] as RouteRecordRaw[];
