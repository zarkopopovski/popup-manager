import { CommonModule } from '@angular/common';
import { Component, inject } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { WebToken } from '../../models/web-token';
import { WebSiteService } from '../../services/web-site.service';
import { User } from '../../models/user';
import { HttpErrorResponse, HttpResponse } from '@angular/common/http';
import { WebTokensResponse } from '../../models/web-tokens-response';
import { BrowserModule } from '@angular/platform-browser';

@Component({
  selector: 'app-dashboard-websites',
  standalone: true,
  imports: [RouterLink, CommonModule, FormsModule],
  templateUrl: './dashboard-websites.component.html',
  styleUrl: './dashboard-websites.component.css'
})
export class DashboardWebsitesComponent {
  private router = inject(Router);
  private webSiteService = inject(WebSiteService)

  public isFormDisplayed: boolean = false;
  public isEditMode: boolean = false;
  public isShowMode: boolean = false;

  public webTitle?: string = '';
  public webDescription?: string = '';
  public webUrl?: string = '';

  public selectedWebToken?: WebToken | null;
  public webTokensList?: WebToken[];

  public userData: User | undefined = undefined;

  public isProgressStarted: boolean = false;

  public isTokenDeleteAction: boolean = false;

  constructor() {
    let userData: string = sessionStorage.getItem("user_data")!;
    if (!userData) {
      this.router.navigateByUrl("/login");
    }
    this.userData = JSON.parse(userData);

    this.listAllWebTokens();
  }

  resetFormData() {
    this.isFormDisplayed = false;
    this.isEditMode = false;
    this.isShowMode = false;
    this.webTitle = '';
    this.webDescription = '';
    this.webUrl = '';
    this.selectedWebToken = null;
  }

  executeWebTokenRequest() {
    console.log("CREATE NEW OR UPDATE EXISTING WEB TOKEN");
    this.isProgressStarted = true;
    if (!this.isEditMode) {
      this.webSiteService.createNewWebToken(this.webTitle!, this.webDescription!, this.webUrl!, this.userData?.tokens?.AccessToken!).subscribe({
        next: (response: any) => {
          console.log(response);
          this.listAllWebTokens();
          this.isProgressStarted = false;
          this.resetFormData();
        },
        error: (err: HttpErrorResponse) => {
        this.isProgressStarted = false;
          console.log(err);
        }
      });
    } else {
      this.webSiteService.updateWebToken(this.webTitle!, this.webDescription!, this.webUrl!, this.userData?.tokens?.AccessToken!, this.selectedWebToken?.api_token!).subscribe({
        next: (response: any) => {
          console.log(response);
          this.listAllWebTokens();
          this.isProgressStarted = false;
          this.resetFormData();
        },
        error: (err: HttpErrorResponse) => {
          this.isProgressStarted = false;
          console.log(err);
        }
      });
    }
  }

  listAllWebTokens() {
    this.webSiteService.listWebTokens(this.userData?.tokens?.AccessToken!).subscribe({
      next: (response: WebTokensResponse) => {
        this.webTokensList = response.data;
      },
      error: (err: HttpErrorResponse) => {
        console.log(err);
      }
    });
  }

  showCurrentTokenDetails(webToken: WebToken) {
    this.selectedWebToken = webToken;
    this.webTitle = webToken.title;
    this.webDescription = webToken.description;
    this.webUrl = webToken.web_url;
    this.isFormDisplayed = true;
    this.isShowMode = true;
  }

  editCurrentToken(webToken: WebToken) {
    console.log(webToken);
    this.selectedWebToken = webToken;
    this.isEditMode = true;
    this.webTitle = webToken.title;
    this.webDescription = webToken.description;
    this.webUrl = webToken.web_url;
    this.isFormDisplayed = true;
  }

  deleteCurrentToken(webToken: WebToken){
    console.log(webToken);
    this.selectedWebToken = webToken;
    this.isTokenDeleteAction = true;
  }

  closeDeleteModal() {
    this.selectedWebToken = null;;
    this.isTokenDeleteAction = false;
  }

  deleteCurrentTokenAction() {
    this.webSiteService.deleteWebToken(this.userData?.tokens?.AccessToken!, this.selectedWebToken?.api_token!).subscribe({
      next: (request: unknown) => {
        this.selectedWebToken = null;
        this.isTokenDeleteAction = false;

        this.listAllWebTokens();
      },
      error: (error: HttpErrorResponse) => {
        console.log(error);
        this.selectedWebToken = null;
        this.isTokenDeleteAction = false;
      }
    });
  }

  showPopupsSite(webToken: WebToken) {
    this.router.navigateByUrl("/dashboard/popup-messages/" + webToken.api_token);
  }
}
