export interface Course {
  id: number;
  title: string;
  description: string;
  created_by: number;
  is_published: boolean;
  created_at: string;
}

export interface CourseModule {
  id: number;
  course_id: number;
  title: string;
  order_index: number;
  created_at: string;
}
