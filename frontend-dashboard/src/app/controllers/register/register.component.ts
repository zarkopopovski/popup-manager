import { Component, inject } from '@angular/core';
import { AuthService } from '../../services/auth.service';
import { HttpErrorResponse } from '@angular/common/http';
import { CommonModule } from '@angular/common';
import { FormsModule, NgModel } from '@angular/forms';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-register',
  standalone: true,
  imports: [RouterLink, CommonModule, FormsModule],
  templateUrl: './register.component.html',
  styleUrl: './register.component.css'
})
export class RegisterComponent {
  private authService = inject(AuthService);

  public username?: string | undefined;
  public password?: string | undefined;

  public isRegistrationSuccessful: boolean = false;
  public isError: boolean = false;
  public errorMessage?: string = 'There seems to be an issue with your registration.';

  constructor(){}

  registerNewUser() {
    this.authService.registerNewUser(this.username, this.password).subscribe({
      next: (response: any) => {
        console.log(response);
        if (response.error_code === '-1') {
          this.isRegistrationSuccessful = true;
        } else {
          this.isError = true;
        }
      },
      error: (err:HttpErrorResponse) => {
        this.isError = true;
        
        if (err.status == 403) {
          this.errorMessage = err.error.error;
        }
      }
    });
  }

}
