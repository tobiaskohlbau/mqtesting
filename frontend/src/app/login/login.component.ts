import { Component, OnInit }   from '@angular/core';
import { Router }      from '@angular/router';

import { AuthService } from '../auth/auth.service';

@Component({
  selector: 'login-component',
  templateUrl: 'login.component.html'
})
export class LoginComponent implements OnInit {
  provider: string[];

  constructor(public authService: AuthService, public router: Router) {
  }

  ngOnInit() {
    this.authService.provider().subscribe(provider => this.provider = provider);
  }

  login(provider: string) {
    this.authService.login(provider);
  }

  logout() {
    this.authService.logout();
  }
}