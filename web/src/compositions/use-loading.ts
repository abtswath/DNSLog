import { Ref, ref } from 'vue'

export default (): [Ref<boolean>, () => void, () => void] => {
    const loading = ref(false);
    const start = () => {
        loading.value = true;
    }
    const finish = () => {
        loading.value = false;
    }

    return [loading, start, finish];
}
