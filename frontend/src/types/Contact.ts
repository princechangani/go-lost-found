export interface Contact {
    itemId: string;
    name: string;
    email: string;
    message: string;
}

export interface ContactResponse {
    status: boolean;
    statusCode: string;
    message: string;
    data: Contact;
}