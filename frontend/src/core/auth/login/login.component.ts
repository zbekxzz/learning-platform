import { Component } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { AuthService } from '../auth.service';
import { catchError, first, throwError } from 'rxjs';
import { LoginResponse } from '../auth.model';
import { LocalStorageService } from '../../../shared/services/local-storage.service';
import { CustomMessageService } from '../../../shared/services/custom-message.service';
import { UserService } from '../../../shared/services/user.service';

@Component({
  selector: 'app-login',
  imports: [ReactiveFormsModule],
  templateUrl: './login.component.html'
})
export class LoginComponent {
  loginForm: FormGroup;

  constructor(
    private fb: FormBuilder,
    private router: Router,
    private authService: AuthService,
    private userService: UserService,
    private localStorageService: LocalStorageService,
    private messageService: CustomMessageService
  ) {
    this.loginForm = this.fb.group({
      username: ['', [Validators.required]],
      password: ['', [Validators.required, Validators.minLength(5)]]
    });
  }

  onSubmit() {
    if (this.loginForm.invalid) {
      this.messageService.showErrorMessage('Заполните все поля!');
      return;
    }

    const username = this.loginForm.get('username')?.value;
    const password = this.loginForm.get('password')?.value;

    this.authService
      .login(username, password)
      .pipe(
        first(),
        catchError((err) => {
          this.messageService.showErrorMessage(
            'Заполненные данные не корректны! Пожалуйста, попробуйте еще раз!'
          );
          return throwError(() => err);
        })
      )
      .subscribe((data: LoginResponse) => {
        this.userService.setUserData(data);
        this.localStorageService.setToken(data.token);
        this.messageService.showSuccessMessage('Вы вошли в систему!');
        this.router.navigate(['main/']);
      });
  }

  goToRegister() {
    this.router.navigate(['/auth/register']);
  }
}
