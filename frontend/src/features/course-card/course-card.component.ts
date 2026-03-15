import {Component, Input} from '@angular/core';
import {Course} from '../../pages/courses/models/course.model';
import {CoursesService} from '../../pages/courses/services/courses.service';

@Component({
  selector: 'app-course-card',
  imports: [],
  templateUrl: './course-card.component.html',
})
export class CourseCardComponent {
  @Input() course!: Course;

  constructor(private coursesService: CoursesService) {}

  enroll() {
    this.coursesService.enroll(this.course.id).subscribe();
  }
}
