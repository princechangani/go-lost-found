import React, {useState, useEffect} from 'react';
import {apiClient} from "@/api/apiClient.ts";
import {API_ENDPOINTS} from "@/api/endpoints.ts";
import {useNavigate} from "react-router-dom";
import {Category} from "@/types/Category.ts";
import {showErrorToast, showSuccessToast} from "@/components/CustomToasts.ts";

// Define the Items interface (you might want to move this to a separate types file)
interface Items {
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
    contact?: string;
}

interface ReportItemProps {
    categories: Category[]
}

const ReportItem: React.FC<ReportItemProps> = ({categories}) => {
    const navigate = useNavigate();
    const [selectedCategory, setSelectedCategory] = useState<Category>({id: '', name: ''});
    const [file, setFile] = useState<File | null>(null);
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [itemType, setItemType] = useState<'lost' | 'found'>('lost');
    const [contact, setContact] = useState('');

    const [item, setItem] = useState<Items>({
        id: "",
        title: '',
        description: '',
        location: '',
        date: '',
        categoryId: '',
        image: '',
        statusType: '',
        type: '',
        time: '',
        contact: ''
    });

    // Fetch categories on component mount


    // Update item when itemType changes
    useEffect(() => {
        setItem(prev => ({
            ...prev,
            type: itemType,
            statusType: itemType === 'lost' ? 'pending' : 'unclaimed'
        }));
    }, [itemType]);

    // Update item when selectedCategory changes
    useEffect(() => {
        setItem(prev => ({
            ...prev,
            categoryId: selectedCategory.id
        }));
    }, [selectedCategory]);

    const validateForm = (): boolean => {
        if (!item.title.trim()) {
            alert('Please enter an item name');
            return false;
        }
        if (!item.description.trim()) {
            alert('Please enter a description');
            return false;
        }
        if (!selectedCategory.id) {
            alert('Please select a category');
            return false;
        }
        if (!item.location.trim()) {
            alert('Please enter a location');
            return false;
        }
        if (!item.date) {
            alert('Please select a date');
            return false;
        }
        if (!item.time) {
            alert('Please select a time');
            return false;
        }
        if (!contact.trim()) {
            alert('Please enter your contact information');
            return false;
        }
        if (!file) {
            alert('Please upload an image');
            return false;
        }
        return true;
    };
   const convertTo12HourTime = (time: string): string => {
        let [hours, minutes] = time.split(':').map(Number);
        let period = 'AM';
        if (hours >= 12) {
            period = 'PM';
        }
        if (hours > 12) {
            hours -= 12;
        }
        if (hours === 0) {
            hours = 12;
        }
        return `${String(hours).padStart(2, '0')}:${String(minutes).padStart(2, '0')} ${period}`;
    };



    const addItem = async () => {
        if (!validateForm()) return;

        setIsSubmitting(true);


        try {
            const formData = new FormData();
            formData.append("title", item.title);
            formData.append("description", item.description);
            formData.append("location", item.location);
            formData.append("date", item.date.toString()); // "YYYY-MM-DD"
            formData.append("categoryId", selectedCategory.id);
            formData.append("statusType", item.statusType);
            formData.append("type", item.type);
            formData.append("email", contact);   // BUG-020 FIX: backend reads "email" field
            formData.append("time", convertTo12HourTime(item.time).toString());
            if (file) {
                formData.append("file", file);   // BUG-019 FIX: matches backend c.FormFile("file")
            }

            const response = await apiClient.post<any>(API_ENDPOINTS.CREATE_ITEM, formData, {
                headers: {
                    'Content-Type': 'multipart/form-data',
                    // BUG-016 FIX: Authorization header is now set by axios interceptor with Bearer prefix
                }
            });

            console.log('Item created successfully:', response.data);
            showSuccessToast("Item reported successfully");

            // Navigate based on item type
            if (itemType === "lost") {
                navigate("/lost-items");
            } else {
                navigate("/found-items");
            }
        } catch (error) {
            console.error("Upload error:", error);
            showErrorToast("Error submitting report. Please try again.");
        } finally {
            setIsSubmitting(false);
        }
    };

    const handleCategoryChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
        const categoryId = e.target.value;
        const category = categories.find(cat => cat.id === categoryId);
        if (category) {
            setSelectedCategory(category);
        }
    };

    const handleItemTypeChange = (type: 'lost' | 'found') => {
        setItemType(type);
        setItem(prev => ({
            ...prev,
            statusType: type === 'lost' ? 'pending' : 'unclaimed'
        }));
    };

    return (
        <div className="max-w-2xl mx-auto p-4">
            <h1 className="text-2xl font-bold text-gray-900 mb-6">Report an Item</h1>

            <div className="bg-white rounded-lg shadow-md p-6">
                {/* Item Type Selection */}
                <div className="flex space-x-4 mb-6">
                    <button
                        type="button"
                        className={`flex-1 py-2 rounded-md transition-colors ${
                            itemType === 'lost'
                                ? 'bg-blue-600 text-white'
                                : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                        }`}
                        onClick={() => handleItemTypeChange('lost')}
                    >
                        Lost Item
                    </button>
                    <button
                        type="button"
                        className={`flex-1 py-2 rounded-md transition-colors ${
                            itemType === 'found'
                                ? 'bg-blue-600 text-white'
                                : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                        }`}
                        onClick={() => handleItemTypeChange('found')}
                    >
                        Found Item
                    </button>
                </div>

                <form className="space-y-6" onSubmit={(e) => e.preventDefault()}>
                    {/* Item Name */}
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Item Name *
                        </label>
                        <input
                            type="text"
                            value={item.title}
                            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                            placeholder="Enter item name"
                            onChange={(e) => setItem({...item, title: e.target.value})}
                            required
                        />
                    </div>

                    {/* Description */}
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Description *
                        </label>
                        <textarea
                            value={item.description}
                            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                            rows={4}
                            placeholder="Describe the item in detail"
                            onChange={(e) => setItem({...item, description: e.target.value})}
                            required
                        />
                    </div>

                    {/* Category */}
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Category *
                        </label>
                        <select
                            value={selectedCategory.id}
                            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                            onChange={handleCategoryChange}
                            required
                        >
                            <option value="">Select a category</option>
                            {categories.map((category) => (
                                <option key={category.id} value={category.id}>
                                    {category.name}
                                </option>
                            ))}
                        </select>
                    </div>

                    {/* Location */}
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Location *
                        </label>
                        <input
                            type="text"
                            value={item.location}
                            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                            placeholder="Where was the item lost/found?"
                            onChange={(e) => setItem({...item, location: e.target.value})}
                            required
                        />
                    </div>

                    {/* Date and Time */}
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-1">
                                Date *
                            </label>
                            <input
                                type="date"
                                value={item.date}
                                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                                onChange={(e) => setItem({...item, date: e.target.value})}
                                required
                            />
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-1">
                                Time *
                            </label>
                            <input
                                type="time"
                                value={item.time}
                                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                                onChange={(e) => setItem({...item, time: e.target.value})}
                                required
                            />
                        </div>
                    </div>

                    {/* Contact Information */}
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Contact Information *
                        </label>
                        <input
                            type="email"
                            value={contact}
                            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                            placeholder="Enter your email"
                            onChange={(e) => setContact(e.target.value)}
                            required
                        />
                    </div>

                    {/* Image Upload */}
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Upload Image *
                        </label>
                        <div
                            className="p-4 border-2 border-dashed border-gray-300 rounded-md hover:border-blue-400 transition-colors">
                            <input
                                type="file"
                                accept="image/*"
                                className="w-full"
                                onChange={(e) => {
                                    if (e.target.files && e.target.files[0]) {
                                        setFile(e.target.files[0]);
                                    }
                                }}
                                required
                            />
                            {file && (
                                <p className="mt-2 text-sm text-gray-600">
                                    Selected: {file.name}
                                </p>
                            )}
                        </div>
                    </div>

                    {/* Submit Button */}
                    <button
                        type="submit"
                        className="w-full bg-blue-600 text-white py-2 rounded-md hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
                        onClick={addItem}
                        disabled={isSubmitting}
                    >
                        {isSubmitting ? 'Submitting...' : 'Submit Report'}
                    </button>
                </form>
            </div>
        </div>
    );
};


export default ReportItem;