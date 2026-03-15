import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { LocalStorageService } from './local-storage.service';
import { Observable } from 'rxjs';

type QueryParams = Record<string, string | number | boolean>;

function toHttpParams(params?: QueryParams): HttpParams {
  let httpParams = new HttpParams();
  if (!params) return httpParams;

  for (const key in params) {
    const value = params[key];
    if (value !== undefined && value !== null) {
      httpParams = httpParams.set(key, String(value));
    }
  }
  return httpParams;
}

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  token: string | null = null;
  constructor(
    private http: HttpClient,
    private localStorageService: LocalStorageService
  ) {
    this.token = this.localStorageService.getToken();
  }

  private getHeaders(): HttpHeaders {
    return new HttpHeaders({
      'Content-Type': 'application/json',
      Authorization: this.token ? `Bearer ${this.token}` : ''
    });
  }

  get<T>(url: string, params?: QueryParams): Observable<T> {
    this.token = this.localStorageService.getToken();
    return this.http.get<T>(url, {
      headers: this.getHeaders(),
      params: toHttpParams(params)
    });
  }

  getSimple<T>(url: string, params?: Record<string, any>): Observable<T> {
    return this.http.get<T>(url, {
      headers: this.getHeaders(),
      params: new HttpParams({ fromObject: params || {} })
    });
  }

  getArrayBuffer(url: string): Observable<ArrayBuffer> {
    return this.http.get(url, {
      headers: this.getHeaders(),
      responseType: 'arraybuffer' as const
    });
  }

  getWithoutHeaders<T>(url: string, params?: QueryParams): Observable<T> {
    return this.http.get<T>(url, {
      params: toHttpParams(params)
    });
  }

  post<T>(url: string, body: unknown): Observable<T> {
    return this.http.post<T>(url, body, { headers: this.getHeaders() });
  }

  postWithoutHeaders<T>(url: string, body: unknown): Observable<T> {
    return this.http.post<T>(url, body, {
      headers: new HttpHeaders({
        'Content-Type': 'application/json'
      })
    });
  }

  postWithoutHeaders2<T>(url: string, body: unknown): Observable<T> {
    return this.http.post<T>(url, body, {
      headers: new HttpHeaders({
        Authorization: this.token ? `Bearer ${this.token}` : ''
      })
    });
  }

  put<T>(url: string, body: unknown): Observable<T> {
    return this.http.put<T>(url, body, { headers: this.getHeaders() });
  }

  delete<T>(url: string): Observable<T> {
    return this.http.delete<T>(url, { headers: this.getHeaders() });
  }
}
