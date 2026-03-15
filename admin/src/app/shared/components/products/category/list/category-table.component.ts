import { CommonModule } from "@angular/common";
import { Component, OnInit, OnDestroy, inject } from "@angular/core";
import { FormsModule } from "@angular/forms";
import { ButtonComponent } from "../../../ui/button/button.component";
import { TableDropdownComponent } from "../../../common/table-dropdown/table-dropdown.component";
import { BadgeComponent } from "../../../ui/badge/badge.component";
import { CategoryQueryService } from "../../../../../core/queries/category-query.service";
import { Category, CategorySortBy } from "../../../../models/category.model";
import { Observable, Subscription, combineLatest } from 'rxjs';
import { map, startWith } from 'rxjs/operators';
import { Router, ActivatedRoute } from '@angular/router';
import { CategoryDetailModalComponent } from "../detail/category-detail-modal.component";
import { CategoryFormModalComponent } from "../form/category-form-modal.component";

@Component({
  selector: "category-table",
  imports: [
    CommonModule,
    FormsModule,
    ButtonComponent,
    TableDropdownComponent,
    BadgeComponent,
    CategoryDetailModalComponent,
    CategoryFormModalComponent,
  ],
  templateUrl: "./category-table.component.html",
  styles: ``,
})
export class CategoryTableComponent implements OnInit, OnDestroy {
  private readonly query = inject(CategoryQueryService);
  private readonly router = inject(Router);
  private readonly route = inject(ActivatedRoute);
  readonly categories$: Observable<Category[]> = this.query.list$().pipe(
    map((res) => res.data),
  );
  readonly loading$: Observable<boolean> = this.query.loading$();
  readonly error$: Observable<any> = this.query.error$();
  readonly currentPage$: Observable<number> = this.query.getCurrentPage$();
  readonly pageSize$: Observable<number> = this.query.getPageSize$();
  readonly sortBy$: Observable<CategorySortBy> = this.query.getSortBy$();
  readonly sortOrder$: Observable<'asc' | 'desc'> = this.query.getSortOrder$();
  readonly status$: Observable<'all' | 'active' | 'inactive'> = this.query.getStatus$();
  readonly pageState$ = combineLatest({
    currentPage: this.query.getCurrentPage$(),
    pageSize: this.query.getPageSize$(),
    totalPages: this.query.list$().pipe(map((r) => r.total_pages), startWith(0)),
    total: this.query.list$().pipe(map((r) => r.total), startWith(0)),
  });
  searchInput = '';
  selectedCategory: Category | null = null;
  isDetailModalOpen = false;
  editingCategory: Category | null = null;
  isFormModalOpen = false;
  formModalMode: 'create-parent' | 'create-child' | 'edit' = 'create-child';
  private readonly subs = new Subscription();

  ngOnInit(): void {
    this.subs.add(
      this.query.getSearch$().subscribe((q) => (this.searchInput = q)),
    );
  }

  ngOnDestroy(): void {
    this.subs.unsubscribe();
  }

  goToPage(page: number): void {
    if (page < 1) return;
    this.updateUrlParams({ page });
  }

  changePageSize(size: number): void {
    if (size <= 0) return;
    this.updateUrlParams({ page: 1, page_size: size });
  }

  onSearch(query: string): void {
    this.updateUrlParams({ search: query, page: 1 });
  }

  onFilterStatusChange(status: 'all' | 'active' | 'inactive'): void {
    this.updateUrlParams({ status, page: 1 });
  }

  onSort(sortBy: CategorySortBy): void {
    this.updateUrlParams({ sort_by: sortBy });
  }

  toggleSortOrder(current: 'asc' | 'desc'): void {
    this.updateUrlParams({ sort_order: current === 'asc' ? 'desc' : 'asc' });
  }

  private updateUrlParams(params: Record<string, any>): void {
    this.router.navigate([], {
      relativeTo: this.route,
      queryParams: params,
      queryParamsHandling: 'merge',
    });
  }

  handleViewMore(item: Category): void {
    this.selectedCategory = item;
    this.isDetailModalOpen = true;
  }

  closeDetailModal(): void {
    this.isDetailModalOpen = false;
    this.selectedCategory = null;
  }

  openCreateParentModal(): void {
    this.editingCategory = null;
    this.formModalMode = 'create-parent';
    this.isFormModalOpen = true;
  }

  openCreateChildModal(): void {
    this.editingCategory = null;
    this.formModalMode = 'create-child';
    this.isFormModalOpen = true;
  }

  openEditModal(item: Category): void {
    this.editingCategory = item;
    this.formModalMode = 'edit';
    this.isFormModalOpen = true;
  }

  closeFormModal(): void {
    this.isFormModalOpen = false;
    this.editingCategory = null;
  }

  handleDelete(item: Category): void {
    if (confirm(`Delete category "${item.name}"?`)) {
      this.query.delete(item.id).subscribe({
        next: () => console.log('Deleted'),
        error: (e) => console.error(e),
      });
    }
  }

  getStatusLabel(isActive: boolean): string {
    return isActive ? "Active" : "Inactive";
  }

  getBadgeColor(isActive: boolean): "success" | "warning" {
    return isActive ? "success" : "warning";
  }
}
