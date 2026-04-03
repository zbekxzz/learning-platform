import { Injectable } from '@angular/core';
import { ApiService } from '../../../shared/services/api.service';
import { Observable } from 'rxjs';
import { Course, CourseModule } from '../models/course.model';

export interface CoursesResponse {
  data: Course[];
  total: number;
}


@Injectable({
  providedIn: 'root'
})
export class CoursesService {

  constructor(private api: ApiService) { }

  getCourses(page: number = 1, limit: number = 10): Observable<CoursesResponse> {
    return this.api.get<CoursesResponse>('/back/courses', { page, limit });
  }

  getCourse(id: number): Observable<Course> {
    return this.api.get<Course>(`/back/courses/${id}`);
  }

  getCourseModules(courseId: number): Observable<CourseModule[]> {
    return this.api.get<CourseModule[]>(`/back/modules/course/${courseId}`);
  }

  enroll(courseId: number) {
    return this.api.post(`/back/enrollments/${courseId}`, {});
  }

}
