<template>
    <NConfigProvider :locale="zhCN">
        <NMessageProvider>
            <router-view></router-view>
        </NMessageProvider>
    </NConfigProvider>
</template>

<script lang="ts" setup>
import { NConfigProvider, NMessageProvider, zhCN } from 'naive-ui';
</script>
