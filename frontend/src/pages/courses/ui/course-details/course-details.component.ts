import { Component, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { CoursesService } from '../../services/courses.service';
import { Course, CourseModule } from '../../models/course.model';
import { forkJoin } from 'rxjs';

@Component({
  selector: 'app-course-details',
  imports: [CommonModule, RouterLink],
  templateUrl: './course-details.component.html',
  styleUrl: './course-details.component.scss'
})
export class CourseDetailsComponent implements OnInit {
  course = signal<Course | null>(null);
  modules = signal<CourseModule[]>([]);
  isLoading = signal<boolean>(true);
  error = signal<string | null>(null);

  constructor(
    private route: ActivatedRoute,
    private coursesService: CoursesService
  ) {}

  ngOnInit() {
    this.route.paramMap.subscribe(params => {
      const id = params.get('id');
      if (id) {
        this.loadCourseData(+id);
      }
    });
  }

  loadCourseData(id: number) {
    this.isLoading.set(true);
    this.error.set(null);

    forkJoin({
      course: this.coursesService.getCourse(id),
      modules: this.coursesService.getCourseModules(id)
    }).subscribe({
      next: (res) => {
        this.course.set(res.course);
        this.modules.set(res.modules || []);
        this.isLoading.set(false);
      },
      error: (err) => {
        console.error('Ошибка загрузки курса', err);
        this.error.set('Не удалось загрузить информацию о курсе. Попробуйте обновить страницу.');
        this.isLoading.set(false);
      }
    });
  }
}
