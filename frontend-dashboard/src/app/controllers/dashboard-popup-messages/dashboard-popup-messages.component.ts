import { Component, inject } from '@angular/core';
import { Router, RouterLink, ActivatedRoute } from '@angular/router';
import { User } from '../../models/user';
import { WebSiteService } from '../../services/web-site.service';
import { HttpErrorResponse, HttpEvent } from '@angular/common/http';
import { PopUpMessage } from '../../models/popup-message';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { BASE_URL } from '../../shared/constants';

@Component({
  selector: 'app-dashboard-popup-messages',
  standalone: true,
  imports: [RouterLink, CommonModule, FormsModule],
  templateUrl: './dashboard-popup-messages.component.html',
  styleUrl: './dashboard-popup-messages.component.css'
})
export class DashboardPopupMessagesComponent {
  private router = inject(Router);
  private route = inject(ActivatedRoute);
  private webSiteService = inject(WebSiteService)

  public isFormDisplayed: boolean = false;
  public isEditMode: boolean = false;
  public isShowMode: boolean = false;

  public webTitle?: string = '';
  public webDescription?: string = '';

  public popupShowAfter: number = 0;
  public popupCloseAfter: number = 0;

  public selectedWebPopup?: PopUpMessage | null;
  public webPopupsList?: PopUpMessage[];

  public userData: User | undefined = undefined;

  public isProgressStarted: boolean = false;

  public isPopupDeleteAction: boolean = false;

  public currentApiToken: string = '';

  public popupType: string = '';
  public popupPos: string = '';

  public imageName: string = '';

  public uploadedLogoFile: any;

  public isTrackable: boolean = false;

  public imagesDomain: string = '';
  
  constructor() {
    let userData: string = sessionStorage.getItem("user_data")!;
    if (!userData) {
      this.router.navigateByUrl("/login");
    }
    this.userData = JSON.parse(userData);

    this.currentApiToken = this.route.snapshot.paramMap.get('apiToken')!;

    this.imagesDomain = BASE_URL;

    this.listAllPopupMessages(this.currentApiToken);
  }

  resetFormData() {
    this.isFormDisplayed = false;
    this.isEditMode = false;
    this.isShowMode = false;
    this.webTitle = '';
    this.webDescription = '';
    this.popupType = '';
    this.popupPos = '';
    this.imageName = '';
    this.selectedWebPopup = null;
    this.uploadedLogoFile = null;

    this.popupShowAfter = 0;
    this.popupCloseAfter = 0;

    this.isTrackable = false;
  }

  listAllPopupMessages(apiToken: string) {
    this.webSiteService.listPopupMessages(this.userData?.tokens?.AccessToken!, apiToken).subscribe({
      next: (response: any) => {
        this.webPopupsList = response.data;
      },
      error: (err: HttpErrorResponse) => {
        console.log(err);
      }
    });
  }

  showCurrentPopupDetails(webPopup: PopUpMessage) {
    this.selectedWebPopup = webPopup;
    this.webTitle = webPopup.title;
    this.webDescription = webPopup.description;
    this.popupShowAfter = webPopup.show_time!;
    this.popupCloseAfter = webPopup.close_time!;

    this.popupType = '' + webPopup.pop_type!;
    this.popupPos = '' + webPopup.popup_pos!;

    this.imageName = Object.values(webPopup.image_name!)[0];

    this.isTrackable = webPopup.isTrackable!;

    this.isFormDisplayed = true;
    this.isShowMode = true;
  }

  returnPopTypePerNumber(popType: number) {
    switch(popType) {
      case 1:
        return '-';
      case 2:
        return 'Success';
      case 3:
        return 'Info';
      case 4:
        return 'Warning';
      case 5:
        return 'Error';
      default:
        return '-';
    }
  }

  uploadSelectedFile(event: any, index: HTMLInputElement) {
    let fileForUpload = event.target.files[0];
    let fileExtension = fileForUpload.name.replace(/^.*\./, '');

    let acept = ['jpg','jpeg', 'gif', 'png', 'ico'];

    if ((acept.indexOf(fileExtension.toLowerCase()) == -1) || (fileForUpload.size > (0.5 * 1024 * 1024))) {
      alert('The file is more then 500k big or incorrect format');
      return;
    }

    this.uploadedLogoFile = event.target.files[0];
  }

  executeWebPopupRequest() {
    console.log("CREATE NEW OR UPDATE EXISTING WEB TOKEN");
    this.isProgressStarted = true;
    if (!this.isEditMode) {
      this.webSiteService.createNewWebPopup(
        this.webTitle!, 
        this.webDescription!, 
        this.popupShowAfter, 
        this.popupCloseAfter, 
        Number(this.popupType), 
        this.userData?.tokens?.AccessToken!, 
        this.currentApiToken,
        Number(this.popupPos),
        this.uploadedLogoFile,
        this.isTrackable).subscribe({
        next: (response: any) => {
          console.log(response);
          this.listAllPopupMessages(this.currentApiToken);
          this.isProgressStarted = false;
          this.resetFormData();
        },
        error: (err: HttpErrorResponse) => {
        this.isProgressStarted = false;
          alert(err.error.error);
        }
      });
    } else {
      this.webSiteService.updateWebPopup(
        this.webTitle!, 
        this.webDescription!, 
        this.popupShowAfter, 
        this.popupCloseAfter, 
        Number(this.popupType), 
        this.userData?.tokens?.AccessToken!, 
        this.currentApiToken, 
        this.selectedWebPopup?.id!,
        this.popupPos,
        this.uploadedLogoFile,
        this.isTrackable).subscribe({
        next: (response: any) => {
          console.log(response);
          this.listAllPopupMessages(this.currentApiToken);
          this.isProgressStarted = false;
          this.resetFormData();
        },
        error: (err: HttpErrorResponse) => {
          this.isProgressStarted = false;
          alert(err.error.error);
        }
      });
    }
  }

  editCurrentPopup(webPopup: PopUpMessage) {
    this.selectedWebPopup = webPopup;
    this.isEditMode = true;
    this.webTitle = webPopup.title;
    this.webDescription = webPopup.description;
    this.popupShowAfter = webPopup.show_time!;
    this.popupCloseAfter = webPopup.close_time!;

    this.popupType = '' + webPopup.pop_type!;
    this.popupPos = '' + webPopup.popup_pos!;

    this.imageName = Object.values(webPopup.image_name!)[0];

    this.isTrackable = webPopup.isTrackable!;

    this.isFormDisplayed = true;
  }

  deleteCurrentPopup(webPopup: PopUpMessage) {
    this.selectedWebPopup = webPopup;
    this.isPopupDeleteAction = true;
  }

  deleteCurrentPopupAction() {
    this.webSiteService.deleteWebPopup(this.userData?.tokens?.AccessToken!, this.currentApiToken, this.selectedWebPopup?.id!).subscribe({
      next: (request: unknown) => {
        this.selectedWebPopup = null;
        this.isPopupDeleteAction = false;

        this.listAllPopupMessages(this.currentApiToken);
      },
      error: (error: HttpErrorResponse) => {
        console.log(error);
        this.selectedWebPopup = null;
        this.isPopupDeleteAction = false;
      }
    });
  }

  closeDeleteModal() {
    this.selectedWebPopup = null;;
    this.isPopupDeleteAction = false;
  }
}
