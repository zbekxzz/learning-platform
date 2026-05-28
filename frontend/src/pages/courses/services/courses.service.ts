import { Injectable } from '@angular/core';
import { ApiService } from '../../../shared/services/api.service';
import { Observable } from 'rxjs';
import { Course, CourseChapter, CourseModule, CourseStructureItem, CourseWithAuthor, ModuleMaterial } from '../models/course.model';

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
    return this.api.get<CoursesResponse>('/back/courses/', { page, limit });
  }

  getCourse(id: number): Observable<Course> {
    return this.api.get<Course>(`/back/courses/${id}`);
  }

  getCourseModules(courseId: number): Observable<CourseModule[]> {
    return this.api.get<CourseModule[]>(`/back/modules/course/${courseId}`);
  }

  getCourseStructure(courseId: number): Observable<CourseStructureItem[]> {
    return this.api.get<CourseStructureItem[]>(`/back/courses/${courseId}/structure`);
  }

  getModuleMaterials(moduleId: number): Observable<ModuleMaterial[]> {
    return this.api.get<ModuleMaterial[]>(`/back/modules/${moduleId}/materials`);
  }

  enroll(courseId: number) {
    return this.api.post(`/back/enrollments/${courseId}`, {});
  }

  // Admin: get all courses with author names
  getAllForAdmin(): Observable<CourseWithAuthor[]> {
    return this.api.get<CourseWithAuthor[]>('/back/courses/admin/all');
  }

  // Teacher: get own courses
  getTeacherCourses(): Observable<Course[]> {
    return this.api.get<Course[]>('/back/courses/teacher/my');
  }

  // Toggle publish/unpublish
  togglePublish(courseId: number): Observable<any> {
    return this.api.put(`/back/courses/${courseId}/publish`, {});
  }

  // Create a new course
  createCourse(data: { title: string; description: string }): Observable<Course> {
    return this.api.post<Course>('/back/courses/', data);
  }

  // Create a new chapter
  addChapter(data: { course_id: number; title: string; order_index: number }): Observable<CourseChapter> {
    return this.api.post<CourseChapter>('/back/chapters/', data);
  }

  // Create a new module
  addModule(data: { chapter_id: number; title: string; order_index: number }): Observable<CourseModule> {
    return this.api.post<CourseModule>('/back/modules/', data);
  }

  // Create a new material
  addMaterial(data: { module_id: number; title: string; type: string; content?: string; order_index: number }): Observable<ModuleMaterial> {
    return this.api.post<ModuleMaterial>('/back/modules/material', data);
  }

  // Reorder materials
  reorderMaterials(moduleId: number, updates: { id: number; order_index: number }[]): Observable<any> {
    return this.api.put(`/back/modules/${moduleId}/materials/reorder`, { updates });
  }

}

