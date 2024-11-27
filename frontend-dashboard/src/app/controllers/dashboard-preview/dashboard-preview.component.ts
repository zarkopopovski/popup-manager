import { Component, inject } from '@angular/core';
import { Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { User } from '../../models/user';
import { AuthService } from '../../services/auth.service';
import { WebSiteService } from '../../services/web-site.service';
import { HttpErrorResponse } from '@angular/common/http';
import { BasicStat } from '../../models/basic-stat';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-dashboard-preview',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './dashboard-preview.component.html',
  styleUrl: './dashboard-preview.component.css'
})
export class DashboardPreviewComponent {
  private router = inject(Router);
  private authService = inject(AuthService);
  private webSiteService = inject(WebSiteService)

  private logoutSubscription?: Subscription;

  public userData: User | undefined = undefined;

  public numWebSites: number = 0;
  public numWebPopups: number = 0;
  public numClicks: number = 0;

  public basicStatsVisits: BasicStat[] = [];

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

    this.getSimpleStats();
    this.getSimpleStatsVisits();
  }

  getSimpleStats() {
    this.webSiteService.getPopupSimpleStats(this.userData?.tokens?.AccessToken!).subscribe({
      next: (res: any) => {
        this.numWebSites = res.data.num_site_tokens;
        this.numWebPopups = res.data.num_web_popups;
        this.numClicks = res.data.num_clicks;
      },
      error: (err: HttpErrorResponse) => {
        console.log(err.error.error);
      }
    });
  }

  getSimpleStatsVisits() {
    this.webSiteService.getPopupSimpleStatsVisits(this.userData?.tokens?.AccessToken!, 10).subscribe({
      next: (res: any) => {
        this.basicStatsVisits = res.data;
      },
      error: (err: HttpErrorResponse) => {
        console.log(err.error.error);
      }
    });
  }
}
