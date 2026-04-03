import { Component, OnInit, signal } from '@angular/core';
import { CoursesService } from '../../services/courses.service';
import { Course } from '../../models/course.model';
import { RouterLink } from '@angular/router';
import { DatePipe } from '@angular/common';

@Component({
  selector: 'app-courses-list',
  imports: [
    RouterLink,
    DatePipe
  ],
  templateUrl: './courses-list.component.html',
})
export class CoursesListComponent implements OnInit {
  courses = signal<Course[]>([]);
  isLoading = signal<boolean>(true);
  error = signal<string | null>(null);

  page = 1;
  total = 0;
  limit = 10;

  constructor(private coursesService: CoursesService) { }

  ngOnInit() {
    this.loadCourses();
  }

  loadCourses() {
    this.isLoading.set(true);
    this.error.set(null);

    this.coursesService.getCourses(this.page, this.limit).subscribe({
      next: (res) => {
        this.courses.set(res.data || []);
        this.total = res.total || 0;
        this.isLoading.set(false);
      },
      error: (err) => {
        console.error('Ошибка получения курсов:', err);
        this.error.set('Произошла ошибка при загрузке курсов. Попробуйте еще раз.');
        this.isLoading.set(false);
      }
    });
  }
}
