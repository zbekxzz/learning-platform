import { Component, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ProfileService } from '../../services/profile.service';
import { UserProfile } from '../../models/profile.model';
import { Course } from '../../../courses/models/course.model';
import { RouterLink } from '@angular/router';
import { CoursesService } from '../../../courses/services/courses.service';

@Component({
  selector: 'app-profile-page',
  standalone: true,
  imports: [CommonModule, RouterLink],
  templateUrl: './profile-page.component.html',
  styleUrl: './profile-page.component.scss'
})
export class ProfilePageComponent implements OnInit {
  profile = signal<UserProfile | null>(null);
  isLoading = signal<boolean>(true);
  error = signal<string | null>(null);

  myCourses = signal<Course[]>([]);
  isLoadingCourses = signal<boolean>(true);
  coursesError = signal<string | null>(null);

  certificates = signal<any[]>([]);
  isLoadingCertificates = signal<boolean>(true);
  certificatesError = signal<string | null>(null);

  constructor(
    private profileService: ProfileService,
    private coursesService: CoursesService
  ) {}

  ngOnInit() {
    this.loadProfile();
    this.loadMyCourses();
    this.loadCertificates();
  }

  loadProfile() {
    this.isLoading.set(true);
    this.error.set(null);

    this.profileService.getProfile().subscribe({
      next: (data) => {
        this.profile.set(data);
        this.isLoading.set(false);
      },
      error: (err) => {
        console.error('Error fetching profile:', err);
        this.error.set('Профиль деректерін жүктеу мүмкін болмады.');
        this.isLoading.set(false);
      }
    });
  }

  loadMyCourses() {
    this.isLoadingCourses.set(true);
    this.coursesError.set(null);

    this.profileService.getMyCourses().subscribe({
      next: (data) => {
        this.myCourses.set(data || []);
        this.isLoadingCourses.set(false);
      },
      error: (err) => {
        console.error('Error fetching my courses:', err);
        this.coursesError.set('Курстарыңызды жүктеу мүмкін болмады.');
        this.isLoadingCourses.set(false);
      }
    });
  }

  loadCertificates() {
    this.isLoadingCertificates.set(true);
    this.certificatesError.set(null);

    this.coursesService.getMyCertificates().subscribe({
      next: (data) => {
        this.certificates.set(data || []);
        this.isLoadingCertificates.set(false);
      },
      error: (err) => {
        console.error('Error fetching certificates:', err);
        this.certificatesError.set('Сертификаттарыңызды жүктеу мүмкін болмады.');
        this.isLoadingCertificates.set(false);
      }
    });
  }
}
