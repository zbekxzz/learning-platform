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

  logout(): void {
    this.localStorage.removeTokenAndRedirect();
  }
}
