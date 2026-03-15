import { Routes } from '@angular/router';
import { AuthGuard } from '../core/auth/auth.guard';

export const routes: Routes = [
  {
    path: '',
    pathMatch: 'full',
    redirectTo: 'auth'
  },
  {
    path: 'auth',
    loadChildren: () => import('../core/auth/auth.module').then((m) => m.AuthModule)
  },
  {
    path: 'main',
    loadChildren: () => import('../core/main/main.module').then((m) => m.MainModule)
  },
  {
    path: '**',
    pathMatch: 'full',
    redirectTo: 'auth'
  }
];
