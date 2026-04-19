import { Component, OnInit, signal, computed, effect } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { DomSanitizer, SafeResourceUrl } from '@angular/platform-browser';
import { FormsModule } from '@angular/forms';
import { CoursesService } from '../../services/courses.service';
import { Course, CourseModule, CourseStructureItem, ModuleMaterial, ModuleStatus } from '../../models/course.model';
import { forkJoin } from 'rxjs';
import { TestService } from '../../services/test.service';
import { TestItem, TestQuestion, TestAnswerOption } from '../../models/test.model';

@Component({
  selector: 'app-course-learning',
  standalone: true,
  imports: [CommonModule, RouterLink, FormsModule],
  templateUrl: './course-learning.component.html',
  styleUrl: './course-learning.component.scss'
})
export class CourseLearningComponent implements OnInit {
  courseId = signal<number>(0);
  moduleId = signal<number>(0);

  course = signal<Course | null>(null);
  structure = signal<CourseStructureItem[]>([]);
  isLoading = signal<boolean>(true);

  // Accordion UI State
  expandedChapters = signal<Set<number>>(new Set());

  // Progress & Lesson Status
  moduleStatuses = signal<Record<number, ModuleStatus>>({}); // map of moduleId -> status
  courseProgress = computed(() => {
    const statuses = this.moduleStatuses();
    const total = Object.keys(statuses).length;
    if (total === 0) return 0;
    const completed = Object.values(statuses).filter(s => s === 'completed').length;
    return Math.round((completed / total) * 100);
  });

  // Active Module Content
  materials = signal<ModuleMaterial[]>([]);
  
  // Test/Quiz State
  testData = signal<TestItem | null>(null);
  testQuestions = signal<TestQuestion[]>([]);
  testAnswers = signal<Record<number, TestAnswerOption[]>>({});
  
  // User Answers Payload
  userAnswers = signal<Record<number, any>>({});

  // UI Interactive State
  activeTabIndexes = signal<Record<number, number>>({});
  quizSubmitted = signal<boolean>(false);
  quizScore = signal<number | null>(null);
  quizError = signal<string | null>(null);

  // Derived properties for template convenience
  currentChapterTitle = signal<string>('');
  currentModuleTitle = signal<string>('');

