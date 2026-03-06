import React from 'react';
import { Link } from 'react-router-dom';
import { AlertCircle, CheckCircle } from 'lucide-react';

const Home: React.FC = () => {


  return (
    <div className="space-y-12">
      <section className="text-center space-y-4">
        <h1 className="text-4xl font-bold text-gray-900">
          Lost Something? Found Something?
        </h1>
        <p className="text-xl text-gray-600 max-w-2xl mx-auto">
          Our college lost and found portal helps connect lost items with their owners.
        </p>
      </section>

      <div className="grid md:grid-cols-2 gap-8 max-w-4xl mx-auto">
        <Link
          to="/lost-items"
          className="p-6 bg-white rounded-xl shadow-md hover:shadow-lg transition-shadow"
        >
          <div className="flex items-center space-x-4">
            <AlertCircle className="h-8 w-8 text-red-500" />
            <div>
              <h2 className="text-xl font-semibold text-gray-900">Report Lost Item</h2>
              <p className="text-gray-600">List your lost item and get notified when found</p>
            </div>
          </div>
        </Link>

        <Link
          to="/found-items"
          className="p-6 bg-white rounded-xl shadow-md hover:shadow-lg transition-shadow"
        >
          <div className="flex items-center space-x-4">
            <CheckCircle className="h-8 w-8 text-green-500" />
            <div>
              <h2 className="text-xl font-semibold text-gray-900">Report Found Item</h2>
              <p className="text-gray-600">Help others by reporting items you've found</p>
            </div>
          </div>
        </Link>
      </div>

      <section className="bg-white rounded-xl shadow-md p-8 max-w-4xl mx-auto">
        <h2 className="text-2xl font-semibold text-gray-900 mb-6">How It Works</h2>
        <div className="grid md:grid-cols-3 gap-6">
          <div className="text-center">
            <div className="bg-blue-100 rounded-full w-12 h-12 flex items-center justify-center mx-auto mb-4">
              <span className="text-blue-600 font-bold">1</span>
            </div>
            <h3 className="font-semibold mb-2">Report</h3>
            <p className="text-gray-600">Submit details about lost or found items</p>
          </div>
          <div className="text-center">
            <div className="bg-blue-100 rounded-full w-12 h-12 flex items-center justify-center mx-auto mb-4">
              <span className="text-blue-600 font-bold">2</span>
            </div>
            <h3 className="font-semibold mb-2">Match</h3>
            <p className="text-gray-600">We'll help match lost items with found ones</p>
          </div>
          <div className="text-center">
            <div className="bg-blue-100 rounded-full w-12 h-12 flex items-center justify-center mx-auto mb-4">
              <span className="text-blue-600 font-bold">3</span>
            </div>
            <h3 className="font-semibold mb-2">Reunite</h3>
            <p className="text-gray-600">Get connected to retrieve your belongings</p>
          </div>
        </div>
      </section>
    </div>
  );
};

export default Home;