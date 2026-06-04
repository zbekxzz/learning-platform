import { Component } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { catchError, first, throwError } from 'rxjs';
import { AuthService } from '../auth.service';
import { CustomMessageService } from '../../../shared/services/custom-message.service';

@Component({
  selector: 'app-restore',
  standalone: true,
  imports: [ReactiveFormsModule],
  templateUrl: './restore.component.html'
})
export class RestoreComponent {
  restoreForm: FormGroup;
  step = 1;

  constructor(
    private fb: FormBuilder,
    private router: Router,
    private authService: AuthService,
    private messageService: CustomMessageService
  ) {
    this.restoreForm = this.fb.group({
      email: ['', [Validators.required, Validators.email]],
      code: ['', []],
      password: ['', []],
      confirm_password: ['', []]
    });
  }

  onSubmit() {
    const email = this.restoreForm.get('email')?.value;

    if (this.step === 1) {
      const emailControl = this.restoreForm.get('email');
      if (!emailControl || emailControl.invalid) {
        this.messageService.showErrorMessage('Введите корректный email!');
        return;
      }

      this.authService
        .requestRestoreCode(email)
        .pipe(
          first(),
          catchError((err) => {
            const errMsg = err.error?.error || 'Произошла ошибка при отправке кода.';
            this.messageService.showErrorMessage(errMsg);
            return throwError(() => err);
          })
        )
        .subscribe(() => {
          this.messageService.showSuccessMessage('Код подтверждения отправлен на вашу почту.');
          this.step = 2;
        });

    } else if (this.step === 2) {
      const code = this.restoreForm.get('code')?.value;
      if (!code || code.trim() === '') {
        this.messageService.showErrorMessage('Введите код подтверждения!');
        return;
      }

      this.authService
        .verifyRestoreCode(email, code)
        .pipe(
          first(),
          catchError((err) => {
            const errMsg = err.error?.error || 'Неверный код или время действия истекло!';
            this.messageService.showErrorMessage(errMsg + ' Возврат на отправку кода.');
            // Reset to step 1
            this.step = 1;
            this.restoreForm.get('code')?.setValue('');
            return throwError(() => err);
          })
        )
        .subscribe(() => {
          this.messageService.showSuccessMessage('Код успешно подтвержден!');
          this.step = 3;
        });

    } else if (this.step === 3) {
      const code = this.restoreForm.get('code')?.value;
      const password = this.restoreForm.get('password')?.value;
      const confirmPassword = this.restoreForm.get('confirm_password')?.value;

      if (!password || password.length < 5) {
        this.messageService.showErrorMessage('Пароль должен быть не менее 5 символов!');
        return;
      }

      if (password !== confirmPassword) {
        this.messageService.showErrorMessage('Пароли не совпадают!');
        return;
      }

      this.authService
        .resetPassword(email, code, password, confirmPassword)
        .pipe(
          first(),
          catchError((err) => {
            const errMsg = err.error?.error || 'Произошла ошибка при сбросе пароля.';
            this.messageService.showErrorMessage(errMsg);
            return throwError(() => err);
          })
        )
        .subscribe(() => {
          this.messageService.showSuccessMessage('Пароль успешно изменен!');
          this.router.navigate(['/auth/login']);
        });
    }
  }

  goToLogin() {
    this.router.navigate(['/auth/login']);
  }
}
