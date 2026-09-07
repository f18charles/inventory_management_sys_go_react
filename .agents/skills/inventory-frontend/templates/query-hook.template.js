// src/api/products.js
//
// One file per backend resource. Export query hooks (reads) and mutation
// hooks (writes) together so a component only imports from one place.
// Swap "products" / "Product" for the actual resource when using this.

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { api } from './client';

const productsKey = (params) => ['products', params ?? {}];
const productKey = (id) => ['products', id];

export function useProducts(params = {}) {
  return useQuery({
    queryKey: productsKey(params),
    queryFn: () => api.get('/products', { params }),
  });
}

export function useProduct(id) {
  return useQuery({
    queryKey: productKey(id),
    queryFn: () => api.get(`/products/${id}`),
    enabled: Boolean(id),
  });
}

export function useCreateProduct() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload) => api.post('/products', payload),
    onSuccess: () => {
      // Invalidate the resource's list queries only — not the whole cache.
      queryClient.invalidateQueries({ queryKey: ['products'] });
    },
  });
}

export function useUpdateProduct(id) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload) => api.patch(`/products/${id}`, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: productKey(id) });
      queryClient.invalidateQueries({ queryKey: ['products'] });
    },
  });
}

// Usage in a component:
//
// const { data: products, isLoading, isError, error } = useProducts({ categoryId });
// const createProduct = useCreateProduct();
//
// if (isLoading) return <Spinner />;
// if (isError) return <ErrorMessage error={error} />;
