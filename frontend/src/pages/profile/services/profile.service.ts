import { Injectable } from '@angular/core';
import { ApiService } from '../../../shared/services/api.service';
import { Observable } from 'rxjs';
import { UserProfile } from '../models/profile.model';
import { Course } from '../../courses/models/course.model';

@Injectable({
  providedIn: 'root'
})
export class ProfileService {
  constructor(private api: ApiService) {}

  getProfile(): Observable<UserProfile> {
    return this.api.get<UserProfile>('/back/user/profile');
  }

  getMyCourses(): Observable<Course[]> {
    return this.api.get<Course[]>('/back/enrollments/my');
  }
}
