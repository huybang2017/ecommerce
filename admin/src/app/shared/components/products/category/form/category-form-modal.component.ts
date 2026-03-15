import { CommonModule } from "@angular/common";
import {
  Component,
  EventEmitter,
  Input,
  OnChanges,
  Output,
  SimpleChanges,
  inject,
} from "@angular/core";
import { ModalComponent } from "../../../ui/modal/modal.component";
import { ButtonComponent } from "../../../ui/button/button.component";
import { InputFieldComponent } from "../../../form/input/input-field.component";
import { TextAreaComponent } from "../../../form/input/text-area.component";
import { SwitchComponent } from "../../../form/input/switch.component";
import { LabelComponent } from "../../../form/label/label.component";
import { SelectComponent, Option } from "../../../form/select/select.component";
import { Category } from "../../../../models/category.model";
import { CategoryQueryService } from "../../../../../core/queries/category-query.service";

interface CategoryForm {
  name: string;
  slug: string;
  description: string;
  image_url: string;
  parent_id: number | null;
  is_active: boolean;
}

@Component({
  selector: "category-form-modal",
  imports: [
    CommonModule,
    ModalComponent,
    ButtonComponent,
    InputFieldComponent,
    TextAreaComponent,
    SwitchComponent,
    LabelComponent,
    SelectComponent,
  ],
  templateUrl: "./category-form-modal.component.html",
  styles: ``,
})
export class CategoryFormModalComponent implements OnChanges {
  private readonly query = inject(CategoryQueryService);

  @Input() isOpen = false;
  @Input() category: Category | null = null;
  @Input() mode: 'create-parent' | 'create-child' | 'edit' = 'create-child';
  @Output() close = new EventEmitter<void>();

  form: CategoryForm = this.emptyForm();
  saving = false;
  errors: Partial<Record<keyof CategoryForm, string>> = {};
  parentCategoryOptions: Option[] = [];
  loadingParents = false;
  private isSlugDirty = false;

  constructor(private CategoryQueryService: CategoryQueryService) {}

  get isCreatingParent(): boolean {
    return this.mode === 'create-parent';
  }

  get isEditing(): boolean {
    return this.mode === 'edit';
  }

  get modalTitle(): string {
    if (this.isEditing) return 'Edit Category';
    if (this.isCreatingParent) return 'Create Parent Category';
    return 'Create Category';
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (changes["isOpen"]?.currentValue === true) {
      this.form = this.category
        ? this.fromCategory(this.category)
        : this.emptyForm();
      this.saving = false;
      this.errors = {};
      this.isSlugDirty = !!this.category;
      if (this.isCreatingParent) {
        this.form.parent_id = null;
        this.parentCategoryOptions = [];
      } else {
        this.loadParentCategories();
      }
    }
  }

  loadParentCategories(): void {
    this.loadingParents = true;
    this.CategoryQueryService.listParentCategoriesAdmin$().subscribe({
      next: (data) => {
        this.parentCategoryOptions = [
          { label: "Root Category", value: "" },
          ...data.map((cat) => ({
            label: cat.name,
            value: String(cat.id),
          })),
        ];
        this.loadingParents = false;
      },
      error: (err) => {
        console.error("Failed to load parent categories:", err);
        this.loadingParents = false;
      },
    });
  }

  get parentValue(): string {
    return this.form.parent_id ? String(this.form.parent_id) : "";
  }

  onNameChange(value: string | number): void {
    this.form.name = String(value);
    if (!this.isSlugDirty) {
      this.form.slug = this.slugify(this.form.name);
    }
  }

  onSlugChange(value: string | number): void {
    this.form.slug = String(value);
    this.isSlugDirty = true;
  }

  onParentIdChange(value: string): void {
    const num = Number(value);
    this.form.parent_id = value !== "" && !isNaN(num) && num > 0 ? num : null;
  }

  onSubmit(): void {
    if (!this.validate()) return;

    this.saving = true;
    const payload: Partial<Category> = {
      name: this.form.name,
      slug: this.form.slug || undefined,
      description: this.form.description,
      image_url: this.form.image_url || undefined,
      parent_id: this.isCreatingParent ? null : this.form.parent_id,
      is_active: this.form.is_active,
    };

    let obs$;
    if (this.isEditing) {
      obs$ = this.query.update(this.category!.id, payload);
    } else if (this.isCreatingParent) {
      obs$ = this.query.createCategoryParent(payload);
    } else {
      obs$ = this.query.create(payload);
    }

    obs$.subscribe({
      next: () => {
        this.saving = false;
        this.close.emit();
      },
      error: (err) => {
        console.error('Failed to save category:', err);
        this.saving = false;
      },
    });
  }

  onClose(): void {
    if (!this.saving) {
      this.close.emit();
    }
  }

  private validate(): boolean {
    this.errors = {};
    if (!this.form.name.trim()) {
      this.errors["name"] = "Name is required.";
    }
    return Object.keys(this.errors).length === 0;
  }

  private slugify(text: string): string {
    return text
      .toLowerCase()
      .trim()
      .replace(/[^\w\s-]/g, "")
      .replace(/[\s_]+/g, "-")
      .replace(/^-+|-+$/g, "");
  }

  private emptyForm(): CategoryForm {
    return {
      name: "",
      slug: "",
      description: "",
      image_url: "",
      parent_id: null,
      is_active: true,
    };
  }

  private fromCategory(c: Category): CategoryForm {
    return {
      name: c.name,
      slug: c.slug,
      description: c.description,
      image_url: c.image_url ?? "",
      parent_id: c.parent_id ?? null,
      is_active: c.is_active,
    };
  }
}
