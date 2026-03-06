// src/api/apiClient.ts
import axiosInstance from './axiosInstance';

export const apiClient = {
    get: <T>(url: string, params?: any) =>
        axiosInstance.get<T>(url, { params }),

    post: <T>(url: string, data?: any  , config?: any) =>
        axiosInstance.post<T>(url, data , config),

    put: <T>(url: string, data?: any) =>
        axiosInstance.put<T>(url, data),

    delete: <T>(url: string) =>
        axiosInstance.delete<T>(url),
};
