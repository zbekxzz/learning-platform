import { Injectable } from '@angular/core';
import { CanActivate, ActivatedRouteSnapshot, RouterStateSnapshot, Router } from '@angular/router';
import { Observable, map } from 'rxjs';
import { ProfileService } from '../../pages/profile/services/profile.service';
import { UserRole } from '../../pages/profile/models/profile.model';

@Injectable({
  providedIn: 'root'
})
export class RoleGuard implements CanActivate {

  constructor(
    private profileService: ProfileService,
    private router: Router
  ) {}

  canActivate(
    route: ActivatedRouteSnapshot,
    state: RouterStateSnapshot
  ): Observable<boolean> {
    
    // Получаем список ролей, которым разрешен доступ к этому роуту
    const allowedRoles = route.data['roles'] as UserRole[];

    return this.profileService.getProfile().pipe(
      map(profile => {
        if (!profile || !allowedRoles.includes(profile.role)) {
          // Если роли не совпадают, отправляем на базовую страницу
          this.router.navigate(['/main/courses']);
          return false;
        }
        return true;
      })
    );
  }
}
