import React, { useEffect } from "react";
import {apiClient} from "@/api/apiClient.ts";
import {API_ENDPOINTS} from "@/api/endpoints.ts";
import {Contact, ContactResponse} from "@/types/Contact.ts";
import {ExceptionResponse} from "@/types/ErrorResponse.ts";
import {showErrorToast, showSuccessToast} from "@/components/CustomToasts.ts";
import {Items} from "@/types/Item.ts";

interface ContactDialogProps {
    open: boolean;
    item: Items|null;
    onClose: () => void;
}

const ContactDialog: React.FC<ContactDialogProps> = ({  open, item, onClose }) => {
    const [contactData, setContactData] = React.useState<Contact>({
        email: "",
        name: "",
        message: "",
        itemId: item?.id ?? "",
    });

    // BUG-022 FIX: sync itemId whenever the item prop changes (e.g. user clicks a different card)
    useEffect(() => {
        setContactData(prev => ({ ...prev, itemId: item?.id ?? "" }));
    }, [item]);



    const validation = () => {
        if (contactData.email === "" || contactData.name === "" || contactData.message === "") {
            return false;
        }
        return true;
    };

    const createContact = async () => {
        console.log(contactData);
        if(!validation()) {
            showErrorToast("Please fill all the fields");
            return;
        }


        try {
            const response = await apiClient.post<ContactResponse|ExceptionResponse>(API_ENDPOINTS.CREATE_CONTACT, contactData,{
                headers: {
                    'Content-Type': 'application/json'
                }
            });
            // BUG-013 FIX: backend now returns 201 for successful contact creation
            if (response.status === 201) {
                showSuccessToast("Contact created successfully");
                if(item?.type==="lost") {
                    item.statusType = "found"
                } else {
                    if(item != null) {
                        item.statusType = "claimed";
                    }
                }
                setContactData({
                    email: "",
                    name: "",
                    message: "",
                    itemId: ""
                });
                onClose();
            } else {
                showErrorToast(response.data.message);
            }
        } catch (error: any) {
            showErrorToast(error?.response?.data?.message ?? "Something went wrong. Please try again.");
            console.error(error);
        }
    };

    if (!open) return null;

    return (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white rounded-lg p-6 w-full max-w-md mx-4">
                <h2 className="text-lg font-semibold mb-4">Contact Owner</h2>

                <div className="space-y-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Email *
                        </label>
                        <input
                            type="email"
                            value={contactData.email}
                            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                            placeholder="Enter your Email"
                            onChange={(e) => setContactData({...contactData, email: e.target.value})}
                            required
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Name *
                        </label>
                        <input
                            type="text"
                            value={contactData.name}
                            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                            placeholder="Enter your Name"
                            onChange={(e) => setContactData({...contactData, name: e.target.value})}
                            required
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Message *
                        </label>
                        <textarea
                            value={contactData.message}
                            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                            placeholder="Enter your Message"
                            onChange={(e) => setContactData({...contactData, message: e.target.value})}
                            rows={4}
                            required
                        />
                    </div>
                </div>

                <div className="flex justify-end gap-2 mt-6">
                    <button
                        onClick={onClose}
                        className="bg-gray-300 px-4 py-2 rounded-md hover:bg-gray-400"
                    >
                        Cancel
                    </button>
                    <button
                        onClick={createContact}
                        className="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700"
                    >
                        Submit
                    </button>
                </div>
            </div>
        </div>
    );
};

export default ContactDialog;