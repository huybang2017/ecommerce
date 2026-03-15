import { Injectable } from "@angular/core";
import { Observable } from "rxjs";
import { filter } from "rxjs/operators";
import {
  Product,
  ProductCreate,
  ProductUpdate,
} from "../../shared/models/product.model";
import { ProductService } from "../services/product.service";
import {
  QueryEngine,
  QueryInstance,
  cachePlugin,
  retryPlugin,
  loggingPlugin,
} from "../query-engine";

@Injectable({ providedIn: "root" })
export class ProductQueryService {
  private listQuery: QueryInstance<Product[]>;

  constructor(
    private engine: QueryEngine,
    private api: ProductService,
  ) {
    this.listQuery = this.engine.query<Product[]>({
      key: "products",
      fetchFn: () => this.api.list(),
      plugins: [
        cachePlugin({ ttl: 60_000, staleTime: 10_000 }),
        retryPlugin({ maxRetries: 2, delay: 1000 }),
        loggingPlugin({ verbose: false }),
      ],
    });
  }

  list$(): Observable<Product[]> {
    return this.listQuery.data$.pipe(filter((v): v is Product[] => v !== null));
  }

  loading$(): Observable<boolean> {
    return this.listQuery.loading$;
  }

  error$(): Observable<any> {
    return this.listQuery.error$;
  }

  refetch(): void {
    this.listQuery.refetch();
  }

  invalidate(): void {
    this.listQuery.invalidate();
  }

  snapshot(): Product[] | null {
    return this.listQuery.getSnapshot();
  }

  create(payload: ProductCreate): Observable<Product> {
    return this.engine
      .mutation<Product, ProductCreate>({
        mutationFn: (p) => this.api.create(p),
        invalidateKeys: ["products"],
      })
      .mutate(payload);
  }

  update(id: number, payload: ProductUpdate): Observable<Product> {
    return this.engine
      .mutation<Product, { id: number; payload: ProductUpdate }>({
        mutationFn: (v) => this.api.update(v.id, v.payload),
        invalidateKeys: ["products"],
      })
      .mutate({ id, payload });
  }

  delete(id: number): Observable<void> {
    const rollback = this.listQuery.optimisticUpdate((cur) =>
      (cur ?? []).filter((p) => p.id !== id),
    );

    return this.engine
      .mutation<void, number>({
        mutationFn: (pid) => this.api.remove(pid),
        invalidateKeys: ["products"],
        onError: () => rollback(),
      })
      .mutate(id);
  }
}
