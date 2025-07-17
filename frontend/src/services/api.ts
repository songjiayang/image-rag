import axios from 'axios';
import type { Record, SearchResponse, PaginatedResponse, CreateRecordRequest, UpdateRecordRequest } from '@/types';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor
api.interceptors.request.use(
  (config) => {
    // Add any auth tokens here if needed
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor
api.interceptors.response.use(
  (response) => response,
  (error) => {
    console.error('API Error:', error.response?.data || error.message);
    return Promise.reject(error);
  }
);

export const recordService = {
  async createRecord(data: CreateRecordRequest, images: File[]): Promise<Record> {
    const formData = new FormData();
    formData.append('name', data.name);
    if (data.description) {
      formData.append('description', data.description);
    }
    
    images.forEach((image, index) => {
      formData.append('images', image);
    });

    const response = await api.post<Record>('/records', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
    return response.data;
  },

  async getRecords(page: number = 1, limit: number = 10): Promise<PaginatedResponse<Record>> {
    const response = await api.get<PaginatedResponse<Record>>('/records', {
      params: { page, limit },
    });
    return response.data;
  },

  async getRecord(id: number): Promise<Record> {
    const response = await api.get<Record>(`/records/${id}`);
    return response.data;
  },

  async updateRecord(id: number, data: UpdateRecordRequest): Promise<Record> {
    const response = await api.put<Record>(`/records/${id}`, data);
    return response.data;
  },

  async deleteRecord(id: number): Promise<{ message: string }> {
    const response = await api.delete<{ message: string }>(`/records/${id}`);
    return response.data;
  },

  async addImageToRecord(recordId: number, image: File): Promise<Record['images'][0]> {
    const formData = new FormData();
    formData.append('image', image);

    const response = await api.post<Record['images'][0]>(`/records/${recordId}/images`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
    return response.data;
  },

  async deleteImage(imageId: number): Promise<{ message: string }> {
    const response = await api.delete<{ message: string }>(`/images/${imageId}`);
    return response.data;
  },
};

export const searchService = {
  async searchImages(image: File, topK: number = 10): Promise<SearchResponse> {
    const formData = new FormData();
    formData.append('image', image);

    const response = await api.post<SearchResponse>('/search', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
      params: { top_k: topK },
    });
    return response.data;
  },

  async findSimilar(imageId: number, topK: number = 10): Promise<SearchResponse> {
    const response = await api.get<SearchResponse>(`/search/similar/${imageId}`, {
      params: { top_k: topK },
    });
    return response.data;
  },

  async advancedSearch(image: File, params: {
    q?: string;
    record_name?: string;
    min_distance?: number;
    max_distance?: number;
    top_k?: number;
  }): Promise<SearchResponse> {
    const formData = new FormData();
    formData.append('image', image);

    const response = await api.post<SearchResponse>('/search/advanced', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
      params,
    });
    return response.data;
  },
};

export const healthService = {
  async checkHealth(): Promise<{ status: string }> {
    const response = await api.get<{ status: string }>('/health');
    return response.data;
  },
};

export default api;