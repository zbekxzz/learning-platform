import { Component, OnInit, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { CoursesService } from '../../../courses/services/courses.service';
import { CourseWithAuthor } from '../../../courses/models/course.model';

@Component({
  selector: 'app-admin-course-management',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './admin-course-management.component.html',
  styleUrl: './admin-course-management.component.scss'
})
export class AdminCourseManagementComponent implements OnInit {
  courses = signal<CourseWithAuthor[]>([]);
  isLoading = signal<boolean>(true);
  error = signal<string | null>(null);
  publishedCount = computed(() => this.courses().filter(c => c.is_published).length);
  hiddenCount = computed(() => this.courses().filter(c => !c.is_published).length);
  togglingId = signal<number | null>(null);

  constructor(private coursesService: CoursesService) {}

  ngOnInit() {
    this.loadCourses();
  }

  loadCourses() {
    this.isLoading.set(true);
    this.error.set(null);
    this.coursesService.getAllForAdmin().subscribe({
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

  togglePublish(course: CourseWithAuthor) {
    this.togglingId.set(course.id);
    this.coursesService.togglePublish(course.id).subscribe({
      next: () => {
        // Update locally without reloading
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
}
