import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { HttpApiService } from './http-api.service';
import { Product, ProductCreate, ProductUpdate } from '../../shared/models/product.model';

@Injectable({ providedIn: 'root' })
export class ProductService {
  private base = '/products';

  constructor(private http: HttpApiService) {}

  list(params?: Record<string, any>): Observable<Product[]> {
    return this.http.get<Product[]>(this.base, { params });
  }

  getById(id: number): Observable<Product> {
    return this.http.get<Product>(`${this.base}/${id}`);
  }

  create(payload: ProductCreate): Observable<Product> {
    return this.http.post<Product>(this.base, payload);
  }

  update(id: number, payload: ProductUpdate): Observable<Product> {
    return this.http.put<Product>(`${this.base}/${id}`, payload);
  }

  remove(id: number): Observable<void> {
    return this.http.delete<void>(`${this.base}/${id}`);
  }
}
