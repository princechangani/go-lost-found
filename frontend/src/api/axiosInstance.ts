// src/api/axiosInstance.ts
import axios from 'axios';

const axiosInstance = axios.create({
  //  baseURL:'https://backend-v1-gcer.onrender.com/api/v1',
    baseURL: import.meta.env.VITE_API_URL,
    timeout: 30000, // 30s — needed for multipart file uploads
    headers: {
        'Content-Type': 'application/json',
    },
});

// Optional: Interceptors
axiosInstance.interceptors.request.use((config) => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

axiosInstance.interceptors.response.use(
    (response) => response,
    (error) => {
        console.error('API Error:', error);
        return Promise.reject(error);
    }
);

export default axiosInstance;
