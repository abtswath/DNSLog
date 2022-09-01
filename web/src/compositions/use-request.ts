import useLoading from './use-loading';
import request from '../libs/request';
import { AxiosRequestConfig, AxiosResponse } from 'axios';
import { useMessage } from 'naive-ui';
import { useRouter } from 'vue-router';

export interface Response<T> {
    code: number
    message: string
    data: T
}

export default <T>() => {
    const [loading, start, finish] = useLoading();

    const message = useMessage();
    const router = useRouter();

    const httpRequest = <D>(url: string, config?: AxiosRequestConfig<D>) => {
        start();
        return new Promise<T>((resolve, reject) => {
            request.request<Response<T>>({
                url,
                ...config
            })
                .then((response) => {
                    if (response.data.code !== 0) {
                        throw new Error(response.data.message);
                    }
                    resolve(response.data.data);
                })
                .catch((err) => {
                    if (err.response && err.response.status && err.response.status === 401) {
                        router.push({
                            name: 'login'
                        })
                            .then(() => {
                            })
                            .catch(() => {
                            });
                    } else {
                        message.error(err.message);
                    }
                    reject(err);
                })
                .finally(finish);
        });
    }

    return {
        request: httpRequest,
        loading
    };
}
