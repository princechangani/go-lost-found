export interface Category {
        id: string;
        name: string;
    }

export interface CategoryResponse {
        status: boolean;
        statusCode: string;
        message: string;
        data: Category[];
    }

