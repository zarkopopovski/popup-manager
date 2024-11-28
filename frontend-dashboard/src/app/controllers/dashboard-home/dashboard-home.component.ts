import { Component, inject } from '@angular/core';
import { Router, RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';
import { User } from '../../models/user';
import { Subscription } from 'rxjs';
import { AuthService } from '../../services/auth.service';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-dashboard-home',
  standalone: true,
  imports: [RouterOutlet, RouterLink, RouterLinkActive, CommonModule],
  templateUrl: './dashboard-home.component.html',
  styleUrl: './dashboard-home.component.css'
})
export class DashboardHomeComponent {
  private router = inject(Router);
  private authService = inject(AuthService);

  private logoutSubscription?: Subscription;

  public userData: User | undefined = undefined;

  constructor() {
    this.logoutSubscription = this.authService.logoutSubject.subscribe(data => {
      if (data == true) {
        sessionStorage.removeItem("user_data");
        this.router.navigateByUrl("/login");
      }
    });

    let userData: string = sessionStorage.getItem("user_data")!;
    if (!userData) {
      this.router.navigateByUrl("/login");
    }
    this.userData = JSON.parse(userData);
  }

  toggleBar() {
    /**
     * Easy selector helper function
     */
    const select = (el: string, all = false): HTMLElement | HTMLElement[] => {
      el = el.trim();
      if (all) {
        return Array.from(document.querySelectorAll(el));
      } else {
        return document.querySelector(el) as HTMLElement;
      }
    };

    /**
     * Easy event listener function
     */
    const on = (
      type: string,
      el: string,
      listener: EventListenerOrEventListenerObject,
      all = false
    ): void => {
      if (all) {
        (select(el, all) as HTMLElement[]).forEach((e) =>
          e.addEventListener(type, listener)
        );
      } else {
        (select(el, all) as HTMLElement).addEventListener(type, listener);
      }
    };

    /**
     * Sidebar toggle
     */
    const body = select('body') as HTMLElement;
    body.classList.toggle('toggle-sidebar');
  }

  signOutFromDashboard() {
    sessionStorage.removeItem("user_data");
  }
}
