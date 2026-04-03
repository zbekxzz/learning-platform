import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { MainComponent } from './main/main.component';
import { RoleGuard } from '../guards/role.guard';

const routes: Routes = [
  {
    path: '',
    component: MainComponent,
    children: [
      {
        path: 'courses',
        loadChildren: () => import('../../pages/courses/courses.module').then(m => m.CoursesModule)
      },
      {
        path: 'profile',
        loadComponent: () => import('../../pages/profile/ui/profile-page/profile-page.component').then(c => c.ProfilePageComponent)
      },
      {
        path: 'teacher/courses',
        canActivate: [RoleGuard],
        data: { roles: ['teacher'] },
        loadComponent: () => import('../../pages/teacher/ui/course-management/course-management.component').then(c => c.TeacherCourseManagementComponent)
      },
      {
        path: 'admin/courses',
        canActivate: [RoleGuard],
        data: { roles: ['admin'] },
        loadComponent: () => import('../../pages/admin/ui/course-management/admin-course-management.component').then(c => c.AdminCourseManagementComponent)
      },
      {
        path: 'admin/users',
        canActivate: [RoleGuard],
        data: { roles: ['admin'] },
        loadComponent: () => import('../../pages/admin/ui/user-management/user-management.component').then(c => c.AdminUserManagementComponent)
      }
    ]
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class MainRoutingModule {}
