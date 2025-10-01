import { HttpClient } from '@angular/common/http';
import { Component } from '@angular/core';

@Component({
  selector: 'app-applicants',
  templateUrl: './applicants.component.html',
  styleUrls: ['./applicants.component.css']
})
export class ApplicantsComponent {
  applications: ApplicantByCompany[] = [];
  searchTerm: string = '';

  constructor(private http: HttpClient) {}

  ngOnInit(): void {
    this.loadApplications();
  }

  loadApplications() {
    
    var ucn = localStorage.getItem("eupravaUcn")

    this.http.get<ApplicantByCompany[]>(`http://localhost:8000/applicant/company/` + ucn)
      .subscribe(data => {
        console.log(data)
        this.applications = data;
      }, err => {
        console.log(err)
      });
  }

  onAcceptApplication(app: ApplicantByCompany){
    this.http.post(`http://localhost:8000/employee`, app, { withCredentials: true   }).subscribe(res =>
      window.location.reload(), err => console.log(err)
    )
}
}

export interface ApplicantByCompany {
    "position_name": string,
    "ad_title": string,
    "job_ad_id": string,
    "citizen_ucn": string,
    "name": string,
    "email": string,
    "education": string[],
    "work_experience": string[],
    "description": string
  }