import { Injectable } from '@angular/core';
import { Http, RequestOptions, Headers } from '@angular/http';

import { Observable } from 'rxjs/Observable';

import { JsonInterface } from './json.interface';

@Injectable()
export class RestService {

  private url: string
  private headers: Headers;
  private requestOptions: RequestOptions;

  constructor(public http: Http) {
    this.headers = new Headers();
    this.headers.append('Content-Type', 'application/json');
    this.requestOptions = new RequestOptions({ headers: this.headers });
  }

  setURL(url: string): void {
    this.url = url;
  }

  create(model: JsonInterface): Observable<any> {
    return this.http.post(this.url, model.json(), this.requestOptions).map(res => res.json());
  }

  read(id?: string): Observable<any> {
    if (id) {
      return this.http.get(this.url + '/' + id, this.requestOptions).map(res => res.json());
    }

    return this.http.get(this.url, this.requestOptions).map(res => res.json());
  }

  update(model: JsonInterface): Observable<any> {
    return this.http.put(this.url + '/' + model.id, model.json(), this.requestOptions).map(res => res.json());
  }

  delete(id?: string): Observable<any> {
    if (id) {
      return this.http.delete(this.url + '/' + id, this.requestOptions);
    }
    return this.http.get(this.url, this.requestOptions).map(res => res.json());
  }
}
