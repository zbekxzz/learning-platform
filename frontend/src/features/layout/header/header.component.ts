import { Component } from '@angular/core';
import { Router, RouterLink } from '@angular/router';
import { LocalStorageService } from '../../../shared/services/local-storage.service';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [RouterLink],
  templateUrl: './header.component.html',
})
export class HeaderComponent {

  constructor(
    private router: Router,
    private storage: LocalStorageService
  ) {}

  logout() {
    this.storage.removeTokenAndRedirect();
    this.router.navigate(['/auth/login']);
  }

}
