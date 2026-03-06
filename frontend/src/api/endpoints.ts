// src/api/endpoints.ts


export const API_ENDPOINTS = {

    GET_ALL_ITEMS: "/items/all",
    GET_ITEM_BY_ID: (id: string) => `/items/${id}`,
    CREATE_ITEM: "/items",
    GET_ALL_CATEGORIES: "/category/all",
    CREATE_CONTACT: "/contact/create",
    GET_ALL_NOTIFICATIONS: "/notifications/all",
    LOGIN: "/user/login",
    REGISTER: "/user/register",
    GET_USER_BY_ID: (id: string) => `/users/${id}`,
    UPDATE_USER: (id: string) => `/users/${id}`,
    DELETE_USER: (id: string) => `/users/${id}`,
};
