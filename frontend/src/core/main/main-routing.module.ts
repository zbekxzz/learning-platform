import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { MainComponent } from './main/main.component';

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
      }
    ]
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class MainRoutingModule {}
