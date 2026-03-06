import React from 'react';
import {Link, useNavigate} from 'react-router-dom';
import { Search } from 'lucide-react';
import { useSearch } from '../context/SearchContext';

const Navbar: React.FC = () => {
  const { searchQuery, setSearchQuery } = useSearch();

  const handleSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchQuery(e.target.value);
  };


  const navigate = useNavigate();
 const logOut = () => {
    localStorage.removeItem("token");
    navigate("/");
    window.location.reload();

  };

  return (
      <nav className="bg-white shadow-md">
        <div className="container mx-auto px-4">
          <div className="flex items-center justify-between h-16">
            <Link to="/" className="text-2xl font-bold text-blue-600">
              Lost & Found
            </Link>

            <div className="relative flex-1 max-w-xl mx-8">
              <div className="relative">
                <input
                    type="text"
                    value={searchQuery}
                    onChange={handleSearch}
                    placeholder="Search for items..."
                    className="w-full px-4 py-2 pl-10 pr-4 rounded-lg border border-gray-300 focus:outline-none focus:border-blue-500"
                />
                <Search className="absolute left-3 top-2.5 h-5 w-5 text-gray-400" />
              </div>
            </div>

            <div className="flex space-x-4">
              <Link
                  to="/lost-items"
                  className="text-gray-600 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium"
              >
                Lost Items
              </Link>
              <Link
                  to="/found-items"
                  className="text-gray-600 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium"
              >
                Found Items
              </Link>
              <Link
                  to="/report-item"
                  className="bg-blue-600 text-white px-4 py-2 rounded-md text-sm font-medium hover:bg-blue-700"
              >
                Report Item
              </Link>
              <button
                  onClick={logOut}
                  className=" text-black px-4 py-2 rounded-md text-sm font-medium hover:bg-blue-700 border-2 border-blue-500"
              >
                Logout
              </button>
            </div>
          </div>
        </div>
      </nav>
  );
};

export default Navbar;