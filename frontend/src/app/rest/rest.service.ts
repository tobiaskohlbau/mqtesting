import { Injectable } from '@angular/core';
import { Http, RequestOptions, Headers } from '@angular/http';
import { Router } from '@angular/router';

import { Observable } from 'rxjs/Observable';

import { JsonInterface } from './json.interface';

import { AuthService } from '../auth/auth.service';

@Injectable()
export class RestService {

  private url: string
  private headers: Headers;

  constructor(private http: Http, private authService: AuthService, private router: Router) {
    this.headers = new Headers();
    this.headers.append('Content-Type', 'application/json');
  }

  setURL(url: string): void {
    this.url = url;
  }

  create(model: JsonInterface): Observable<any> {
    return this.http.post(this.url, model.json(), this.authService.requestOptions()).map(res => res.json()).catch(this.handleError);
  }

  read(id?: string): Observable<any> {
    if (id) {
      return this.http.get(this.url + '/' + id, this.authService.requestOptions()).map(res => res.json()).catch(this.handleError);
    }

    return this.http.get(this.url, this.authService.requestOptions()).map(res => res.json()).catch(this.handleError);
  }

  update(model: JsonInterface): Observable<any> {
    return this.http.put(this.url + '/' + model.id, model.json(), this.authService.requestOptions()).map(res => res.json()).catch(this.handleError);
  }

  delete(id?: string): Observable<any> {
    if (id) {
      return this.http.delete(this.url + '/' + id, this.authService.requestOptions()).catch(this.handleError);
    }
    return this.http.get(this.url, this.authService.requestOptions()).map(res => res.json()).catch(this.handleError);
  }


  private handleError(error: Response) {
    if (error.status == 401) {
      console.log('authorization required');
    }
    return Observable.throw('Internal server error');
  }
}
