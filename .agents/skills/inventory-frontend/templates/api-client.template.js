// src/api/client.js
//
// Single shared axios instance for the whole app. Nothing else in the
// codebase should import axios or call it directly — always go through
// this instance (or, better, through a resource hook in src/api/).

import axios from 'axios';
import { useAuthStore } from '../stores/authStore';

export class ApiError extends Error {
  constructor(code, message, status) {
    super(message);
    this.name = 'ApiError';
    this.code = code;
    this.status = status;
  }
}

export const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
});

// Attach the auth token, read fresh from the store on every request
// (not captured once at module load) so it stays correct after login/logout.
apiClient.interceptors.request.use((config) => {
  const token = useAuthStore.getState().token;
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Unwrap the backend's { data } envelope on success, and normalize
// { error: { code, message } } into an ApiError on failure — so calling
// code never touches axios's response/error shape directly.
apiClient.interceptors.response.use(
  (response) => response.data?.data,
  (error) => {
    const backendError = error.response?.data?.error;
    if (backendError) {
      return Promise.reject(
        new ApiError(backendError.code ?? 'UNKNOWN_ERROR', backendError.message, error.response.status)
      );
    }
    // Network error, timeout, or a non-JSON/unexpected failure.
    return Promise.reject(new ApiError('NETWORK_ERROR', error.message, error.response?.status));
  }
);

// Convenience wrappers — use these instead of calling apiClient.get/post
// directly so payload/method handling stays consistent everywhere.
export const api = {
  get: (path, config) => apiClient.get(path, config),
  post: (path, payload) => apiClient.post(path, payload),
  patch: (path, payload) => apiClient.patch(path, payload),
  delete: (path) => apiClient.delete(path),
};
