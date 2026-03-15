import { Component } from '@angular/core';
import {CoursesService} from '../../services/courses.service';
import {Course} from '../../models/course.model';
import {CourseCardComponent} from '../../../../features/course-card/course-card.component';

@Component({
  selector: 'app-courses-list',
  imports: [
    CourseCardComponent
  ],
  templateUrl: './courses-list.component.html',
})
export class CoursesListComponent {
  courses: Course[] = [];

  page = 1;
  total = 0;
  limit = 10;

  constructor(private coursesService: CoursesService) {}

  ngOnInit() {
    this.loadCourses();
  }

  loadCourses() {

    this.coursesService.getCourses(this.page, this.limit)
      .subscribe((res: any) => {

        this.courses = res.data;
        this.total = res.total;

      });

  }
}
