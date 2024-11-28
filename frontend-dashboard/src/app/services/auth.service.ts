import { HttpClient, HttpErrorResponse, HttpHeaders } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable, Subject, catchError, map, throwError } from 'rxjs';

import { 
  BASE_URL, 
  API_PATH,
  USER_PATH,
  REGISTER,
  LOGIN,
  REFRESH_TOKEN, 
  CONFIRMATION,
  RESET_PASSWORD,
  CHANGE_PASSWORD} from '../shared/constants';

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  private readonly http = inject(HttpClient)

  public logoutSubject: Subject<boolean> = new Subject<boolean>();

  constructor() { }

  registerNewUser(username: string | undefined, password: string | undefined): Observable<any> {
    let httpHeaders = new HttpHeaders()
    .set('Content-Type', 'application/json');

    const httpOptions = {
      headers: httpHeaders,
    };

    const userData = {
      'email': username,
      'password': password,
    };

    return this.http.post(BASE_URL+API_PATH+REGISTER, userData, httpOptions);
  }

  loginUser(username: string | undefined, password: string | undefined): Observable<any> {
    let httpHeaders = new HttpHeaders()
    .set('Content-Type', 'application/json');

    const httpOptions = {
      headers: httpHeaders,
    };

    const userData = {
      'email': username,
      'password': password,
    };

    return this.http.post(BASE_URL+API_PATH+LOGIN, userData, httpOptions);
  }

  refreshTokenRequest(refreshToken: string): Observable<any> {
    let httpHeaders = new HttpHeaders()
    .set('Content-Type', 'application/json');

    const httpOptions = {
      headers: httpHeaders,
    };

    return this.http.get<any>(BASE_URL+API_PATH+USER_PATH+REFRESH_TOKEN+'/'+refreshToken, httpOptions);
  }

  confirmRegistration(confirmationToken: string): Observable<any> {
    let httpHeaders = new HttpHeaders()
    .set('Content-Type', 'application/json');

    const httpOptions = {
      headers: httpHeaders,
    };

    return this.http.get<any>(BASE_URL+API_PATH+CONFIRMATION+'/'+confirmationToken, httpOptions);
  }

  resetPassword(username: string | undefined): Observable<any> {
    let httpHeaders = new HttpHeaders()
    .set('Content-Type', 'application/json');

    const httpOptions = {
      headers: httpHeaders,
    };

    const userData = {
      'email': username,
    };

    return this.http.post(BASE_URL+API_PATH+RESET_PASSWORD, userData, httpOptions);
  }

  changeUserPassword(accessToken: string, password: string | undefined): Observable<any> {
    let httpHeaders = new HttpHeaders()
    .set('Content-Type', 'application/json')
    .set('Authorization', 'Bearer ' + accessToken);

    const httpOptions = {
      headers: httpHeaders,
    };

    const userData = {
      'password': password,
    };

    return this.http.post(BASE_URL+API_PATH+USER_PATH+CHANGE_PASSWORD, userData, httpOptions);
  }

  pushToLogout(logoutStatus: boolean) {
    this.logoutSubject.next(true);
  }
}
