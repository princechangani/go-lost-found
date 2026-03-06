import React from 'react';

const Footer: React.FC = () => {
  return (
    <footer className="bg-white shadow-md mt-auto">
      <div className="container mx-auto px-4 py-6">
        <div className="flex flex-col md:flex-row justify-between items-center">
          <div className="text-gray-600 text-sm">
            © 2024 College Lost & Found. All rights reserved.
          </div>
          <div className="flex space-x-6 mt-4 md:mt-0">
            <a href="#" className="text-gray-600 hover:text-blue-600 text-sm">
              Privacy Policy
            </a>
            <a href="#" className="text-gray-600 hover:text-blue-600 text-sm">
              Terms of Service
            </a>
            <a href="#" className="text-gray-600 hover:text-blue-600 text-sm">
              Contact Us
            </a>
          </div>
        </div>
      </div>
    </footer>
  );
};

export default Footer;