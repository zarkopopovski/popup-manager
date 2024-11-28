import { HttpErrorResponse, HttpHandlerFn, HttpInterceptorFn, HttpRequest } from "@angular/common/http";
import { inject } from "@angular/core";
import { Subject, catchError, switchMap } from "rxjs";
import { User } from "../models/user";
import { AuthService } from "./auth.service";

let refreshingToken$: Subject<void> | null = null;

let userData: User | undefined = undefined;

export const refreshTokenInterceptor: HttpInterceptorFn = (
  req: HttpRequest<unknown>,
  next: HttpHandlerFn
) => {
  const authService = inject(AuthService);

  return next(req).pipe(
    catchError((err) => {
      if (err.status === 401) {    
        if (err.url.includes('/refresh-token')) {
          authService.pushToLogout(true);
        }    
        if (!refreshingToken$) {
          refreshingToken$ = new Subject<void>();
          let userDataStr: string = sessionStorage.getItem("user_data")!;
          if (userDataStr) {
            let userData: User = JSON.parse(userDataStr);
            
            authService.refreshTokenRequest(userData.tokens?.RefreshToken!).subscribe({
              next: (resp: any) => {
                console.log(resp)
                refreshingToken$!.next();
                refreshingToken$!.complete();
                refreshingToken$ = null;
              },
              error: (err: HttpErrorResponse) => {
                authService.pushToLogout(true);
              },
            });
          }
        }
        return refreshingToken$.pipe(switchMap(() => next(req)));
      }

      throw err;
    })
  );
};