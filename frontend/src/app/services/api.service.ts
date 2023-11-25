import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { HttpClient, HttpParams, HttpResponse } from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  private _apiRoot: string = '/api';

  constructor(
    protected _http: HttpClient,
  ) { }

  get(path: string, params?: HttpParams): Observable<object> {
    const options: any = {};
    if (params) {
      options.params = params;
    }
    return this._http.get(`${this._apiRoot}/${path}`, options);
  }

  protected _getMutateRequestHeaders(): any {
    return {
      contentType: 'application/json',
    };
  }

  post(path: string, body: any | null): Observable<HttpResponse<object>> {
    return this._http.post(`${this._apiRoot}/${path}`, body, {
      headers: this._getMutateRequestHeaders(),
      responseType: 'json',
      observe: 'response',
    });
  }

  login(username: string, password: string): Observable<HttpResponse<any>> {
    return this.post('auth/login/', {
      'username': username,
      'password': password,
    }) as Observable<HttpResponse<any>>;
  }

  logout(): Observable<HttpResponse<object>> {
    return this.post('auth/logout/', null);
  }

  register(
    username: string,
    password: string,
  ): Observable<HttpResponse<any>> {
    return this.post(
      `auth/register/`,
      { username: username, password: password },
    ) as Observable<HttpResponse<any>>;
  }
}
