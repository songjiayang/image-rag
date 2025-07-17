export interface Record {
  id: number;
  name: string;
  description: string;
  images: Image[];
  created_at: string;
  updated_at: string;
}

export interface Image {
  id: number;
  record_id: number;
  filename: string;
  path: string;
  vector_id: string;
  created_at: string;
}

export interface SearchResult {
  record_id: number;
  record_name: string;
  description: string;
  image_id: number;
  filename: string;
  distance: number;
}

export interface SearchResponse {
  results: SearchResult[];
  count: number;
  message: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
}

export interface CreateRecordRequest {
  name: string;
  description?: string;
}

export interface UpdateRecordRequest {
  name?: string;
  description?: string;
}