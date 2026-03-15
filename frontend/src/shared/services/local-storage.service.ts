import { Injectable } from '@angular/core';
import { Router } from '@angular/router';

@Injectable({
  providedIn: 'root'
})
export class LocalStorageService {
  private readonly TOKEN_KEY: string = 'Token';

  constructor(private router: Router) {}

  setToken(token: string): void {
    localStorage.setItem(this.TOKEN_KEY, token);
  }

  getToken(): string | null {
    return localStorage.getItem(this.TOKEN_KEY);
  }

  removeTokenAndRedirect(): void {
    if (this.getToken()) {
      localStorage.removeItem(this.TOKEN_KEY);
    }

    this.router.navigate(['/auth']);
  }
}
