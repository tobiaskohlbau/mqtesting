import { Injectable } from '@angular/core';
import { Http, RequestOptions, Headers } from '@angular/http';
import { Router } from '@angular/router';

import { Observable } from 'rxjs/Observable';
import { Subject } from 'rxjs/Subject';
import 'rxjs/add/observable/of';
import 'rxjs/add/operator/do';
import 'rxjs/add/operator/delay';

import { PostAuthService } from './post-auth.service';

@Injectable()
export class AuthService {
  user: JSON;
  isAuthenticated: boolean = false;

  constructor(private postAuthService: PostAuthService, private http: Http) {
    this.postAuthService.addPostAuthListener((event) => this.onPostAuth(event as MessageEvent));
    // this.check().subscribe(res => this.isAuthenticated = res);
  }

  get token(): string {
    return localStorage.getItem("mqtesting_jwt_token");
  }

  login(provider: string) {
    window.open('https://localhost/oauth/authenticate/' + provider, '', 'width=800,width=600');
  }

  logout() {
    localStorage.removeItem("mqtesting_jwt_token");
    this.isAuthenticated = false;
  }

  onPostAuth(event: MessageEvent) {
    if (event.data.status === 200) {
      localStorage.setItem("mqtesting_jwt_token", event.data.jwt);
      this.check().subscribe(res => this.isAuthenticated = res);
    }
  }

  check(): Observable<boolean> {
    return this.http.get("https://localhost/api/user", this.requestOptions()).map(res => {
      this.user = res.json();
      return res.status === 200;
    });
  }

  provider(): Observable<string[]> {
    return this.http.get("https://localhost/api/provider").map(res => res.json());
  }

  requestOptions(): RequestOptions {
    let headers = new Headers({'Authorization': 'Bearer ' + this.token });
    let options = new RequestOptions({ headers: headers});
    return options;
  }
}