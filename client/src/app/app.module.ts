import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

import { ToastrModule } from 'ngx-toastr';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { WelcomePageComponent } from './components/welcome-page/welcome-page.component';
import { NavBarComponent } from './components/nav-bar/nav-bar.component';
import { ApiService } from './services/api.service';
import { ConfigService } from './services/config.service';
import { HttpClient, HttpClientModule, HttpHandler } from '@angular/common/http';
import { AuthService } from './services/auth.service';

import { NgxCaptchaModule } from 'ngx-captcha';
import { ProfilePageComponent } from './components/profile-page/profile-page.component';

@NgModule({
  declarations: [
    AppComponent,
    WelcomePageComponent,
    NavBarComponent,
    ProfilePageComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    FormsModule,
    ReactiveFormsModule,
    HttpClientModule,
    NgxCaptchaModule,
    ToastrModule.forRoot({
      positionClass:"toast-top-center",
      preventDuplicates: true,
      closeButton: true
    }),
    BrowserAnimationsModule
  ],
  providers: [ApiService,
    ConfigService,
    AuthService],
  bootstrap: [AppComponent]
})
export class AppModule { }
