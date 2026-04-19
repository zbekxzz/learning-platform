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

export interface CourseChapter {
  id: number;
  course_id: number;
  title: string;
  order_index: number;
  created_at: string;
}

export interface CourseStructureItem {
  chapter: CourseChapter;
  modules: CourseModule[] | null;
}

export type ModuleStatus = 'locked' | 'available' | 'completed';

export interface ModuleMaterial {
  id: number;
  module_id: number;
  title: string;
  type: string; // 'text', 'video', 'pdf', 'interactive', etc.
  file_url: string | null;
  external_url: string | null;
  content: any; 
  order_index: number;
  created_at: string;
}

