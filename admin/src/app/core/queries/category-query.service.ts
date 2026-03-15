import { Injectable } from "@angular/core";
import { Observable, BehaviorSubject } from "rxjs";
import { filter } from "rxjs/operators";
import {
  Category,
  CategoryListParams,
  CategorySortBy,
  PaginatedResponse,
} from "../../shared/models/category.model";
import { CategoryService } from "../services/category.service";
import {
  QueryEngine,
  QueryInstance,
  cachePlugin,
  retryPlugin,
  loggingPlugin,
} from "../query-engine";

@Injectable({ providedIn: "root" })
export class CategoryQueryService {
  private listQuery: QueryInstance<PaginatedResponse<Category>>;
  private currentPageSubject$ = new BehaviorSubject<number>(1);
  private pageSizeSubject$ = new BehaviorSubject<number>(10);
  private searchSubject$ = new BehaviorSubject<string>("");
  private statusSubject$ = new BehaviorSubject<"all" | "active" | "inactive">(
    "all",
  );
  private sortBySubject$ = new BehaviorSubject<CategorySortBy>("name");
  private sortOrderSubject$ = new BehaviorSubject<"asc" | "desc">("asc");

  constructor(
    private engine: QueryEngine,
    private api: CategoryService,
  ) {
    this.listQuery = this.engine.query<PaginatedResponse<Category>>({
      key: "categories-admin",
      fetchFn: () =>
        this.api.list({
          page: this.currentPageSubject$.value,
          limit: this.pageSizeSubject$.value,
          search: this.searchSubject$.value || undefined,
          status: this.statusSubject$.value,
          sort_by: this.sortBySubject$.value,
          sort_order: this.sortOrderSubject$.value,
        }),
      plugins: [
        cachePlugin({ ttl: 5 * 60_000, staleTime: 0 }),
        retryPlugin({ maxRetries: 2 }),
        loggingPlugin({ verbose: false }),
      ],
    });
  }

  list$(): Observable<PaginatedResponse<Category>> {
    return this.listQuery.data$.pipe(
      filter((v): v is PaginatedResponse<Category> => v !== null),
    );
  }

  listParentCategoriesAdmin$(): Observable<Category[]> {
    return this.api.listParentCategoriesAdmin();
  }

  loading$(): Observable<boolean> {
    return this.listQuery.loading$;
  }

  error$(): Observable<any> {
    return this.listQuery.error$;
  }

  getCurrentPage$(): Observable<number> {
    return this.currentPageSubject$.asObservable();
  }

  getPageSize$(): Observable<number> {
    return this.pageSizeSubject$.asObservable();
  }

  getSearch$(): Observable<string> {
    return this.searchSubject$.asObservable();
  }

  getStatus$(): Observable<"all" | "active" | "inactive"> {
    return this.statusSubject$.asObservable();
  }

  getSortBy$(): Observable<CategorySortBy> {
    return this.sortBySubject$.asObservable();
  }

  getSortOrder$(): Observable<"asc" | "desc"> {
    return this.sortOrderSubject$.asObservable();
  }

  applyFilters(params: {
    page: number;
    page_size: number;
    search: string;
    sort_by: CategorySortBy;
    sort_order: "asc" | "desc";
    status: "all" | "active" | "inactive";
  }): void {
    this.currentPageSubject$.next(params.page);
    this.pageSizeSubject$.next(params.page_size);
    this.searchSubject$.next(params.search);
    this.sortBySubject$.next(params.sort_by);
    this.sortOrderSubject$.next(params.sort_order);
    this.statusSubject$.next(params.status);
    this.refetch();
  }

  setPage(page: number): void {
    if (page >= 1) {
      this.currentPageSubject$.next(page);
      this.refetch();
    }
  }

  setPageSize(size: number): void {
    if (size > 0) {
      this.pageSizeSubject$.next(size);
      this.currentPageSubject$.next(1);
      this.refetch();
    }
  }

  setSearch(query: string): void {
    this.searchSubject$.next(query);
    this.currentPageSubject$.next(1);
    this.refetch();
  }

  setStatus(status: "all" | "active" | "inactive"): void {
    this.statusSubject$.next(status);
    this.currentPageSubject$.next(1);
    this.refetch();
  }

  setSortBy(sortBy: CategorySortBy): void {
    this.sortBySubject$.next(sortBy);
    this.refetch();
  }

  setSortOrder(order: "asc" | "desc"): void {
    this.sortOrderSubject$.next(order);
    this.refetch();
  }

  refetch(): void {
    this.listQuery.refetch();
  }

  invalidate(): void {
    this.listQuery.invalidate();
  }

  snapshot(): PaginatedResponse<Category> | null {
    return this.listQuery.getSnapshot();
  }

  create(payload: Partial<Category>): Observable<Category> {
    return this.engine
      .mutation<Category, Partial<Category>>({
        mutationFn: (p) => this.api.create(p),
        invalidateKeys: ["categories-admin"],
      })
      .mutate(payload);
  }

  createCategoryParent(payload: Partial<Category>): Observable<Category> {
    return this.engine
      .mutation<Category, Partial<Category>>({
        mutationFn: (p) => this.api.createCategoryParent(p),
        invalidateKeys: ["categories-admin"],
      })
      .mutate(payload);
  }

  update(id: number, payload: Partial<Category>): Observable<Category> {
    return this.engine
      .mutation<Category, { id: number; payload: Partial<Category> }>({
        mutationFn: (v) => this.api.update(v.id, v.payload),
        invalidateKeys: ["categories-admin"],
      })
      .mutate({ id, payload });
  }

  delete(id: number): Observable<void> {
    const rollback = this.listQuery.optimisticUpdate((cur) => {
      const prev: PaginatedResponse<Category> = cur ?? {
        data: [],
        limit: this.pageSizeSubject$.value,
        page: this.currentPageSubject$.value,
        total: 0,
        total_pages: 0,
      };

      return {
        ...prev,
        data: prev.data.filter((c) => c.id !== id),
      };
    });

    return this.engine
      .mutation<void, number>({
        mutationFn: (cid) => this.api.remove(cid),
        invalidateKeys: ["categories-admin"],
        onError: () => rollback(),
      })
      .mutate(id);
  }
}
