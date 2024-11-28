import { CommonModule } from '@angular/common';
import { Component, inject } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { HttpErrorResponse } from '@angular/common/http';
import { User } from '../../models/user';

@Component({
  selector: 'app-dashboard-change-password',
  standalone: true,
  imports: [RouterLink, CommonModule, FormsModule],
  templateUrl: './dashboard-change-password.component.html',
  styleUrl: './dashboard-change-password.component.css'
})
export class DashboardChangePasswordComponent {
  private router = inject(Router);
  private authService = inject(AuthService);

  public userData: User | undefined = undefined;

  public isProgressStarted: boolean = false;  

  public passwordChanged: boolean = false;
  public isError: boolean = false;

  public newPassword: string = '';
  public newPasswordConfirm: string = '';

  constructor() {
    let userData: string = sessionStorage.getItem("user_data")!;
    if (!userData) {
      this.router.navigateByUrl("/login");
    }
    this.userData = JSON.parse(userData);
  }

  resetFormData() {
    this.newPassword = '';
    this.newPasswordConfirm = '';
  }

  executeChangePasswordRequest() {
    this.isProgressStarted = true;
    this.authService.changeUserPassword(this.userData?.tokens?.AccessToken!, this.newPassword).subscribe({
      next: (response: any) => {
        this.isProgressStarted = false;
        this.passwordChanged = true;
      },
      error: (err: HttpErrorResponse) => {
        this.isProgressStarted = false;
        this.isError = true;
      }
    });
  }
}
