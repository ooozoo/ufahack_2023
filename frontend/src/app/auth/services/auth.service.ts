import { Injectable } from '@angular/core';
import { HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';
import { tap } from 'rxjs/operators'
import { ApiService } from '../../services/api.service';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  constructor(
    protected _api: ApiService,
  ) { }

  login(username: string, password: string): Observable<HttpResponse<any>> {
    return this._api.login(username, password).pipe(tap((response)=> {
      this.setUser(response.body);
    }));
  }

  logout(): Observable<HttpResponse<object>> {
    return this._api.logout().pipe(tap(() => {
      this.removeUser();
    }));
  }

  register(
    username: string,
    password: string,
  ): Observable<HttpResponse<any>> {
    return this._api.register(username, password).pipe(tap((response)=> {
      this.setUser(response.body);
    }));
  }

  isAuthenticated(): boolean {
    return !!this.getUser();
  }

  getUser(): any | null {
    const strUser = localStorage.getItem('user');
    if (strUser) {
      try {
        return JSON.parse(strUser);
      } catch {
        return null;
      }
    }
    return null;
  }

  setUser(user: any): void {
    localStorage.setItem('user', JSON.stringify(user));
  }

  removeUser() {
    localStorage.removeItem('user');
  }
}
