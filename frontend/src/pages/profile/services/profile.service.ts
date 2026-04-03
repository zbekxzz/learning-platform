import { Injectable } from '@angular/core';
import { ApiService } from '../../../shared/services/api.service';
import { Observable, tap, of } from 'rxjs';
import { UserProfile } from '../models/profile.model';
import { Course } from '../../courses/models/course.model';

@Injectable({
  providedIn: 'root'
})
export class ProfileService {
  private cachedProfile: UserProfile | null = null;

  constructor(private api: ApiService) {}

  getProfile(): Observable<UserProfile> {
    if (this.cachedProfile) {
      return of(this.cachedProfile);
    }
    return this.api.get<UserProfile>('/back/user/profile').pipe(
      tap(profile => this.cachedProfile = profile)
    );
  }

  getCurrentProfile(): UserProfile | null {
    return this.cachedProfile;
  }

  clearCache(): void {
    this.cachedProfile = null;
  }

  getMyCourses(): Observable<Course[]> {
    return this.api.get<Course[]>('/back/enrollments/my');
  }
}
