import React, { useState, useEffect } from "react";
import { apiClient } from "@/api/apiClient.ts";
import { LoginResponse } from "@/types/LoginResponse.ts";
import { API_ENDPOINTS } from "@/api/endpoints.ts";
import { useNavigate } from "react-router-dom";

interface AuthPageProps {
    onLoginSuccess?: () => void;
}

const AuthPage: React.FC<AuthPageProps> = ({ onLoginSuccess }) => {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [error, setError] = useState<string | null>(null);
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate();

    useEffect(() => {
        const token = localStorage.getItem("token");
        if (token) {
            onLoginSuccess?.();
            navigate("/home");
        }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    const handleLogin = async (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);
        setLoading(true);

        try {
            const response = await apiClient.post<LoginResponse>(API_ENDPOINTS.LOGIN, {
                email,
                password,
            });

            localStorage.setItem("token", response.data.User.bearerToken.trim());
            console.log(response.data.User.bearerToken);
            onLoginSuccess?.();
            navigate("/home");
        } catch (error: any) {
            setError("Invalid email or password.");
            console.error("Login failed:", error);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-100 via-indigo-200 to-purple-200">
            <div className="bg-white shadow-2xl rounded-2xl p-8 w-full max-w-md">
                <h1 className="text-3xl font-bold text-center text-gray-800 mb-6">
                    Welcome Back 👋
                </h1>
                <p className="text-center text-gray-500 mb-8">
                    Please sign in to continue
                </p>

                <form onSubmit={handleLogin} className="space-y-6">
                    <div>
                        <label className="block text-gray-700 font-medium mb-1">Email</label>
                        <input
                            type="email"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            required
                            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
                            placeholder="Enter your email"
                        />
                    </div>

                    <div>
                        <label className="block text-gray-700 font-medium mb-1">Password</label>
                        <input
                            type="password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            required
                            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
                            placeholder="Enter your password"
                        />
                    </div>

                    {error && <p className="text-red-500 text-sm text-center">{error}</p>}

                    <button
                        type="submit"
                        disabled={loading}
                        className={`w-full py-3 mt-2 text-white font-semibold rounded-lg transition ${
                            loading
                                ? "bg-indigo-400 cursor-not-allowed"
                                : "bg-indigo-600 hover:bg-indigo-700"
                        }`}
                    >
                        {loading ? "Logging in..." : "Login"}
                    </button>
                </form>

                <p className="text-center text-gray-500 text-sm mt-6">
                    Don’t have an account?{" "}
                    <a href="#" className="text-indigo-600 hover:underline">
                        Sign up
                    </a>
                </p>
            </div>
        </div>
    );
};

export default AuthPage;
