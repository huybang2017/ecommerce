import { Component, OnInit, OnDestroy } from "@angular/core";
import { CategoryTableComponent } from "../../../shared/components/products/category/list/category-table.component";
import { ComponentCardComponent } from "../../../shared/components/common/component-card/component-card.component";
import { ActivatedRoute, Router } from "@angular/router";
import { CategoryQueryService } from "../../../core/queries/category-query.service";
import { CategorySortBy } from "../../../shared/models/category.model";
import { Subscription } from "rxjs";

@Component({
  selector: "category",
  imports: [CategoryTableComponent, ComponentCardComponent],
  templateUrl: "./category.component.html",
  styles: ``,
})
export class CategoryComponent implements OnInit, OnDestroy {
  private subs = new Subscription();

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private query: CategoryQueryService,
  ) {}

  ngOnInit() {
    this.subs.add(
      this.route.queryParams.subscribe((params) => {
        // Parse and validate params
        const page = Number(params["page"]);
        const page_size = Number(params["page_size"]);
        const search = params["search"] ?? "";
        const sort_by = params["sort_by"] ?? "name";
        const sort_order = params["sort_order"] ?? "asc";
        const status = params["status"] ?? "all";

        // Validate page and page_size
        const isPageValid = !isNaN(page) && page >= 1;
        const isPageSizeValid = !isNaN(page_size) && page_size > 0;

        // Redirect if params are missing or invalid
        if (!isPageValid || !isPageSizeValid) {
          this.router.navigate([], {
            relativeTo: this.route,
            queryParams: {
              page: isPageValid ? page : 1,
              page_size: isPageSizeValid ? page_size : 10,
              search: search || undefined,
              status,
              sort_by,
              sort_order,
            },
            queryParamsHandling: "merge",
            replaceUrl: true,
          });
          return;
        }

        // Sync all params to query service in one atomic call → single refetch
        this.query.applyFilters({
          page,
          page_size,
          search,
          sort_by: sort_by as CategorySortBy,
          sort_order: sort_order as "asc" | "desc",
          status: status as "all" | "active" | "inactive",
        });
      }),
    );
  }

  ngOnDestroy() {
    this.subs.unsubscribe();
  }
}
