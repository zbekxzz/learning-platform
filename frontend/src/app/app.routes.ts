import { Routes } from '@angular/router';
import { AuthGuard } from '../core/auth/auth.guard';
import { NoAuthGuard } from '../core/auth/no-auth.guard';

export const routes: Routes = [
  {
    path: '',
    pathMatch: 'full',
    redirectTo: 'main'
  },
  {
    path: 'auth',
    canActivate: [NoAuthGuard],
    loadChildren: () => import('../core/auth/auth.module').then((m) => m.AuthModule)
  },
  {
    path: 'main',
    canActivate: [AuthGuard],
    loadChildren: () => import('../core/main/main.module').then((m) => m.MainModule)
  },
  {
    path: '**',
    pathMatch: 'full',
    redirectTo: 'main'
  }
];
