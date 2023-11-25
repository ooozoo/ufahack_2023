import {Inject, Injectable} from '@angular/core';
import {
  HttpRequest,
  HttpHandler,
  HttpEvent,
  HttpInterceptor, HttpErrorResponse
} from '@angular/common/http';
import {catchError, Observable, throwError} from 'rxjs';
import {TuiAlertService} from "@taiga-ui/core";
import {AuthService} from "../services/auth.service";
import {Router} from "@angular/router";

@Injectable()
export class AuthInterceptor implements HttpInterceptor {

  constructor(
    @Inject(TuiAlertService) private _alerts: TuiAlertService,
    private _router: Router,
    private _authService: AuthService,
  ) { }

  intercept(request: HttpRequest<unknown>, next: HttpHandler): Observable<HttpEvent<unknown>> {
    try {
      return next.handle(request).pipe(
        catchError((error: HttpErrorResponse) => {
          if (error.status === 401) {
            this._authService.removeUser();
            const next = window.location.pathname;
            if (next === '/' || next === '/login') {
              this._router.navigate(['/login']);
            } else {
              this._router.navigate(['/login'], {queryParams: {next: next}});
            }
          }
          if (error.status > 499) {
            this.showErrorNotification();
          }
          return throwError(error);
        }),
      );
    } catch {
      return next.handle(request).pipe(
        catchError((error: HttpErrorResponse) => {
          if (error.status > 499) {
            this.showErrorNotification();
          }

          return throwError(error);
        }),
      );
    }
  }

  showErrorNotification(): void {
    this._alerts
      .open('Внутренняя ошибка сервера. Попробуйте позже.', {label: ' Ошибка'})
      .subscribe();
  }
}
