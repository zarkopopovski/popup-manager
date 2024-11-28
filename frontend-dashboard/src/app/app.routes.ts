import { Routes } from '@angular/router';
import { LoginComponent } from './controllers/login/login.component';
import { RegisterComponent } from './controllers/register/register.component';
import { DashboardHomeComponent } from './controllers/dashboard-home/dashboard-home.component';
import { DashboardPreviewComponent } from './controllers/dashboard-preview/dashboard-preview.component';
import { DashboardPopupMessagesComponent } from './controllers/dashboard-popup-messages/dashboard-popup-messages.component';
import { DashboardWebsitesComponent } from './controllers/dashboard-websites/dashboard-websites.component';
import { RegisterConfirmComponent } from './controllers/register-confirm/register-confirm.component';
import { ForgetPasswordComponent } from './controllers/forget-password/forget-password.component';
import { DashboardChangePasswordComponent } from './controllers/dashboard-change-password/dashboard-change-password.component';

export const routes: Routes = [
    {path: '', component: LoginComponent},
    {path: 'login', component: LoginComponent},
    {path: 'register', component: RegisterComponent},
    {path: 'confirm-registartion/:confirmationToken', component: RegisterConfirmComponent},
    {path: 'forget-password', component: ForgetPasswordComponent},
    {path: 'dashboard', component: DashboardHomeComponent, children: [
        {path: '', component: DashboardPreviewComponent},
        {path: 'web-sites', component: DashboardWebsitesComponent},
        {path: 'change-password', component: DashboardChangePasswordComponent},
        {path: 'popup-messages/:apiToken', component: DashboardPopupMessagesComponent},
    ]},
];
