import { Injectable } from '@angular/core';
import { Observable, tap } from 'rxjs';
import { ApiService } from '../../shared/services/api.service';
import { LoginResponse, RegisterRequest } from './auth.model';
import { LocalStorageService } from '../../shared/services/local-storage.service';

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  constructor(
    private api: ApiService,
    private localStorage: LocalStorageService
  ) {}

  login(email: string, password: string): Observable<LoginResponse> {
    const payload = { email, password };

    return this.api.post<LoginResponse>('/back/auth/login', payload).pipe(
      tap(res => {
        this.localStorage.setToken(res.token);
      })
    );
  }

  register(data: RegisterRequest): Observable<any> {
    return this.api.post('/back/auth/register', data);
  }

  requestRestoreCode(email: string): Observable<any> {
    return this.api.post('/back/auth/restore/request', { email });
  }

  verifyRestoreCode(email: string, code: string): Observable<any> {
    return this.api.post('/back/auth/restore/verify', { email, code });
  }

  resetPassword(email: string, code: string, password: string, confirm_password: string): Observable<any> {
    return this.api.post('/back/auth/restore/reset', { email, code, password, confirm_password });
  }

  logout(): void {
    this.localStorage.removeTokenAndRedirect();
  }
}

