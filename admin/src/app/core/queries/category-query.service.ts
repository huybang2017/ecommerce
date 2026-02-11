import { Injectable } from '@angular/core';
import { Observable, BehaviorSubject, combineLatest } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { Category, CategoryListParams, PaginatedResponse } from '../../shared/models/category.model';
import { CategoryService } from '../services/category.service';
import { QueryEngine, QueryInstance, cachePlugin, retryPlugin, loggingPlugin } from '../query-engine';

@Injectable({ providedIn: 'root' })
export class CategoryQueryService {
  private listQuery: QueryInstance<PaginatedResponse<Category>>;

  // For admin view: pagination + filtering state
  private currentPageSubject$ = new BehaviorSubject<number>(1);
  private pageSizeSubject$ = new BehaviorSubject<number>(10);
  private searchSubject$ = new BehaviorSubject<string>('');
  private isActiveSubject$ = new BehaviorSubject<boolean | null>(null);
  private sortBySubject$ = new BehaviorSubject<'name' | 'created_at' | 'updated_at'>('name');
  private sortOrderSubject$ = new BehaviorSubject<'asc' | 'desc'>('asc');

  constructor(
    private engine: QueryEngine,
    private api: CategoryService,
  ) {
    this.listQuery = this.engine.query<PaginatedResponse<Category>>({
      key: 'categories-admin',
      fetchFn: () =>
        this.api.list({
          page: this.currentPageSubject$.value,
          page_size: this.pageSizeSubject$.value,
          search: this.searchSubject$.value || undefined,
          is_active: this.isActiveSubject$.value,
          sort_by: this.sortBySubject$.value,
          sort_order: this.sortOrderSubject$.value,
        }),
      plugins: [
        cachePlugin({ ttl: 5 * 60_000, staleTime: 30_000 }),
        retryPlugin({ maxRetries: 2 }),
        loggingPlugin({ verbose: false }),
      ],
    });
  }

  // ─── Reactive Selectors ───────────────────────────────────

  list$(): Observable<PaginatedResponse<Category>> {
    return this.listQuery.data$.pipe(
      filter((v): v is PaginatedResponse<Category> => v !== null),
    );
  }

  loading$(): Observable<boolean> {
    return this.listQuery.loading$;
  }

  error$(): Observable<any> {
    return this.listQuery.error$;
  }

  // ─── Pagination State ─────────────────────────────────────

  getCurrentPage$(): Observable<number> {
    return this.currentPageSubject$.asObservable();
  }

  getPageSize$(): Observable<number> {
    return this.pageSizeSubject$.asObservable();
  }

  getSearch$(): Observable<string> {
    return this.searchSubject$.asObservable();
  }

  getSortBy$(): Observable<'name' | 'created_at' | 'updated_at'> {
    return this.sortBySubject$.asObservable();
  }

  getSortOrder$(): Observable<'asc' | 'desc'> {
    return this.sortOrderSubject$.asObservable();
  }

  // ─── Filter/Sort Actions ──────────────────────────────────

  setPage(page: number): void {
    if (page >= 1) {
      this.currentPageSubject$.next(page);
      this.refetch();
    }
  }

  setPageSize(size: number): void {
    if (size > 0) {
      this.pageSizeSubject$.next(size);
      this.currentPageSubject$.next(1); // Reset to page 1 when changing page size
      this.refetch();
    }
  }

  setSearch(query: string): void {
    this.searchSubject$.next(query);
    this.currentPageSubject$.next(1); // Reset pagination on search
    this.refetch();
  }

  setIsActive(active: boolean | null): void {
    this.isActiveSubject$.next(active);
    this.currentPageSubject$.next(1);
    this.refetch();
  }

  setSortBy(sortBy: 'name' | 'created_at' | 'updated_at'): void {
    this.sortBySubject$.next(sortBy);
    this.refetch();
  }

  setSortOrder(order: 'asc' | 'desc'): void {
    this.sortOrderSubject$.next(order);
    this.refetch();
  }

  // ─── Actions ──────────────────────────────────────────────

  refetch(): void {
    this.listQuery.refetch();
  }

  invalidate(): void {
    this.listQuery.invalidate();
  }

  snapshot(): PaginatedResponse<Category> | null {
    return this.listQuery.getSnapshot();
  }



  // ─── Mutations ────────────────────────────────────────────

  create(payload: Partial<Category>): Observable<Category> {
    return this.engine
      .mutation<Category, Partial<Category>>({
        mutationFn: (p) => this.api.create(p),
        invalidateKeys: ['categories-admin'],
      })
      .mutate(payload);
  }

  update(id: number, payload: Partial<Category>): Observable<Category> {
    return this.engine
      .mutation<Category, { id: number; payload: Partial<Category> }>({
        mutationFn: (v) => this.api.update(v.id, v.payload),
        invalidateKeys: ['categories-admin'],
      })
      .mutate({ id, payload });
  }

  delete(id: number): Observable<void> {
    const rollback = this.listQuery.optimisticUpdate((cur) => {
      const prev: PaginatedResponse<Category> =
        cur ?? {
          data: [],
          pagination: {
            page: this.currentPageSubject$.value,
            page_size: this.pageSizeSubject$.value,
            total: 0,
            total_pages: 0,
          },
        };

      return {
        ...prev,
        data: prev.data.filter((c) => c.id !== id),
      };
    });

    return this.engine
      .mutation<void, number>({
        mutationFn: (cid) => this.api.remove(cid),
        invalidateKeys: ['categories-admin'],
        onError: () => rollback(),
      })
      .mutate(id);
  }
}
