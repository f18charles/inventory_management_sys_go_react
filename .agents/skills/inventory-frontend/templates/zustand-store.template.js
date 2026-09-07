// src/stores/authStore.js
//
// Zustand stores hold CLIENT state only — auth session, UI toggles,
// in-progress multi-step form state. Never cache server data
// (products, sales, inventory, etc.) here; that's TanStack Query's job.
//
// Shape convention: state and actions live together in one object,
// not raw setters exposed to components.

import { create } from 'zustand';
import { persist } from 'zustand/middleware';

export const useAuthStore = create(
  persist(
    (set) => ({
      user: null,
      token: null,

      login: (user, token) => set({ user, token }),

      logout: () => set({ user: null, token: null }),

      updateUser: (partial) =>
        set((state) => ({ user: { ...state.user, ...partial } })),
    }),
    {
      name: 'auth-storage', // localStorage key
      partialize: (state) => ({ user: state.user, token: state.token }),
    }
  )
);

// Example of a non-persisted, UI-only store — same pattern, no `persist`:
//
// export const useUiStore = create((set) => ({
//   isSidebarOpen: true,
//   activeModal: null,
//   toggleSidebar: () => set((s) => ({ isSidebarOpen: !s.isSidebarOpen })),
//   openModal: (name) => set({ activeModal: name }),
//   closeModal: () => set({ activeModal: null }),
// }));
