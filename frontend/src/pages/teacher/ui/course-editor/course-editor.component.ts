import { Component, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { FormBuilder, FormGroup, FormArray, ReactiveFormsModule, Validators } from '@angular/forms';
import { CoursesService } from '../../../courses/services/courses.service';
import { CourseStructureItem, Course } from '../../../courses/models/course.model';
import { TestService } from '../../../courses/services/test.service';

@Component({
  selector: 'app-course-editor',
  standalone: true,
  imports: [CommonModule, RouterLink, ReactiveFormsModule],
  templateUrl: './course-editor.component.html',
  styleUrl: './course-editor.component.scss'
})
export class CourseEditorComponent implements OnInit {
  courseId!: number;
  course = signal<Course | null>(null);
  structure = signal<CourseStructureItem[]>([]);
  isLoading = signal<boolean>(true);
  error = signal<string | null>(null);

  // Modals
  isChapterModalOpen = signal<boolean>(false);
  isModuleModalOpen = signal<boolean>(false);
  activeChapterId = signal<number | null>(null);

  chapterForm: FormGroup;
  moduleForm: FormGroup;
  isSaving = signal<boolean>(false);

  expandedChapters = signal<Set<number>>(new Set());

  // Test Builder States
  finalTest = signal<any | null>(null);
  isTestModalOpen = signal<boolean>(false);
  testForm!: FormGroup;
  editingTestId = signal<number | null>(null);

  constructor(
    private route: ActivatedRoute,
    private coursesService: CoursesService,
    private testService: TestService,
    private fb: FormBuilder
  ) {
    this.chapterForm = this.fb.group({
      title: ['', Validators.required]
    });
    this.moduleForm = this.fb.group({
      title: ['', Validators.required]
    });
    this.initTestForm();
  }

  ngOnInit() {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.courseId = parseInt(id, 10);
      this.loadCourseData();
    }
  }

  loadCourseData() {
    this.isLoading.set(true);
    // Load course details
    this.coursesService.getCourse(this.courseId).subscribe({
      next: (c) => this.course.set(c),
      error: () => console.error('Failed to load course')
    });

    // Load structure
    this.coursesService.getCourseStructure(this.courseId).subscribe({
      next: (st) => {
        this.structure.set(st || []);
        // expand all chapters by default
        this.expandedChapters.set(new Set(st.map(s => s.chapter.id)));
        this.isLoading.set(false);
      },
      error: () => {
        this.error.set('Құрылымды жүктеу мүмкін болмады.');
        this.isLoading.set(false);
      }
    });

    this.loadFinalTest();
  }

  loadFinalTest() {
    this.testService.startFinalTest(this.courseId).subscribe({
      next: (res) => {
        if (res && res.test) {
          this.finalTest.set(res.test);
        } else {
          this.finalTest.set(null);
        }
      },
      error: () => {
        this.finalTest.set(null);
      }
    });
  }

  toggleChapter(chapterId: number) {
    const current = new Set(this.expandedChapters());
    if (current.has(chapterId)) {
      current.delete(chapterId);
    } else {
      current.add(chapterId);
    }
    this.expandedChapters.set(current);
  }

  // --- Chapter Flow ---
  openChapterModal() {
    this.chapterForm.reset();
    this.isChapterModalOpen.set(true);
  }

  closeChapterModal() {
    this.isChapterModalOpen.set(false);
  }

  saveChapter() {
    if (this.chapterForm.invalid) {
      this.chapterForm.markAllAsTouched();
      return;
    }
    
    this.isSaving.set(true);
    const orderIndex = this.structure().length + 1;
    const title = this.chapterForm.value.title;

    this.coursesService.addChapter({ course_id: this.courseId, title, order_index: orderIndex }).subscribe({
      next: () => {
        this.isSaving.set(false);
        this.closeChapterModal();
        this.loadCourseData();
      },
      error: () => {
        this.isSaving.set(false);
        alert('Қате пайда болды');
      }
    });
  }

  // --- Module Flow ---
  openModuleModal(chapterId: number) {
    this.activeChapterId.set(chapterId);
    this.moduleForm.reset();
    this.isModuleModalOpen.set(true);
  }

  closeModuleModal() {
    this.isModuleModalOpen.set(false);
    this.activeChapterId.set(null);
  }

  saveModule() {
    if (this.moduleForm.invalid) {
      this.moduleForm.markAllAsTouched();
      return;
    }

    const chapterId = this.activeChapterId();
    if (!chapterId) return;

    this.isSaving.set(true);
    
    // Find chapter to get order index
    const chapterItem = this.structure().find(s => s.chapter.id === chapterId);
    const orderIndex = chapterItem && chapterItem.modules ? chapterItem.modules.length + 1 : 1;
    const title = this.moduleForm.value.title;

    this.coursesService.addModule({ chapter_id: chapterId, title, order_index: orderIndex }).subscribe({
      next: () => {
        this.isSaving.set(false);
        this.closeModuleModal();
        // ensure chapter is expanded
        const expanded = new Set(this.expandedChapters());
        expanded.add(chapterId);
        this.expandedChapters.set(expanded);
        
        this.loadCourseData();
      },
      error: () => {
        this.isSaving.set(false);
        alert('Қате пайда болды');
      }
    });
  }

  initTestForm() {
    this.testForm = this.fb.group({
      title: ['', Validators.required],
      timeLimit: [30, [Validators.required, Validators.min(1)]],
      maxAttempts: [2, [Validators.required, Validators.min(1)]],
      questions: this.fb.array([])
    });
  }

  get questions(): FormArray {
    return this.testForm.get('questions') as FormArray;
  }

  addQuestion() {
    const qGroup = this.fb.group({
      type: ['mcq', Validators.required],
      questionText: ['', Validators.required],
      orderIndex: [this.questions.length + 1],
      correctText: [''],
      answers: this.fb.array([]),
      pairs: this.fb.array([])
    });
    this.questions.push(qGroup);
    this.addAnswer(this.questions.length - 1);
    this.addAnswer(this.questions.length - 1);
  }

  removeQuestion(qIndex: number) {
    this.questions.removeAt(qIndex);
  }

  getAnswers(qIndex: number): FormArray {
    return this.questions.at(qIndex).get('answers') as FormArray;
  }

  addAnswer(qIndex: number) {
    const answers = this.getAnswers(qIndex);
    answers.push(this.fb.group({
      text: ['', Validators.required],
      isCorrect: [false]
    }));
  }

  removeAnswer(qIndex: number, aIndex: number) {
    this.getAnswers(qIndex).removeAt(aIndex);
  }

  getPairs(qIndex: number): FormArray {
    return this.questions.at(qIndex).get('pairs') as FormArray;
  }

  addPair(qIndex: number) {
    const pairs = this.getPairs(qIndex);
    pairs.push(this.fb.group({
      left: ['', Validators.required],
      right: ['', Validators.required]
    }));
  }

  removePair(qIndex: number, pIndex: number) {
    this.getPairs(qIndex).removeAt(pIndex);
  }

  onQuestionTypeChange(qIndex: number) {
    const q = this.questions.at(qIndex) as FormGroup;
    const type = q.get('type')?.value;
    
    // reset arrays
    const answers = q.get('answers') as FormArray;
    while (answers.length) answers.removeAt(0);
    
    const pairs = q.get('pairs') as FormArray;
    while (pairs.length) pairs.removeAt(0);

    q.patchValue({ correctText: '' });

    if (type === 'mcq') {
      this.addAnswer(qIndex);
      this.addAnswer(qIndex);
    } else if (type === 'matching') {
      this.addPair(qIndex);
    }
  }

  openAddTestModal() {
    this.editingTestId.set(null);
    this.initTestForm();
    this.addQuestion();
    this.isTestModalOpen.set(true);
  }

  openEditTestModal(test: any) {
    this.editingTestId.set(test.id);
    this.initTestForm();
    this.isSaving.set(true);

    this.testService.getTestDetails(test.id).subscribe({
      next: (res) => {
        this.isSaving.set(false);
        this.isTestModalOpen.set(true);

        this.testForm.patchValue({
          title: res.test.title,
          timeLimit: res.test.time_limit,
          maxAttempts: res.test.max_attempts
        });

        // Clear existing questions
        while (this.questions.length) this.questions.removeAt(0);

        // Populate questions
        (res.questions || []).forEach((q: any) => {
          const qGroup = this.fb.group({
            type: [q.type, Validators.required],
            questionText: [q.question_text, Validators.required],
            orderIndex: [q.order_index],
            correctText: [''],
            answers: this.fb.array([]),
            pairs: this.fb.array([])
          });

          this.questions.push(qGroup);
          const qIndex = this.questions.length - 1;

          if (q.type === 'mcq') {
            const answers = res.answers[q.id] || [];
            answers.forEach((ans: any) => {
              this.getAnswers(qIndex).push(this.fb.group({
                text: [ans.text, Validators.required],
                isCorrect: [ans.is_correct || false]
              }));
            });
          } else if (q.type === 'matching') {
            const pairs = res.answers[q.id] || [];
            pairs.forEach((p: any) => {
              this.getPairs(qIndex).push(this.fb.group({
                left: [p.left, Validators.required],
                right: [p.right, Validators.required]
              }));
            });
          } else if (q.type === 'open') {
            qGroup.patchValue({ correctText: res.answers[q.id] || '' });
          }
        });
      },
      error: (err) => {
        this.isSaving.set(false);
        alert('Тест мәліметтерін жүктеу мүмкін болмады.');
      }
    });
  }

  closeTestModal() {
    this.isTestModalOpen.set(false);
  }

  saveTest() {
    if (this.testForm.invalid) {
      this.testForm.markAllAsTouched();
      return;
    }

    const { title, timeLimit, maxAttempts, questions } = this.testForm.value;

    this.isSaving.set(true);

    const savePayload = {
      course_id: this.courseId,
      type: 'final',
      title,
      time_limit: timeLimit,
      max_attempts: maxAttempts,
      questions: questions.map((q: any) => ({
        type: q.type,
        question_text: q.questionText,
        order_index: q.orderIndex,
        correct_text: q.type === 'open' ? q.correctText : undefined,
        answers: q.type === 'mcq' ? q.answers.map((a: any) => ({ text: a.text, is_correct: a.isCorrect })) : undefined,
        pairs: q.type === 'matching' ? q.pairs.map((p: any) => ({ left: p.left, right: p.right })) : undefined
      }))
    };

    if (this.editingTestId()) {
      this.testService.deleteTest(this.editingTestId()!).subscribe({
        next: () => {
          this.createTest(savePayload);
        },
        error: () => {
          this.isSaving.set(false);
          alert('Ескі тестті өшіру мүмкін болмады.');
        }
      });
    } else {
      this.createTest(savePayload);
    }
  }

  createTest(payload: any) {
    this.testService.createFullTest(payload).subscribe({
      next: () => {
        this.isSaving.set(false);
        this.isTestModalOpen.set(false);
        this.loadFinalTest();
      },
      error: (err) => {
        this.isSaving.set(false);
        alert('Тестті сақтау кезінде қате пайда болды.');
      }
    });
  }

  deleteTestConfirmation(testId: number) {
    if (!confirm('Тестті өшіргіңіз келетініне сенімдісіз бе?')) return;
    this.isSaving.set(true);
    this.testService.deleteTest(testId).subscribe({
      next: () => {
        this.isSaving.set(false);
        this.loadFinalTest();
      },
      error: () => {
        this.isSaving.set(false);
        alert('Тестті өшіру сәтсіз аяқталды.');
      }
    });
  }
}
