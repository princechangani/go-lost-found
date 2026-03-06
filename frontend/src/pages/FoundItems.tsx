import React, {useEffect, useState} from 'react';
import { Filter } from 'lucide-react';
import FilterPanel, {FilterOptions} from "@/components/FilterPanel.tsx";
import {useSearch} from "@/context/SearchContext.tsx";
import {apiClient} from "@/api/apiClient.ts";
import {API_ENDPOINTS} from "@/api/endpoints.ts";
import ContactDialog from "@/components/ContactDialog.tsx";
import {Category} from "@/types/Category.ts";
import {useRefresh} from "@/context/RefreshContext.tsx";
import {Items, ItemsResponse} from "@/types/Item.ts";

interface LostItemsProps {
    categories: Category[]
}

const FoundItems: React.FC<LostItemsProps> = ({categories}) => {

    const { refresh } = useRefresh();


    const [allFoundItems, setAllFoundItems] = useState<Items[]>([]);
    const [filteredItems, setFilteredItems] = useState<Items[]>([]);
    const [openDialog, setOpenDialog] = React.useState(false);
    const [selectedItem, setSelectedItem] = useState<Items | null>(null);
  const [isFilterOpen, setIsFilterOpen] = useState(false);
  const [activeFilters, setActiveFilters] = useState<FilterOptions>({
    category: 'All',
    dateFrom: '',
    dateTo: '',
    status: 'All'
  });
  const { searchQuery } = useSearch();



  const handleFilterApply = (filters: FilterOptions) => {
    setActiveFilters(filters);
  };
  useEffect(() => {
    apiClient.get<ItemsResponse>(API_ENDPOINTS.GET_ALL_ITEMS)
        .then((res) => {
          const foundItems = res.data.data
              .filter((item) => item.type.toString().toLowerCase() === 'found')
              .map((item) => ({
                  ...item,
                  date: item.date ? new Date(item.date).toLocaleDateString('en-CA') : '',
              }));
            setAllFoundItems(foundItems);

        })
        .catch((err) => console.error(err));
  }, [refresh]);


    useEffect(() => {
        let filtered = [...allFoundItems];

        // Apply category filter
        if (activeFilters.category !== 'All') {
            filtered = filtered.filter(item => item.categoryId === activeFilters.category);
        }

        // Apply date filter
        if (activeFilters.dateFrom && activeFilters.dateTo) {
            const from = new Date(activeFilters.dateFrom);
            const to = new Date(activeFilters.dateTo);
            filtered = filtered.filter(item => {
                const itemDate = new Date(item.date);
                return itemDate >= from && itemDate <= to;
            });
        }

        // Apply status filter
        if (activeFilters.status !== 'All') {
            filtered = filtered.filter(item => item.statusType === activeFilters.status.toLowerCase());
        }

        // Apply search query
        if (searchQuery) {
            const lowerSearch = searchQuery.toLowerCase();
            filtered = filtered.filter(item =>
                item.title.toLowerCase().includes(lowerSearch) ||
                item.description.toLowerCase().includes(lowerSearch) ||
                item.location.toLowerCase().includes(lowerSearch)
            );
        }

        setFilteredItems(filtered);
    }, [allFoundItems, activeFilters, searchQuery]);


  return (
      <div className="space-y-6">
        <div className="flex justify-between items-center">
          <h1 className="text-2xl font-bold text-gray-900">Found Items</h1>
          <button
              onClick={() => setIsFilterOpen(true)}
              className="flex items-center space-x-2 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
          >
            <Filter className="h-4 w-4" />
            <span>Filter</span>
          </button>
        </div>

        {/* Search results count */}
        {searchQuery && (
            <div className="text-sm text-gray-600">
              Found {filteredItems.length} results for "{searchQuery}"
            </div>
        )}

        {/* Active filters display */}
        {(activeFilters.category !== 'All' ||
            activeFilters.dateFrom ||
            activeFilters.dateTo ||
            activeFilters.status !== 'All') && (
            <div className="bg-blue-50 p-4 rounded-md">
              <h3 className="text-sm font-medium text-blue-800 mb-2">Active Filters:</h3>
              <div className="flex flex-wrap gap-2">
                {activeFilters.category !== 'All' && (
                    <span className="px-3 py-1 bg-blue-100 text-blue-800 rounded-full text-sm">
                Category: {activeFilters.category}
              </span>
                )}
                {activeFilters.dateFrom && activeFilters.dateTo && (
                    <span className="px-3 py-1 bg-blue-100 text-blue-800 rounded-full text-sm">
                Date: {activeFilters.dateFrom} to {activeFilters.dateTo}
              </span>
                )}
                {activeFilters.status !== 'All' && (
                    <span className="px-3 py-1 bg-blue-100 text-blue-800 rounded-full text-sm">
                Status: {activeFilters.status}
              </span>
                )}
              </div>
            </div>
        )}

        {filteredItems.length === 0 ? (
            <div className="text-center py-12">
              <p className="text-gray-500 text-lg">No items found matching your criteria</p>
            </div>
        ) : (
            <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
              {filteredItems.map((item) => (
                  <div key={item.id} className="bg-white rounded-lg shadow-md overflow-hidden">
                    <img
                        src={item.image}
                        alt={item.title}
                        className="w-full h-48 object-cover"
                    />
                    <div className="p-6">
                      <div className="flex justify-between items-start">
                        <h3 className="text-lg font-semibold text-gray-900">{item.title}</h3>
                        <span
                            className={`px-2 py-1 rounded-full text-xs font-medium ${
                                item.statusType === 'claimed'
                                    ? 'bg-green-100 text-green-800'
                                    : 'bg-blue-100 text-blue-800'
                            }`}
                        >
                    {item.statusType === 'claimed' ? 'Claimed' : 'Unclaimed'}
                  </span>
                      </div>
                      <p className="mt-2 text-gray-600">{item.description}</p>
                      <div className="mt-4 space-y-2">
                        <div className="flex items-center text-sm text-gray-500">
                          <span className="font-medium">Location:</span>
                          <span className="ml-2">{item.location}</span>
                        </div>
                          <div className="flex items-center text-sm text-gray-500 pt-2" >

                              <div className="flex items-center text-sm text-gray-500 pr-2">
                                  <span className="font-medium">Date:</span>
                                  <span className="ml-2">{item.date}</span>
                              </div>
                              <div className="flex items-center text-sm text-gray-500">
                                  <span className="font-medium">Time:</span>
                                  <span className="ml-2">{item.time}</span>
                              </div>
                          </div>
                        <div className="flex items-center text-sm text-gray-500">
                          <span className="font-medium">Category:</span>
                          <span className="ml-2">{categories.find((category) => category.id === item.categoryId)?.name}</span>
                        </div>
                      </div>
                      <button className="mt-4 w-full bg-blue-600 text-white py-2 rounded-md hover:bg-blue-700"
                              onClick={() => {
                                  setSelectedItem(item);
                                  setOpenDialog(true);
                              }}
                      >
                        Claim Item
                      </button>
                    </div>
                  </div>
              ))}
            </div>
        )}
          {selectedItem && (
              <ContactDialog
                  item={selectedItem}
                  open={openDialog}
                  onClose={() => {
                      setOpenDialog(false);
                      setSelectedItem(null);
                  }}
              />
          )}
        <FilterPanel
            isOpen={isFilterOpen}
            onClose={() => setIsFilterOpen(false)}
            onApplyFilters={handleFilterApply}
            initialFilters={activeFilters}
            categories={categories}
        />
      </div>
  );
};
export default FoundItems;

