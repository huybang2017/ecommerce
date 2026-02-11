import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Category, PaginatedResponse } from '../../shared/models/category.model';
import { HttpApiService } from './http-api.service';

@Injectable({ providedIn: 'root' })
export class CategoryService {
  constructor(private api: HttpApiService) {}

  list(params?: any): Observable<PaginatedResponse<Category>> {
    return this.api.get<PaginatedResponse<Category>>('/categories/admin', { params });
  }

  getById(id: number): Observable<Category> {
    return this.api.get<Category>(`/categories/${id}`);
  }

  create(payload: Partial<Category>): Observable<Category> {
    return this.api.post<Category>('/categories', payload);
  }

  update(id: number, payload: Partial<Category>): Observable<Category> {
    return this.api.put<Category>(`/categories/${id}`, payload);
  }

  remove(id: number): Observable<void> {
    return this.api.delete<void>(`/categories/${id}`);
  }
}
