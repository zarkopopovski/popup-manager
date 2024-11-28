import { CommonModule } from '@angular/common';
import { Component, inject } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { RouterLink, Router } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { HttpErrorResponse } from '@angular/common/http';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [RouterLink, CommonModule, FormsModule],
  templateUrl: './login.component.html',
  styleUrl: './login.component.css'
})
export class LoginComponent {
  private authService = inject(AuthService);
  private router = inject(Router);

  public email: string | undefined;
  public password: string | undefined;

  public isProgressStarted: boolean = false;

  constructor(){}

  loginToDashboard() {
    this.isProgressStarted = true;
    this.authService.loginUser(this.email, this.password).subscribe({
      next: (response: any) => {
        this.isProgressStarted = false;
        sessionStorage.setItem("user_data", JSON.stringify(response.data));
        
        this.router.navigateByUrl("/dashboard");
      },
      error: (err: HttpErrorResponse) => {
        this.isProgressStarted = false;
        console.log(err);
      }
    });
  }
}
