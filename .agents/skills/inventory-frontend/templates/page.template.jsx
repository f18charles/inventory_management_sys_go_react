// src/pages/ProductsPage.jsx
//
// Route-level page: composes query hooks + reusable components.
// Always handles loading and error states explicitly.

import { useState } from 'react';
import { useProducts } from '../api/products';
import { ProductTable } from '../components/ProductTable';
import { ProductFilters } from '../components/ProductFilters';
import { Spinner } from '../components/Spinner';
import { ErrorMessage } from '../components/ErrorMessage';

export function ProductsPage() {
  const [categoryId, setCategoryId] = useState(null);
  const { data: products, isLoading, isError, error } = useProducts({ categoryId });

  return (
    <div className="p-6 space-y-4">
      <h1 className="text-2xl font-semibold">Products</h1>

      <ProductFilters categoryId={categoryId} onCategoryChange={setCategoryId} />

      {isLoading && <Spinner />}
      {isError && <ErrorMessage error={error} />}
      {!isLoading && !isError && <ProductTable products={products} />}
    </div>
  );
}
