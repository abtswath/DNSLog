<template>
    <NLayout>
        <NLayoutContent :native-scrollbar="false"
            :content-style="{ height: '100vh', maxWidth: '1200px', padding: '24px' }">
            <NSpace vertical>
                <NSpace :vertical="false">
                    <NButton :loading="logoutLoading" @click="logout">退出登录</NButton>
                    <NForm label-placement="left" ref="formRef" inline @submit.prevent="search">
                        <NFormItem label="域名" path="domain">
                            <NInput v-model:value="formData.domain" />
                        </NFormItem>
                        <NFormItem label="IP" path="ip">
                            <NInput v-model:value="formData.ip" />
                        </NFormItem>
                        <NButton attr-type="submit">搜索</NButton>
                    </NForm>
                </NSpace>
                <NDataTable :loading="loading" :columns="columns" :data="data"></NDataTable>
                <NPagination @update:page="handlePageUpdate" :page="props.page" :pageCount="pageCount" />
            </NSpace>
        </NLayoutContent>
    </NLayout>
</template>

<script lang="ts" setup>
import { NLayout, NLayoutContent, NSpace, NButton, NForm, NFormItem, NInput, NDataTable, DataTableColumns, NPagination, FormInst } from 'naive-ui';
import { reactive, Ref, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import useRequest from '../compositions/use-request';

const props = defineProps<{
    page?: number
    domain?: string
    ip?: string
}>();

const formRef = ref<FormInst>();
const pageCount = ref(0);

interface Record {
    domain: string
    ip: string
    created_at: string
}

const columns: DataTableColumns = [
    {
        key: 'domain',
        title: '域名'
    },
    {
        key: 'ip',
        title: 'IP'
    },
    {
        key: 'created_at',
        title: '时间'
    }
];
const data = ref<Record[]>([]) as Ref<Record[]>;

const { request, loading } = useRequest<{ size: number, data: Record[] }>();

const load = () => {
    request('/logs', {
        params: {
            page: props.page,
            domain: props.domain,
            ip: props.ip
        }
    })
        .then((response) => {
            pageCount.value = Math.ceil(response.size / 30);
            data.value = response.data || [];
        })
        .catch(() => { });
}
load();
const route = useRoute();
const router = useRouter();

const formData = reactive({
    domain: props.domain || '',
    ip: props.ip || ''
});
const search = () => {
    formRef.value?.validate()
        .then(() => {
            router.push({
                query: {
                    ...route.query,
                    page: 1,
                    domain: formData.domain,
                    ip: formData.ip
                }
            })
        })
        .catch(() => { });
};

const handlePageUpdate = () => {
    router.push({
        query: {
            ...route.query,
            page: props.page || 1,
        }
    })
};

watch(() => route.query, load);

const { request: logoutRequest, loading: logoutLoading } = useRequest();

const logout = () => {
    logoutRequest('/session', { method: 'DELETE' })
        .then(() => { })
        .catch(() => { })
        .finally(() => {
            router.push({
                name: 'login'
            });
        });
}

</script>

<style lang="scss" scoped>
</style>
