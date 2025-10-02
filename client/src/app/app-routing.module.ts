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
import { CreateNewsComponent } from './components/create-news/create-news.component';
import { CompanyComponent } from './components/company/company.component';
import { CvComponent } from './components/cv/cv.component';
import { ApplicationsComponent } from './components/applications/applications.component';
import { ApplicantsComponent } from './components/applicants/applicants.component';
import { JobsComponent } from './components/jobs/jobs.component';
import { EmployeeComponent } from './components/employee/employee.component';

const routes: Routes = [
  {path: '', pathMatch:'full', component: WelcomePageComponent},
  {path: 'welcome-page', component: WelcomePageComponent},
  {path: 'profile-page', component: ProfilePageComponent},
  {path: 'jobad', component: JobadComponent},
  {path: 'jobads/create', component: CreateJobAdComponent},
  {path: 'news', component: NewsComponent},
  {path: 'sign-in', component: LoginComponent},
  {path: 'create-news', component: CreateNewsComponent},
  {path: 'company', component: CompanyComponent},
  {path: 'resume', component: CvComponent},
  {path: 'applications', component: ApplicationsComponent},
  {path: 'applicants', component: ApplicantsComponent},
  {path: 'jobs', component: JobsComponent},
  {path: 'employee', component: EmployeeComponent}
]
;

@NgModule({
  imports: [CommonModule, RouterModule.forRoot(routes,{scrollPositionRestoration: 'enabled'})],
  exports: [RouterModule]
})
export class AppRoutingModule { }
