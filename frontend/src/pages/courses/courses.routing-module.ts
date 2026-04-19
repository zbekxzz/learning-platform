import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import {CoursesListComponent} from './ui/courses-list/courses-list.component';
import {CourseDetailsComponent} from './ui/course-details/course-details.component';
import {CourseLearningComponent} from './ui/course-learning/course-learning.component';

const routes: Routes = [
  {
    path: '',
    pathMatch: 'full',
    redirectTo: 'list'
  },
  {
    path: 'list',
    component: CoursesListComponent
  },
  {
    path: 'course/:id',
    component: CourseDetailsComponent
  },
  {
    path: 'course/:id/learn',
    component: CourseLearningComponent
  },
  {
    path: 'course/:id/learn/:moduleId',
    component: CourseLearningComponent
  }
]

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class CoursesRoutingModule {}
