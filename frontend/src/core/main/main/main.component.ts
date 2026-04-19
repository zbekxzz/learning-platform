import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { SidebarComponent } from '../../../features/layout/sidebar/sidebar.component';
import { HeaderComponent } from '../../../features/layout/header/header.component';
import { LayoutService } from '../../../shared/services/layout.service';

@Component({
  selector: 'app-main',
  standalone: true,
  imports: [RouterOutlet, SidebarComponent, HeaderComponent],
  templateUrl: './main.component.html',
})
export class MainComponent {
  constructor(public layout: LayoutService) {}
}
