export interface Items {
    id: string;
    title: string;
    description: string;
    location: string;
    date: string;
    categoryId: string;
    image: string;
    statusType: string;
    type: string;
    time: string;
}

export interface ItemsResponse {
    status: boolean;
    statusCode: string;
    message: string;
    data: Items[];
}