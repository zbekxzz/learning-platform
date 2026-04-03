import { Component, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { CoursesService } from '../../services/courses.service';
import { Course, CourseModule } from '../../models/course.model';
import { forkJoin } from 'rxjs';
import { ProfileService } from '../../../profile/services/profile.service';

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

  isEnrolling = signal<boolean>(false);
  enrollSuccess = signal<boolean>(false);
  enrollError = signal<string | null>(null);

  constructor(
    private route: ActivatedRoute,
    private coursesService: CoursesService,
    private profileService: ProfileService
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
      modules: this.coursesService.getCourseModules(id),
      myCourses: this.profileService.getMyCourses()
    }).subscribe({
      next: (res) => {
        this.course.set(res.course);
        this.modules.set(res.modules || []);
        
        // Проверяем, записан ли уже пользователь на этот курс
        const alreadyEnrolled = (res.myCourses || []).some((c: Course) => c.id === id);
        this.enrollSuccess.set(alreadyEnrolled);
        
        this.isLoading.set(false);
      },
      error: (err) => {
        console.error('Ошибка загрузки курса', err);
        this.error.set('Не удалось загрузить информацию о курсе. Попробуйте обновить страницу.');
        this.isLoading.set(false);
      }
    });
  }

  enrollInCourse() {
    const currentCourse = this.course();
    if (!currentCourse) return;

    this.isEnrolling.set(true);
    this.enrollError.set(null);

    this.coursesService.enroll(currentCourse.id).subscribe({
      next: () => {
        this.enrollSuccess.set(true);
        this.isEnrolling.set(false);
      },
      error: (err) => {
        console.error('Ошибка записи на курс', err);
        this.enrollError.set('Не удалось записаться на курс. Возможно, вы уже записаны.');
        this.isEnrolling.set(false);
      }
    });
  }
}
