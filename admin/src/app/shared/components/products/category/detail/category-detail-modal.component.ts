import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { ModalComponent } from '../../../ui/modal/modal.component';
import { BadgeComponent } from '../../../ui/badge/badge.component';
import { Category } from '../../../../models/category.model';

@Component({
  selector: 'category-detail-modal',
  imports: [CommonModule, ModalComponent, BadgeComponent],
  templateUrl: './category-detail-modal.component.html',
  styles: ``,
})
export class CategoryDetailModalComponent {
  @Input() isOpen = false;
  @Input() category: Category | null = null;
  @Output() close = new EventEmitter<void>();

  getBadgeColor(isActive: boolean): 'success' | 'warning' {
    return isActive ? 'success' : 'warning';
  }
}
