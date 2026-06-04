import { Component, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { FormBuilder, FormGroup, FormArray, ReactiveFormsModule, Validators } from '@angular/forms';
import { CdkDragDrop, DragDropModule, moveItemInArray } from '@angular/cdk/drag-drop';
import { CoursesService } from '../../../courses/services/courses.service';
import { ModuleMaterial } from '../../../courses/models/course.model';
import { TestService } from '../../../courses/services/test.service';

@Component({
  selector: 'app-module-editor',
  standalone: true,
  imports: [CommonModule, RouterLink, ReactiveFormsModule, DragDropModule],
  templateUrl: './module-editor.component.html',
  styleUrl: './module-editor.component.scss'
})
export class ModuleEditorComponent implements OnInit {
  courseId!: number;
  moduleId!: number;

  materials = signal<ModuleMaterial[]>([]);
  activeMaterial = signal<ModuleMaterial | null>(null);

  isLoading = signal<boolean>(true);
  error = signal<string | null>(null);

  isMaterialModalOpen = signal<boolean>(false);
  materialForm: FormGroup;
  isSaving = signal<boolean>(false);
  isUploading = signal<boolean>(false);
  selectedFile: File | null = null;

  // Test Builder States
  moduleTest = signal<any | null>(null);
  isTestModalOpen = signal<boolean>(false);
  testForm!: FormGroup;
  editingTestId = signal<number | null>(null);

  constructor(
    private route: ActivatedRoute,
    private coursesService: CoursesService,
    private testService: TestService,
    private fb: FormBuilder
  ) {
    this.materialForm = this.fb.group({
      title: ['', Validators.required],
      type: ['text', Validators.required],
      content: [''],
      file_url: [''],
      external_url: [''],
      tabs: this.fb.array([])
    });
    this.initTestForm();
  }

  initTestForm() {
    this.testForm = this.fb.group({
      title: ['', Validators.required],
      timeLimit: [30, [Validators.required, Validators.min(1)]],
      maxAttempts: [2, [Validators.required, Validators.min(1)]],
      questions: this.fb.array([])
    });
  }


  get tabs(): FormArray {
    return this.materialForm.get('tabs') as FormArray;
  }

  addTab() {
    this.tabs.push(this.fb.group({
      title: ['', Validators.required],
      body: ['', Validators.required]
    }));
  }

  removeTab(index: number) {
    this.tabs.removeAt(index);
  }

  onFileSelected(event: Event) {
    const input = event.target as HTMLInputElement;
    if (input.files && input.files.length > 0) {
      this.selectedFile = input.files[0];
      this.uploadSelectedFile();
    }
  }

  uploadSelectedFile() {
    if (!this.selectedFile) return;
    this.isUploading.set(true);
    this.coursesService.uploadFile(this.selectedFile).subscribe({
      next: (res) => {
        this.materialForm.patchValue({ file_url: res.url });
        this.isUploading.set(false);
        alert('Файл сәтті жүктелді!');
      },
      error: (err) => {
        console.error('File upload failed', err);
        this.isUploading.set(false);
        alert('Файлды жүктеу сәтсіз аяқталды');
      }
    });
  }


  ngOnInit() {
    this.courseId = parseInt(this.route.snapshot.paramMap.get('id') || '0', 10);
    this.moduleId = parseInt(this.route.snapshot.paramMap.get('moduleId') || '0', 10);
    
    if (this.moduleId) {
      this.loadMaterials();
      this.loadModuleTest();
    }
  }

  loadModuleTest() {
    this.testService.startTest(this.moduleId).subscribe({
      next: (res) => {
        if (res && res.test) {
          this.moduleTest.set(res.test);
        } else {
          this.moduleTest.set(null);
        }
      },
      error: () => {
        this.moduleTest.set(null);
      }
    });
  }


  loadMaterials() {
    this.isLoading.set(true);
    this.coursesService.getModuleMaterials(this.moduleId).subscribe({
      next: (mats) => {
        this.materials.set(mats || []);
        if (mats && mats.length > 0 && !this.activeMaterial()) {
          this.activeMaterial.set(mats[0]);
        }
        this.isLoading.set(false);
      },
      error: () => {
        this.error.set('Материалдарды жүктеу мүмкін болмады.');
        this.isLoading.set(false);
      }
    });
  }

  selectMaterial(mat: ModuleMaterial) {
    this.activeMaterial.set(mat);
  }

  drop(event: CdkDragDrop<ModuleMaterial[]>) {
    const currentMaterials = [...this.materials()];
    moveItemInArray(currentMaterials, event.previousIndex, event.currentIndex);
    
    // Optimistic UI update
    this.materials.set(currentMaterials);

    // Prepare batch update payload
    const updates = currentMaterials.map((mat, index) => ({
      id: mat.id,
      order_index: index + 1
    }));

    // Send to backend
    this.coursesService.reorderMaterials(this.moduleId, updates).subscribe({
      next: () => {
        console.log('Materials reordered successfully');
      },
      error: (err) => {
        console.error('Failed to reorder materials', err);
        // Revert UI on failure
        this.loadMaterials();
      }
    });
  }

  // --- Add Material Flow ---
  openMaterialModal() {
    this.materialForm.reset({ type: 'text' });
    while (this.tabs.length) {
      this.tabs.removeAt(0);
    }
    this.selectedFile = null;
    this.isMaterialModalOpen.set(true);
  }

  closeMaterialModal() {
    this.isMaterialModalOpen.set(false);
  }

  saveMaterial() {
    if (this.materialForm.invalid) {
      this.materialForm.markAllAsTouched();
      return;
    }

    this.isSaving.set(true);
    const orderIndex = this.materials().length + 1;
    const { title, type, content, file_url, external_url, tabs } = this.materialForm.value;

    let finalContent = content || '';
    if (type === 'interactive') {
      finalContent = JSON.stringify(tabs || []);
    }

    this.coursesService.addMaterial({
      module_id: this.moduleId,
      title,
      type,
      content: finalContent,
      file_url: file_url || undefined,
      external_url: external_url || undefined,
      order_index: orderIndex
    }).subscribe({
      next: (newMat) => {
        this.isSaving.set(false);
        this.closeMaterialModal();
        this.loadMaterials();
        this.activeMaterial.set(newMat);
      },
      error: () => {
        this.isSaving.set(false);
        alert('Қате пайда болды');
      }
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
      module_id: this.moduleId,
      type: 'module',
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
        this.loadModuleTest();
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
        this.loadModuleTest();
      },
      error: () => {
        this.isSaving.set(false);
        alert('Тестті өшіру сәтсіз аяқталды.');
      }
    });
  }
}

