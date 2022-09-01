import axios from 'axios';

const request = axios.create({
    timeout: 30 * 1000,
    baseURL: '/api'
});

let controller: AbortController = new AbortController();
request.interceptors.request.use((config) => {
    config.signal = controller.signal;
    return config;
});

request.interceptors.response.use(response => {
    return response;
}, error => {
    error.message = error?.response?.data?.message || '系统错误，请稍后重试';
    if (error.response && error.response.status === 401) {
        controller.abort();
        controller = new AbortController();
    }
    return Promise.reject(error);
});

export default request;
