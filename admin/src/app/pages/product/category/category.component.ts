import { Component, OnInit, OnDestroy } from "@angular/core";
import { CategoryTableComponent } from "../../../shared/components/products/category/list/category-table.component";
import { ComponentCardComponent } from "../../../shared/components/common/component-card/component-card.component";
import { ActivatedRoute, Router } from "@angular/router";
import { CategoryQueryService } from "../../../core/queries/category-query.service";
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
        // Set defaults for missing params
        if (!params["page"] || !params["page_size"]) {
          this.router.navigate([], {
            relativeTo: this.route,
            queryParams: {
              page: Number(params["page"]) || 1,
              page_size: Number(params["page_size"]) || 10,
              status: params["status"] || "all",
              sort_by: params["sort_by"] || "name",
              sort_order: params["sort_order"] || "asc",
            },
            queryParamsHandling: "merge",
            replaceUrl: true,
          });
          return;
        }

        // Sync params into the query service
        const page = Number(params["page"]) || 1;
        const page_size = Number(params["page_size"]) || 10;
        const search = params["search"] ?? "";
        const sort_by = (params["sort_by"] ?? "name") as "name" | "created_at" | "updated_at";
        const sort_order = (params["sort_order"] ?? "asc") as "asc" | "desc";
        const status = params["status"] ?? "all";

        this.query.setPage(page);
        this.query.setPageSize(page_size);
        this.query.setSearch(search);
        this.query.setSortBy(sort_by);
        this.query.setSortOrder(sort_order);

        // Map status to isActive filter
        let isActive: boolean | null = null;
        if (status === "active") isActive = true;
        else if (status === "inactive") isActive = false;
        this.query.setIsActive(isActive);
      }),
    );
  }

  ngOnDestroy() {
    this.subs.unsubscribe();
  }
}
