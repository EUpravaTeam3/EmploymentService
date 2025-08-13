import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { CommonModule } from '@angular/common';
import { AppComponent } from './app.component';
import { WelcomePageComponent } from './components/welcome-page/welcome-page.component';
import { ProfilePageComponent } from './components/profile-page/profile-page.component';
import { CreateJobAdComponent } from './components/create-jobad/create-jobad.component';
import { JobadComponent } from './components/jobad/jobad.component';
import { NewsComponent } from './components/news/news.component';
import { LoginComponent } from './components/login/login.component';

const routes: Routes = [
  {path: '', pathMatch:'full', component: WelcomePageComponent},
  {path: 'welcome-page', component: WelcomePageComponent},
  {path: 'profile-page', component: ProfilePageComponent},
  {path: 'jobad', component: JobadComponent},
  {path: 'jobads/create', component: CreateJobAdComponent},
  {path: 'news', component: NewsComponent},
  {path: 'sign-in', component: LoginComponent}]
;

@NgModule({
  imports: [CommonModule, RouterModule.forRoot(routes,{scrollPositionRestoration: 'enabled'})],
  exports: [RouterModule]
})
export class AppRoutingModule { }
