// src/routes/_authenticated/products.jsx
//
// Route files stay thin: path definition, auth/guard concerns, and
// delegating to the actual page component. No business logic here.

import { createFileRoute, redirect } from '@tanstack/react-router';
import { useAuthStore } from '../../stores/authStore';
import { ProductsPage } from '../../pages/ProductsPage';

export const Route = createFileRoute('/_authenticated/products')({
  beforeLoad: () => {
    if (!useAuthStore.getState().token) {
      throw redirect({ to: '/login' });
    }
  },
  component: ProductsPage,
});

// For a route with a param, e.g. src/routes/_authenticated/products.$productId.jsx:
//
// export const Route = createFileRoute('/_authenticated/products/$productId')({
//   component: ProductDetailPage,
// });
//
// Inside ProductDetailPage: const { productId } = Route.useParams();
