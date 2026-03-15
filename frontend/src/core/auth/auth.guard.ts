import { Injectable } from '@angular/core';
import { CanActivate, Router } from '@angular/router';
import { LocalStorageService } from '../../shared/services/local-storage.service';

@Injectable({
  providedIn: 'root'
})
export class AuthGuard implements CanActivate {

  constructor(
    private storage: LocalStorageService,
    private router: Router
  ) {}

  canActivate(): boolean {

    const token = this.storage.getToken();

    if (!token) {
      this.router.navigate(['/auth/login']);
      return false;
    }

    return true;
  }
}
