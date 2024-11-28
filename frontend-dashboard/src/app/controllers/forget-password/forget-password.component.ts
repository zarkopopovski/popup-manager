import { CommonModule } from '@angular/common';
import { Component, inject } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { AuthService } from '../../services/auth.service';
import { HttpErrorResponse } from '@angular/common/http';

@Component({
  selector: 'app-forget-password',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './forget-password.component.html',
  styleUrl: './forget-password.component.css'
})
export class ForgetPasswordComponent {
  private authService = inject(AuthService);

  public email: string | undefined;
  public password: string | undefined;

  public isProgressStarted: boolean = false;
  public requestFinished: boolean = false;

  passwordRequest() {
    this.isProgressStarted = true;
    this.authService.resetPassword(this.email).subscribe({
      next: (request: any) => {
        this.isProgressStarted = false;
        this.requestFinished = true;
      },
      error: (err: HttpErrorResponse) => {
        this.isProgressStarted = false;
        this.requestFinished = true;
      }
    });
  }
}
