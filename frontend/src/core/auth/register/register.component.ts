import { Component } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { catchError, first, throwError } from 'rxjs';

import { AuthService } from '../auth.service';
import { CustomMessageService } from '../../../shared/services/custom-message.service';

@Component({
  selector: 'app-register',
  imports: [ReactiveFormsModule],
  templateUrl: './register.component.html'
})
export class RegisterComponent {

  registerForm: FormGroup;

  constructor(
    private fb: FormBuilder,
    private router: Router,
    private authService: AuthService,
    private messageService: CustomMessageService
  ) {

    this.registerForm = this.fb.group({
      full_name: ['', [Validators.required]],
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.required, Validators.minLength(5)]],
      confirm_password: ['', [Validators.required]]
    });

  }

  onSubmit() {

    if (this.registerForm.invalid) {
      this.messageService.showErrorMessage('Заполните все поля!');
      return;
    }

    const { full_name, email, password, confirm_password } = this.registerForm.value;

    if (password !== confirm_password) {
      this.messageService.showErrorMessage('Пароли не совпадают!');
      return;
    }

    this.authService
      .register({
        full_name,
        email,
        password
      })
      .pipe(
        first(),
        catchError((err) => {
          this.messageService.showErrorMessage(
            'Ошибка регистрации. Попробуйте ещё раз.'
          );
          return throwError(() => err);
        })
      )
      .subscribe(() => {

        this.messageService.showSuccessMessage('Регистрация успешна!');

        this.router.navigate(['/auth/login']);

      });

  }

  goToLogin() {
    this.router.navigate(['/auth/login']);
  }

}
