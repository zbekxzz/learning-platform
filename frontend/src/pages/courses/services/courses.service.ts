import { Injectable } from '@angular/core';
import { ApiService } from '../../../shared/services/api.service';
import { Observable } from 'rxjs';
import { Course } from '../models/course.model';

@Injectable({
  providedIn: 'root'
})
export class CoursesService {

  constructor(private api: ApiService) {}

  getCourses(page: number = 1, limit: number = 10) {
    return this.api.get('/back/courses', { page, limit });
  }

  getCourse(id: number): Observable<Course> {
    return this.api.get(`/back/courses/${id}`);
  }

  enroll(courseId: number) {
    return this.api.post(`/back/enrollments/${courseId}`, {});
  }

}
