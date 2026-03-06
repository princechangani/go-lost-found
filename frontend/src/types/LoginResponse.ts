export interface LoginResponse {
    User: User
    httpStatus: number
    message: string
    status: boolean
    timestamp: string
}

export interface User {
    id: string
    email: string
    password: string
    role: string
    fcmToken: any
    bearerToken: string
}