  // Navigation
  prevModuleId = signal<number | null>(null);
  nextModuleId = signal<number | null>(null);

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private coursesService: CoursesService,
    private sanitizer: DomSanitizer,
    private testService: TestService
  ) {
    // Listen to route params changes without reloading the component
    this.route.paramMap.subscribe(params => {
      const cId = Number(params.get('id'));
      const mId = Number(params.get('moduleId'));
      if (cId) this.courseId.set(cId);
      if (mId) {
        this.moduleId.set(mId);
        this.loadModuleContent(mId);
        // Only update pointers if structure is already loaded
        if (this.structure().length > 0) {
          this.updateNavPointers(mId);
        }
      }
    });

  }

  ngOnInit() {
    this.loadCourseData();
  }

  loadCourseData() {
    this.isLoading.set(true);
    const id = this.courseId();

    forkJoin({
      course: this.coursesService.getCourse(id),
      structure: this.coursesService.getCourseStructure(id)
    }).subscribe({
      next: (res) => {
        this.course.set(res.course);
        this.structure.set(res.structure || []);

        // Mock statuses initialization
        // We'll set the first module to 'available', the rest to 'locked' temporarily
        // Wait, for demo purposes, if user has this course, let's make the first 2 available, and some locked.
        const st: Record<number, ModuleStatus> = {};
        let isFirst = true;

        (res.structure || []).forEach(item => {
          (item.modules || []).forEach((mod: CourseModule) => {
            // If first module, make available. Later ones make locked.
            // But if user clicks around, we don't want to enforce locking strictly in frontend, 
            // although we should to demo it. Let's make all of Chapter 1 available.
            if (item.chapter.order_index === 1) {
              st[mod.id] = (mod.order_index <= 2) ? 'completed' : 'available';
            } else {
              st[mod.id] = 'locked';
            }
          });
        });

        this.moduleStatuses.set(st);

        // Handle redirect if no moduleId is provided
        if (!this.moduleId()) {
          let firstModId = 0;
          for (const item of (res.structure || [])) {
            if (item.modules && item.modules.length > 0) {
              firstModId = item.modules[0].id;
              break;
            }
          }
          if (firstModId) {
            this.isLoading.set(false);
            this.router.navigate(['/main/courses/course', id, 'learn', firstModId], { replaceUrl: true });
            return; // Exit here, let the paramMap subscription handle subsequent setup
          } else {
            // In case the course has no modules whatsoever
            this.isLoading.set(false);
            return;
          }
        }

        // Structure is now ready, update nav pointers for current active module
        this.updateNavPointers(this.moduleId());
        this.isLoading.set(false);
      },
      error: (err) => {
        console.error('Ошибка загрузки', err);
        this.isLoading.set(false);
      }
    });
  }

  // --- Module Content Generation ---
  loadModuleContent(moduleId: number) {
    // Reset test/tabs state when loading new module
    this.quizSubmitted.set(false);
    this.quizScore.set(null);
    this.quizError.set(null);
    this.activeTabIndexes.set({});
    this.userAnswers.set({});
    this.testData.set(null);
    this.testQuestions.set([]);
    this.testAnswers.set({});

    // Fetch real materials
    this.coursesService.getModuleMaterials(moduleId).subscribe({
      next: (materials) => {
        // Sort by order_index just in case
        materials.sort((a, b) => a.order_index - b.order_index);
        this.materials.set(materials);
      },
      error: (err) => {
        console.error('Ошибка загрузки материалов темы', err);
        this.materials.set([]);
      }
    });

    // Fetch test
    this.testService.startTest(moduleId).subscribe({
      next: (res) => {
        if (res && res.test) {
          this.testData.set(res.test);
          
          // Safe fallback: if backend returns empty type, assume it is mcq
          const questions = (res.questions || []).map(q => {
            if (!q.type) {
               q.type = 'mcq';
            }
            return q;
          });
          
          this.testQuestions.set(questions);
          this.testAnswers.set(res.answers || {});
        }
      },
      error: (err) => {
        console.warn('Тест для данного модуля не найден или недоступен.');
      }
    });
  }

  // Find prev/next IDs safely traversing chapters and modules
  updateNavPointers(currentModId: number) {
    const flatModules: { id: number, chapter: string, title: string }[] = [];
    const struct = this.structure();

    struct.forEach(item => {
      // Auto-expand the chapter containing the current module
      if (item.modules?.some((m: CourseModule) => m.id === currentModId)) {
        this.expandedChapters.update(set => {
          const newSet = new Set(set);
          newSet.add(item.chapter.id);
          return newSet;
        });
      }

      (item.modules || []).forEach((mod: CourseModule) => {
        flatModules.push({
          id: mod.id,
          chapter: item.chapter.title,
          title: mod.title
        });
      });
    });

    const currIndex = flatModules.findIndex(m => m.id === currentModId);
    if (currIndex !== -1) {
      this.currentChapterTitle.set(flatModules[currIndex].chapter);
      this.currentModuleTitle.set(flatModules[currIndex].title);
      this.prevModuleId.set(currIndex > 0 ? flatModules[currIndex - 1].id : null);
      this.nextModuleId.set(currIndex < flatModules.length - 1 ? flatModules[currIndex + 1].id : null);
    }
  }

  toggleChapterAccordion(chapterId: number) {
    const current = new Set(this.expandedChapters());
    if (current.has(chapterId)) {
      current.delete(chapterId);
    } else {
      current.add(chapterId);
    }
    this.expandedChapters.set(current);
  }

  navigateToModule(modId: number) {
    const statuses = this.moduleStatuses();
    // Block navigation if locked
    if (statuses[modId] === 'locked') return;

    this.router.navigate(['/main/courses/course', this.courseId(), 'learn', modId]);
  }

  // Dynamic Tabs
  setActiveTab(blockId: number, index: number) {
    this.activeTabIndexes.update(dict => {
      return { ...dict, [blockId]: index };
    });
  }

  // User answers integration
  setUserAnswer(questionId: number, answer: any) {
    this.userAnswers.update(current => ({ ...current, [questionId]: answer }));
  }

  onMatchChange(questionId: number, left: string, event: Event) {
    const target = event.target as HTMLSelectElement;
    const val = target.value;
    
    this.userAnswers.update(current => {
      const updated = { ...current };
      if (!updated[questionId]) {
        updated[questionId] = {};
      }
      updated[questionId][left] = val;
      return updated;
    });
  }

  isSubmitDisabled(): boolean {
    const qCount = this.testQuestions().length;
    return Object.keys(this.userAnswers()).length === 0 && qCount > 0;
  }

  // Quiz submission
  submitTest() {
    const test = this.testData();
    if (!test) return;

    const payload = {
      test_id: test.id,
      answers: this.userAnswers()
    };

    this.testService.submitTest(payload).subscribe({
      next: (res) => {
        this.quizSubmitted.set(true);
        this.quizScore.set(res.score);
        this.quizError.set(null);

        // Update status to completed!
        this.moduleStatuses.update(dict => {
          return { ...dict, [this.moduleId()]: 'completed' };
        });

        // Auto-unlock next module
        const nextId = this.nextModuleId();
        if (nextId) {
          this.moduleStatuses.update(dict => {
            if (dict[nextId] === 'locked') {
              return { ...dict, [nextId]: 'available' };
            }
            return dict;
          });
        }
      },
      error: (err) => {
        this.quizError.set('Ошибка при отправке теста. Проверьте правильность заполнения.');
        this.quizSubmitted.set(false);
      }
    });
  }

  // Security for iframe URL
  safeUrl(url: string | null): SafeResourceUrl {
    return this.sanitizer.bypassSecurityTrustResourceUrl(url || '');
  }
}
