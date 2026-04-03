import { Component, OnInit, signal } from '@angular/core';
import { RouterLink, RouterLinkActive } from '@angular/router';
import { ProfileService } from '../../../pages/profile/services/profile.service';
import { UserRole } from '../../../pages/profile/models/profile.model';

@Component({
  selector: 'app-sidebar',
  imports: [RouterLink, RouterLinkActive],
  templateUrl: './sidebar.component.html',
})
export class SidebarComponent implements OnInit {
  role = signal<UserRole | null>(null);

  constructor(private profileService: ProfileService) {}

  ngOnInit() {
    this.profileService.getProfile().subscribe(profile => {
      if (profile) {
        this.role.set(profile.role);
      }
    });
  }
}
