import React, {useEffect, useState} from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import Navbar from './components/GlassmorphismNavbar.tsx';
import Home from './pages/Home';
import LostItems from './pages/LostItems';
import FoundItems from './pages/FoundItems';
import ReportItem from './pages/ReportItem';
import Footer from './components/Footer';
import { SearchProvider } from './context/SearchContext';
import {Category, CategoryResponse} from "@/types/Category.ts";
import {apiClient} from "@/api/apiClient.ts";
import {API_ENDPOINTS} from "@/api/endpoints.ts";
import {RefreshProvider} from "@/context/RefreshContext.tsx";
import AuthPage from './pages/AuthPage.tsx';

const App: React.FC = () => {

    const [categories1, setCategories] = useState<Category[]>([]);
    // BUG-018 FIX: derive auth state directly from localStorage — no stale state
    const [isAuthenticated, setIsAuthenticated] = useState<boolean>(
        () => !!localStorage.getItem('token')
    );

    useEffect(() => {
        if (!isAuthenticated) return;
        const fetchCategories = async () => {
            try {
                const response = await apiClient.get<CategoryResponse>(API_ENDPOINTS.GET_ALL_CATEGORIES);
                setCategories(response.data.data);
            } catch (error) {
                console.error('Error fetching categories:', error);
            }
        };
        fetchCategories();
    }, [isAuthenticated]);

    const handleLoginSuccess = () => {
        setIsAuthenticated(true);
    };

    // BUG-017 FIX: if not authenticated, always render AuthPage regardless of path
    if (!isAuthenticated) {
        return <AuthPage onLoginSuccess={handleLoginSuccess} />;
    }

    return (
        <React.StrictMode>
            <RefreshProvider>
                <SearchProvider>
                    <div className="min-h-screen bg-gray-50">
                        <Navbar />
                        <div className="container mx-auto px-4 py-8">
                            <Routes>
                                {/* BUG-017 FIX: redirect root to /home */}
                                <Route path="/" element={<Navigate to="/home" replace />} />
                                <Route path="/home" element={<Home />} />
                                <Route path="/lost-items" element={<LostItems categories={categories1} />} />
                                <Route path="/found-items" element={<FoundItems categories={categories1} />} />
                                <Route path="/report-item" element={<ReportItem categories={categories1} />} />
                                {/* Catch-all redirect */}
                                <Route path="*" element={<Navigate to="/home" replace />} />
                            </Routes>
                        </div>
                        <Footer />
                    </div>
                </SearchProvider>
            </RefreshProvider>
        </React.StrictMode>
    );
};

export default App;