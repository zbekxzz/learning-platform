import { Component, OnInit, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink } from '@angular/router';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { CoursesService } from '../../../courses/services/courses.service';
import { Course } from '../../../courses/models/course.model';

@Component({
  selector: 'app-teacher-course-management',
  standalone: true,
  imports: [CommonModule, RouterLink, ReactiveFormsModule],
  templateUrl: './course-management.component.html',
  styleUrl: './course-management.component.scss'
})
export class TeacherCourseManagementComponent implements OnInit {
  courses = signal<Course[]>([]);
  isLoading = signal<boolean>(true);
  error = signal<string | null>(null);
  publishedCount = computed(() => this.courses().filter(c => c.is_published).length);
  hiddenCount = computed(() => this.courses().filter(c => !c.is_published).length);
  togglingId = signal<number | null>(null);

  isCreateModalOpen = signal<boolean>(false);
  isCreating = signal<boolean>(false);
  courseForm: FormGroup;

  isStatsModalOpen = signal<boolean>(false);
  statsLoading = signal<boolean>(false);
  statsData = signal<any[]>([]);
  selectedCourseTitle = signal<string>('');

  constructor(
    private coursesService: CoursesService,
    private fb: FormBuilder
  ) {
    this.courseForm = this.fb.group({
      title: ['', Validators.required],
      description: ['', Validators.required]
    });
  }

  ngOnInit() {
    this.loadCourses();
  }

  loadCourses() {
    this.isLoading.set(true);
    this.error.set(null);
    this.coursesService.getTeacherCourses().subscribe({
      next: (courses) => {
        this.courses.set(courses || []);
        this.isLoading.set(false);
      },
      error: (err) => {
        console.error('Failed to load courses', err);
        this.error.set('Курстар тізімін жүктеу мүмкін болмады.');
        this.isLoading.set(false);
      }
    });
  }

  togglePublish(course: Course) {
    this.togglingId.set(course.id);
    this.coursesService.togglePublish(course.id).subscribe({
      next: () => {
        this.courses.update(list =>
          list.map(c => c.id === course.id ? { ...c, is_published: !c.is_published } : c)
        );
        this.togglingId.set(null);
      },
      error: (err) => {
        console.error('Failed to toggle publish', err);
        this.togglingId.set(null);
        alert('Курс күйін өзгерту мүмкін болмады.');
      }
    });
  }

  openCreateForm() {
    this.courseForm.reset();
    this.isCreateModalOpen.set(true);
  }

  closeCreateForm() {
    this.isCreateModalOpen.set(false);
  }

  createCourse() {
    if (this.courseForm.invalid) {
      this.courseForm.markAllAsTouched();
      return;
    }

    this.isCreating.set(true);
    this.coursesService.createCourse(this.courseForm.value).subscribe({
      next: () => {
        this.isCreating.set(false);
        this.closeCreateForm();
        this.loadCourses(); // Reload the list
      },
      error: (err) => {
        console.error('Failed to create course', err);
        this.isCreating.set(false);
        alert('Курсты жасау кезінде қате пайда болды.');
      }
    });
  }

  openStatisticsModal(course: Course) {
    this.selectedCourseTitle.set(course.title);
    this.statsLoading.set(true);
    this.isStatsModalOpen.set(true);
    this.statsData.set([]);

    this.coursesService.getCourseStatistics(course.id).subscribe({
      next: (data) => {
        this.statsData.set(data || []);
        this.statsLoading.set(false);
      },
      error: (err) => {
        console.error('Failed to load course statistics', err);
        alert('Курс статистикасын жүктеу мүмкін болмады.');
        this.statsLoading.set(false);
        this.closeStatisticsModal();
      }
    });
  }

  closeStatisticsModal() {
    this.isStatsModalOpen.set(false);
  }
}
