import { CommonModule } from "@angular/common";
import { Component, OnInit, OnDestroy } from "@angular/core";
import { FormsModule } from "@angular/forms";
import { ButtonComponent } from "../../../ui/button/button.component";
import { TableDropdownComponent } from "../../../common/table-dropdown/table-dropdown.component";
import { BadgeComponent } from "../../../ui/badge/badge.component";
import { CategoryQueryService } from "../../../../../core/queries/category-query.service";
import { Category } from "../../../../models/category.model";
import { Observable, Subscription } from 'rxjs';
import { map } from 'rxjs/operators';
import { Router, ActivatedRoute } from '@angular/router';

@Component({
  selector: "category-table",
  imports: [
    CommonModule,
    FormsModule,
    ButtonComponent,
    TableDropdownComponent,
    BadgeComponent,
  ],
  templateUrl: "./category-table.component.html",
  styles: ``,
})
export class CategoryTableComponent implements OnInit {
  categories$!: Observable<Category[]>;
  loading$!: Observable<boolean>;
  error$!: Observable<any>;

  // Local copies of pagination & filter state for template bindings
  currentPage = 1;
  pageSize = 10;
  searchInput = '';
  isActiveFilter: boolean | null = null;
  sortBy: 'name' | 'created_at' | 'updated_at' = 'name';
  sortOrder: 'asc' | 'desc' = 'asc';

  private subs = new Subscription();

  constructor(
    private query: CategoryQueryService,
    private router: Router,
    private route: ActivatedRoute,
  ) {}

  ngOnInit(): void {
    this.categories$ = this.query.list$().pipe(map((res) => res.data));
    this.loading$ = this.query.loading$();
    this.error$ = this.query.error$();

    // Sync local state with query service observables
    this.subs.add(this.query.getCurrentPage$().subscribe((p) => (this.currentPage = p)));
    this.subs.add(this.query.getPageSize$().subscribe((s) => (this.pageSize = s)));
    this.subs.add(this.query.getSearch$().subscribe((q) => (this.searchInput = q)));
    this.subs.add(this.query.getSortBy$().subscribe((sb) => (this.sortBy = sb)));
    this.subs.add(this.query.getSortOrder$().subscribe((so) => (this.sortOrder = so)));
  }

  // ─── Pagination ───────────────────────────────────────────

  goToPage(page: number): void {
    if (page < 1) return;
    this.query.setPage(page);
    this.updateUrlParams({ page });
  }

  changePageSize(size: number): void {
    if (size <= 0) return;
    this.query.setPageSize(size);
    // reset to page 1 when changing page size
    this.updateUrlParams({ page: 1, page_size: size });
  }

  // ─── Search & Filter ──────────────────────────────────────

  onSearch(query: string): void {
    this.searchInput = query;
    this.query.setSearch(query);
    this.updateUrlParams({ search: query, page: 1 });
  }

  onFilterStatusChange(value: string | null): void {
    if (value === null || value === '') {
      this.isActiveFilter = null;
    } else {
      this.isActiveFilter = value === 'true';
    }
    this.query.setIsActive(this.isActiveFilter);
    const status = this.isActiveFilter === null ? 'all' : this.isActiveFilter ? 'active' : 'inactive';
    this.updateUrlParams({ status, page: 1 });
  }

  // ─── Sorting ──────────────────────────────────────────────

  onSort(sortBy: 'name' | 'created_at' | 'updated_at'): void {
    this.query.setSortBy(sortBy);
    this.updateUrlParams({ sort_by: sortBy });
  }

  toggleSortOrder(): void {
    const next = this.sortOrder === 'asc' ? 'desc' : 'asc';
    this.query.setSortOrder(next);
    this.updateUrlParams({ sort_order: next });
  }

  private updateUrlParams(params: Record<string, any>): void {
    this.router.navigate([], {
      relativeTo: this.route,
      queryParams: params,
      queryParamsHandling: 'merge',
    });
  }

  ngOnDestroy(): void {
    this.subs.unsubscribe();
  }

  // ─── Actions ──────────────────────────────────────────────

  handleViewMore(item: Category): void {
    console.log("View More Category:", item);
  }

  handleDelete(item: Category): void {
    if (confirm(`Delete category "${item.name}"?`)) {
      this.query.delete(item.id).subscribe({
        next: () => console.log('Deleted'),
        error: (e) => console.error(e),
      });
    }
  }

  // ─── Helpers ──────────────────────────────────────────────

  getStatusLabel(isActive: boolean): string {
    return isActive ? "Active" : "Inactive";
  }

  getBadgeColor(isActive: boolean): "success" | "warning" {
    return isActive ? "success" : "warning";
  }
}
