import { Component, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { CoursesService } from '../../../courses/services/courses.service';
import { CourseStructureItem, Course } from '../../../courses/models/course.model';

@Component({
  selector: 'app-course-editor',
  standalone: true,
  imports: [CommonModule, RouterLink, ReactiveFormsModule],
  templateUrl: './course-editor.component.html',
  styleUrl: './course-editor.component.scss'
})
export class CourseEditorComponent implements OnInit {
  courseId!: number;
  course = signal<Course | null>(null);
  structure = signal<CourseStructureItem[]>([]);
  isLoading = signal<boolean>(true);
  error = signal<string | null>(null);

  // Modals
  isChapterModalOpen = signal<boolean>(false);
  isModuleModalOpen = signal<boolean>(false);
  activeChapterId = signal<number | null>(null);

  chapterForm: FormGroup;
  moduleForm: FormGroup;
  isSaving = signal<boolean>(false);

  expandedChapters = signal<Set<number>>(new Set());

  constructor(
    private route: ActivatedRoute,
    private coursesService: CoursesService,
    private fb: FormBuilder
  ) {
    this.chapterForm = this.fb.group({
      title: ['', Validators.required]
    });
    this.moduleForm = this.fb.group({
      title: ['', Validators.required]
    });
  }

  ngOnInit() {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.courseId = parseInt(id, 10);
      this.loadCourseData();
    }
  }

  loadCourseData() {
    this.isLoading.set(true);
    // Load course details
    this.coursesService.getCourse(this.courseId).subscribe({
      next: (c) => this.course.set(c),
      error: () => console.error('Failed to load course')
    });

    // Load structure
    this.coursesService.getCourseStructure(this.courseId).subscribe({
      next: (st) => {
        this.structure.set(st || []);
        // expand all chapters by default
        this.expandedChapters.set(new Set(st.map(s => s.chapter.id)));
        this.isLoading.set(false);
      },
      error: () => {
        this.error.set('Құрылымды жүктеу мүмкін болмады.');
        this.isLoading.set(false);
      }
    });
  }

  toggleChapter(chapterId: number) {
    const current = new Set(this.expandedChapters());
    if (current.has(chapterId)) {
      current.delete(chapterId);
    } else {
      current.add(chapterId);
    }
    this.expandedChapters.set(current);
  }

  // --- Chapter Flow ---
  openChapterModal() {
    this.chapterForm.reset();
    this.isChapterModalOpen.set(true);
  }

  closeChapterModal() {
    this.isChapterModalOpen.set(false);
  }

  saveChapter() {
    if (this.chapterForm.invalid) {
      this.chapterForm.markAllAsTouched();
      return;
    }
    
    this.isSaving.set(true);
    const orderIndex = this.structure().length + 1;
    const title = this.chapterForm.value.title;

    this.coursesService.addChapter({ course_id: this.courseId, title, order_index: orderIndex }).subscribe({
      next: () => {
        this.isSaving.set(false);
        this.closeChapterModal();
        this.loadCourseData();
      },
      error: () => {
        this.isSaving.set(false);
        alert('Қате пайда болды');
      }
    });
  }

  // --- Module Flow ---
  openModuleModal(chapterId: number) {
    this.activeChapterId.set(chapterId);
    this.moduleForm.reset();
    this.isModuleModalOpen.set(true);
  }

  closeModuleModal() {
    this.isModuleModalOpen.set(false);
    this.activeChapterId.set(null);
  }

  saveModule() {
    if (this.moduleForm.invalid) {
      this.moduleForm.markAllAsTouched();
      return;
    }

    const chapterId = this.activeChapterId();
    if (!chapterId) return;

    this.isSaving.set(true);
    
    // Find chapter to get order index
    const chapterItem = this.structure().find(s => s.chapter.id === chapterId);
    const orderIndex = chapterItem && chapterItem.modules ? chapterItem.modules.length + 1 : 1;
    const title = this.moduleForm.value.title;

    this.coursesService.addModule({ chapter_id: chapterId, title, order_index: orderIndex }).subscribe({
      next: () => {
        this.isSaving.set(false);
        this.closeModuleModal();
        // ensure chapter is expanded
        const expanded = new Set(this.expandedChapters());
        expanded.add(chapterId);
        this.expandedChapters.set(expanded);
        
        this.loadCourseData();
      },
      error: () => {
        this.isSaving.set(false);
        alert('Қате пайда болды');
      }
    });
  }
}
