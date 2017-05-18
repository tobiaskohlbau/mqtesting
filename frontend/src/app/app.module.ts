import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { MdCardModule, MdIconModule, MdButtonModule } from '@angular/material';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component'
import { MessageService } from './message/message.service';
import { RestModule } from './rest/rest.module';


import { MessageComponent } from './message/message.component';

import { AuthGuard } from './auth/auth-guard.service';
import { AuthService } from './auth/auth.service';
import { PostAuthService } from './auth/post-auth.service';
import { LoginComponent } from './login/login.component';

@NgModule({
  declarations: [
    AppComponent,
    MessageComponent,
    LoginComponent
  ],
  imports: [
    BrowserModule,
    FormsModule,
    HttpModule,
    AppRoutingModule,
    BrowserAnimationsModule,
    MdCardModule,
    MdIconModule,
    MdButtonModule,
    RestModule
  ],
  providers: [
    MessageService,
    AuthService,
    AuthGuard,
    PostAuthService
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }