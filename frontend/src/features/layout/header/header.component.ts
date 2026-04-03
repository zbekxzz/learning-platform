import { Component, OnInit, signal } from '@angular/core';
import { Router, RouterLink } from '@angular/router';
import { LocalStorageService } from '../../../shared/services/local-storage.service';
import { ProfileService } from '../../../pages/profile/services/profile.service';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [RouterLink],
  templateUrl: './header.component.html',
})
export class HeaderComponent implements OnInit {
  initial = signal<string>('U');

  constructor(
    private router: Router,
    private storage: LocalStorageService,
    private profileService: ProfileService
  ) {}

  ngOnInit() {
    this.profileService.getProfile().subscribe(profile => {
      if (profile && profile.full_name) {
        this.initial.set(profile.full_name[0].toUpperCase());
      }
    });
  }

  logout() {
    this.profileService.clearCache();
    this.storage.removeTokenAndRedirect();
    this.router.navigate(['/auth/login']);
  }

}
