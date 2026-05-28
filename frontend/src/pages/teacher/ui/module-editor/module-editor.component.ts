import { Component, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { CdkDragDrop, DragDropModule, moveItemInArray } from '@angular/cdk/drag-drop';
import { CoursesService } from '../../../courses/services/courses.service';
import { ModuleMaterial } from '../../../courses/models/course.model';

@Component({
  selector: 'app-module-editor',
  standalone: true,
  imports: [CommonModule, RouterLink, ReactiveFormsModule, DragDropModule],
  templateUrl: './module-editor.component.html',
  styleUrl: './module-editor.component.scss'
})
export class ModuleEditorComponent implements OnInit {
  courseId!: number;
  moduleId!: number;

  materials = signal<ModuleMaterial[]>([]);
  activeMaterial = signal<ModuleMaterial | null>(null);

  isLoading = signal<boolean>(true);
  error = signal<string | null>(null);

  isMaterialModalOpen = signal<boolean>(false);
  materialForm: FormGroup;
  isSaving = signal<boolean>(false);

  constructor(
    private route: ActivatedRoute,
    private coursesService: CoursesService,
    private fb: FormBuilder
  ) {
    this.materialForm = this.fb.group({
      title: ['', Validators.required],
      content: ['']
    });
  }

  ngOnInit() {
    this.courseId = parseInt(this.route.snapshot.paramMap.get('id') || '0', 10);
    this.moduleId = parseInt(this.route.snapshot.paramMap.get('moduleId') || '0', 10);
    
    if (this.moduleId) {
      this.loadMaterials();
    }
  }

  loadMaterials() {
    this.isLoading.set(true);
    this.coursesService.getModuleMaterials(this.moduleId).subscribe({
      next: (mats) => {
        this.materials.set(mats || []);
        if (mats && mats.length > 0 && !this.activeMaterial()) {
          this.activeMaterial.set(mats[0]);
        }
        this.isLoading.set(false);
      },
      error: () => {
        this.error.set('Материалдарды жүктеу мүмкін болмады.');
        this.isLoading.set(false);
      }
    });
  }

  selectMaterial(mat: ModuleMaterial) {
    this.activeMaterial.set(mat);
  }

  drop(event: CdkDragDrop<ModuleMaterial[]>) {
    const currentMaterials = [...this.materials()];
    moveItemInArray(currentMaterials, event.previousIndex, event.currentIndex);
    
    // Optimistic UI update
    this.materials.set(currentMaterials);

    // Prepare batch update payload
    const updates = currentMaterials.map((mat, index) => ({
      id: mat.id,
      order_index: index + 1
    }));

    // Send to backend
    this.coursesService.reorderMaterials(this.moduleId, updates).subscribe({
      next: () => {
        console.log('Materials reordered successfully');
      },
      error: (err) => {
        console.error('Failed to reorder materials', err);
        // Revert UI on failure
        this.loadMaterials();
      }
    });
  }

  // --- Add Material Flow ---
  openMaterialModal() {
    this.materialForm.reset();
    this.isMaterialModalOpen.set(true);
  }

  closeMaterialModal() {
    this.isMaterialModalOpen.set(false);
  }

  saveMaterial() {
    if (this.materialForm.invalid) {
      this.materialForm.markAllAsTouched();
      return;
    }

    this.isSaving.set(true);
    const orderIndex = this.materials().length + 1;
    const { title, content } = this.materialForm.value;

    this.coursesService.addMaterial({
      module_id: this.moduleId,
      title,
      type: 'text',
      content: content || '',
      order_index: orderIndex
    }).subscribe({
      next: (newMat) => {
        this.isSaving.set(false);
        this.closeMaterialModal();
        // Option 1: push directly
        // this.materials.update(list => [...list, newMat]);
        // Option 2: reload to ensure consistency
        this.loadMaterials();
        this.activeMaterial.set(newMat);
      },
      error: () => {
        this.isSaving.set(false);
        alert('Қате пайда болды');
      }
    });
  }
}
