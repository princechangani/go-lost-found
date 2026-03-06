import React, { createContext, useContext, useState, ReactNode } from 'react';

// Define the context type
interface RefreshContextType {
    refresh: boolean;
    triggerRefresh: () => void;
}

// Create context with default value
const RefreshContext = createContext<RefreshContextType | undefined>(undefined);

// Custom hook for consuming context
export const useRefresh = (): RefreshContextType => {
    const context = useContext(RefreshContext);
    if (!context) {
        throw new Error('useRefresh must be used within a RefreshProvider');
    }
    return context;
};

// Provider component
interface Props {
    children: ReactNode;
}

export const RefreshProvider: React.FC<Props> = ({ children }) => {
    const [refresh, setRefresh] = useState(false);

    const triggerRefresh = () => setRefresh(prev => !prev); // toggle to force re-run

    return (
        <RefreshContext.Provider value={{ refresh, triggerRefresh }}>
            {children}
        </RefreshContext.Provider>
    );
};
