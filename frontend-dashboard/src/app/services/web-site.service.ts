import { Injectable, inject } from '@angular/core';

import { 
  BASE_URL,
  API_PATH,
  USER_PATH,
  WEB_SITE,
  NOTIFICATION, 
  SIMPLE_STATS,
  LATEST_X_VISITS,
  SIMPLE_STATS_VISITS} from '../shared/constants';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { WebTokensResponse } from '../models/web-tokens-response';

interface WebSiteRequest {
  title?: string;
  description?: string;
  web_url?: string;
}

@Injectable({
  providedIn: 'root'
})
export class WebSiteService {
  private http = inject(HttpClient)

  constructor() { }

  createNewWebToken(title: string, description: string | null, webURL: string, accessToken: string): Observable<any> {
    let httpHeaders = new HttpHeaders()
    .set('Content-Type', 'application/json');
    httpHeaders = httpHeaders.append('Authorization', 'Bearer ' + accessToken);

    const httpOptions = {
      headers: httpHeaders,
    };

    const userData = {
      'title': title,
      'description': description,
      'web_url': webURL
    };

    return this.http.post<WebSiteRequest>(BASE_URL+API_PATH+USER_PATH+WEB_SITE, userData, httpOptions);
  }
  
  updateWebToken(title: string, description: string | null, webURL: string, accessToken: string, webToken: string): Observable<any> {
    let httpHeaders = new HttpHeaders()
    .set('Content-Type', 'application/json')
    .set('Authorization', 'Bearer ' + accessToken);

    const httpOptions = {
      headers: httpHeaders,
    };

    const userData = {
      'title': title,
      'description': description,
      'web_url': webURL
    };

    return this.http.put<any>(BASE_URL+API_PATH+USER_PATH+WEB_SITE+'/'+webToken, userData, httpOptions);
  }

  deleteWebToken(accessToken: string, webToken: string): Observable<any> {
    let httpHeaders = new HttpHeaders()
    .set('Content-Type', 'application/json')
    .set('Authorization', 'Bearer ' + accessToken);

    const httpOptions = {
      headers: httpHeaders,
    };

    return this.http.delete<any>(BASE_URL+API_PATH+USER_PATH+WEB_SITE+'/'+webToken, httpOptions);
  }

  listWebTokens(accessToken: string): Observable<WebTokensResponse> {
    let httpHeaders = new HttpHeaders()
    .set('Content-Type', 'application/json')
    .set('Authorization', 'Bearer ' + accessToken);

    const httpOptions = {
      headers: httpHeaders,
    };

    return this.http.get<WebTokensResponse>(BASE_URL+API_PATH+USER_PATH+WEB_SITE, httpOptions);
  }

  listPopupMessages(accessToken: string, apiToken: string): Observable<any> {
    let httpHeaders = new HttpHeaders()
    .set('Content-Type', 'application/json')
    .set('Authorization', 'Bearer ' + accessToken);

    const httpOptions = {
      headers: httpHeaders,
    };

    return this.http.get<any>(BASE_URL+API_PATH+USER_PATH+NOTIFICATION+'/'+apiToken, httpOptions);
  }

  createNewWebPopup(title: string, description: string | null, showTime: number, closeTime: number, popupType: number, accessToken: string, apiToken: string, popupPos: number, imageFile: any, isTrackable: boolean): Observable<any> {
    let httpHeaders = new HttpHeaders()
    httpHeaders = httpHeaders.append('Authorization', 'Bearer ' + accessToken);

    let params = new FormData();
    params.set('api_token', apiToken);
    params.set('title', title);
    params.set('description', description!);
    params.set('show_time', String(showTime));
    params.set('close_time', String(closeTime));
    params.set('type', String(popupType));
    params.set('popup_pos', String(popupPos));
    params.set('file', imageFile);
    params.set('is_trackable', String((isTrackable)?1:0));

    const httpOptions = {
      headers: httpHeaders,
    };

    return this.http.post<WebSiteRequest>(BASE_URL+API_PATH+USER_PATH+NOTIFICATION, params, httpOptions);
  }
  
  updateWebPopup(title: string, description: string | null, showTime: number, closeTime: number, popupType: number, accessToken: string, apiToken: string, notificationID: number, popupPos: string, imageFile: any, isTrackable: boolean): Observable<any> {
    let httpHeaders = new HttpHeaders()
    .set('Authorization', 'Bearer ' + accessToken);

    const httpOptions = {
      headers: httpHeaders,
    };

    let params = new FormData();
    params.set('api_token', apiToken);
    params.set('title', title);
    params.set('description', description!);
    params.set('show_time', String(showTime));
    params.set('close_time', String(closeTime));
    params.set('type', String(popupType));
    params.set('popup_pos', popupPos);
    params.set('file', imageFile);
    params.set('is_trackable', String((isTrackable)?1:0));

    return this.http.put<any>(BASE_URL+API_PATH+USER_PATH+NOTIFICATION+'/'+apiToken+'/'+notificationID, params, httpOptions);
  }

  deleteWebPopup(accessToken: string, apiToken: string, notificationID: number): Observable<any> {
    let httpHeaders = new HttpHeaders()
    .set('Content-Type', 'application/json')
    .set('Authorization', 'Bearer ' + accessToken);

    const httpOptions = {
      headers: httpHeaders,
    };

    return this.http.delete<any>(BASE_URL+API_PATH+USER_PATH+NOTIFICATION+'/'+apiToken+'/'+notificationID, httpOptions);
  }

  getPopupSimpleStats(accessToken: string): Observable<any> {
    let httpHeaders = new HttpHeaders()
    .set('Content-Type', 'application/json')
    .set('Authorization', 'Bearer ' + accessToken);

    const httpOptions = {
      headers: httpHeaders,
    };

    return this.http.get<any>(BASE_URL+API_PATH+USER_PATH+SIMPLE_STATS, httpOptions);
  }

  getPopupSimpleStatsVisits(accessToken: string, latestXVisits: number): Observable<any> {
    let httpHeaders = new HttpHeaders()
    .set('Content-Type', 'application/json')
    .set('Authorization', 'Bearer ' + accessToken);

    const httpOptions = {
      headers: httpHeaders,
    };

    return this.http.get<any>(BASE_URL+API_PATH+USER_PATH+SIMPLE_STATS_VISITS+'/'+latestXVisits, httpOptions);
  }
}
