import { Component, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { AdminUsersService } from '../../services/admin-users.service';
import { UserProfile, UserRole } from '../../../profile/models/profile.model';

@Component({
  selector: 'app-admin-user-management',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './user-management.component.html',
  styleUrl: './user-management.component.scss'
})
export class AdminUserManagementComponent implements OnInit {
  users = signal<UserProfile[]>([]);
  isLoading = signal<boolean>(true);
  error = signal<string | null>(null);

  isFormOpen = signal<boolean>(false);
  isSaving = signal<boolean>(false);
  editingUserId = signal<number | null>(null);

  userForm: FormGroup;

  roles: UserRole[] = ['student', 'teacher', 'admin'];

  constructor(
    private usersService: AdminUsersService,
    private fb: FormBuilder
  ) {
    this.userForm = this.fb.group({
      full_name: ['', Validators.required],
      email: ['', [Validators.required, Validators.email]],
      role: ['student', Validators.required],
      password: ['']
    });
  }

  ngOnInit() {
    this.loadUsers();
  }

  loadUsers() {
    this.isLoading.set(true);
    this.error.set(null);
    this.usersService.getUsers().subscribe({
      next: (res) => {
        // Защита: вдруг бек возвращает { data: [...] }
        const usersArray = Array.isArray(res) ? res : (res.data || []);
        this.users.set(usersArray);
        this.isLoading.set(false);
      },
      error: (err) => {
        console.error('Failed to load users', err);
        this.error.set('Не удалось загрузить список пользователей.');
        this.isLoading.set(false);
      }
    });
  }

  openCreateForm() {
    this.editingUserId.set(null);
    this.userForm.reset({ role: 'student' });
    this.userForm.get('password')?.setValidators([Validators.required]);
    this.userForm.get('password')?.updateValueAndValidity();
    this.isFormOpen.set(true);
  }

  openEditForm(user: UserProfile) {
    this.editingUserId.set(user.id);
    this.userForm.patchValue({
      full_name: user.full_name,
      email: user.email,
      role: user.role,
      password: ''
    });
    this.userForm.get('password')?.clearValidators();
    this.userForm.get('password')?.updateValueAndValidity();
    this.isFormOpen.set(true);
  }

  closeForm() {
    this.isFormOpen.set(false);
  }

  saveUser() {
    if (this.userForm.invalid) {
      this.userForm.markAllAsTouched();
      return;
    }

    this.isSaving.set(true);
    const formValue = this.userForm.value;
    
    // Удаляем password, если он пустой при редактировании
    const payload = { ...formValue };
    if (!payload.password) {
      delete payload.password;
    }

    const currentId = this.editingUserId();

    const request$ = currentId 
      ? this.usersService.updateUser(currentId, payload)
      : this.usersService.createUser(payload);

    request$.subscribe({
      next: () => {
        this.isSaving.set(false);
        this.closeForm();
        this.loadUsers();
      },
      error: (err) => {
        console.error('Failed to save user', err);
        // Тут можно было бы показать локальную ошибку формы
        this.isSaving.set(false);
        alert('Ошибка при сохранении пользователя');
      }
    });
  }

  deleteUser(id: number) {
    if (!confirm('Вы уверены, что хотите удалить этого пользователя?')) {
      return;
    }

    this.usersService.deleteUser(id).subscribe({
      next: () => {
        this.loadUsers(); // Перезагружаем список
      },
      error: (err) => {
        console.error('Failed to delete user', err);
        alert('Не удалось удалить пользователя.');
      }
    });
  }
}
