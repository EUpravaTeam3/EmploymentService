import { HttpClient } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { ApplicantByUser } from 'src/app/model/applicantByUser';

@Component({
  selector: 'app-applications',
  templateUrl: './applications.component.html',
  styleUrls: ['./applications.component.css']
})
export class ApplicationsComponent implements OnInit {
  applications: ApplicantByUser[] = [];
  searchTerm: string = '';

  constructor(private http: HttpClient) {}

  ngOnInit(): void {
    this.loadApplications();
  }

  loadApplications() {
    
    var ucn = localStorage.getItem("eupravaUcn")

    this.http.get<ApplicantByUser[]>(`http://localhost:8080/applicant/` + ucn)
      .subscribe(data => {
        this.applications = data;
      });
  }

  get filteredApplications(): ApplicantByUser[] {
    if (!this.searchTerm) return this.applications;
    return this.applications.filter(app =>
      app.job_ad_title.toLowerCase().includes(this.searchTerm.toLowerCase())
    );
  }

  onDeleteApplication(applicantId: string | undefined) {
    this.http.delete(`http://localhost:8080/applicant/${applicantId}`)
      .subscribe(() => {
        this.applications = this.applications.filter(app => app.applicant_id !== applicantId);
      });
  }
}
