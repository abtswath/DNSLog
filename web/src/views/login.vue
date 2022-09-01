<template>
    <NLayout position="absolute">
        <NLayoutContent position="absolute">
            <NCard class="login-card">
                <template #header>
                    登录
                </template>
                <NForm ref="formRef" :label-width="0" :model="formData" :rules="rules" label-placement="left"
                    label-align="left" @submit.prevent="login">
                    <NFormItem path="username">
                        <NInput placeholder="用户名" v-model:value="formData.username" />
                    </NFormItem>
                    <NFormItem path="password">
                        <NInput placeholder="密码" type="password" v-model:value="formData.password"></NInput>
                    </NFormItem>
                    <NFormItem>
                        <NButton type="primary" :loading="loading" block attr-type="submit">登录</NButton>
                    </NFormItem>
                </NForm>
            </NCard>
        </NLayoutContent>
    </NLayout>
</template>

<script lang="ts" setup>
import { NLayout, NLayoutContent, NCard, NForm, NFormItem, NInput, NButton, FormInst, FormRules } from 'naive-ui';
import { reactive, ref } from 'vue';
import useRequest from '../compositions/use-request';
import router from '../router';

const formRef = ref<FormInst>();

const rules: FormRules = {
    username: [{ required: true, message: '请输入用户名' }],
    password: [{ required: true, message: '请输入密码' }]
};

const formData = reactive({
    username: '',
    password: ''
});

const { request, loading } = useRequest();

const login = () => {
    formRef.value?.validate()
        .then(() => {
            request('/session', { method: 'POST', data: formData })
                .then(() => {
                    router.push({
                        name: 'index'
                    });
                })
                .catch(() => { })
        })
        .catch(() => { });
}

</script>

<style lang="scss" scoped>
.login-card {
    position: absolute;
    left: 0;
    right: 0;
    margin: auto;
    top: 50%;
    transform: translateY(-50%);
    max-width: 350px;
}
</style>
