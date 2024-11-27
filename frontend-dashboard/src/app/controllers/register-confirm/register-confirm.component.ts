import { Component, inject } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { HttpErrorResponse } from '@angular/common/http';

@Component({
  selector: 'app-register-confirm',
  standalone: true,
  imports: [],
  templateUrl: './register-confirm.component.html',
  styleUrl: './register-confirm.component.css'
})
export class RegisterConfirmComponent {
  private route = inject(ActivatedRoute);
  private authService = inject(AuthService);

  public confirmationSuccessful: boolean = false;
  public isError: boolean = false;
  public errorMessage?: string = '';

  constructor(){
    let confirmationToken = this.route.snapshot.paramMap.get('confirmationToken')!;

    this.confirmRegistrationWithToken(confirmationToken);
  }

  confirmRegistrationWithToken(confirmationToken: string) {
    this.authService.confirmRegistration(confirmationToken).subscribe({
      next: (response: any) => {
        this.confirmationSuccessful = true;
        this.errorMessage = 'Registration confirmation was successful. You can now log in to your account.';
      },
      error: (err: HttpErrorResponse) => {
        this.isError = true;
        this.errorMessage = err.error.error;
      }
    });
  }

}
