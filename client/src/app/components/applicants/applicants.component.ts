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

    this.http.get<ApplicantByCompany[]>(`http://localhost:8080/applicant/company/` + ucn)
      .subscribe(data => {
        this.applications = data;
      });
  }

  onAcceptApplication(app: ApplicantByCompany){
    this.http.post(`http://localhost:8080/applicant/employ`, app).subscribe(res =>
      window.location.reload(), err => alert(err)
    )
}
}

export interface ApplicantByCompany {
    "position_name": string,
    "ad_title": string,
    "citizen_ucn": string,
    "name": string,
    "email": string,
    "education": string[],
    "work_experience": string[],
    "description": string
  }