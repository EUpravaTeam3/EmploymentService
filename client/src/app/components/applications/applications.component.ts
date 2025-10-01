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
  ucn: string = '';

  constructor(private http: HttpClient) {}

  ngOnInit(): void {
    this.loadApplications();
  }

  loadApplications() {
    
    this.ucn = localStorage.getItem("eupravaUcn")!

    this.http.get<ApplicantByUser[]>(`http://localhost:8000/applicant/` + this.ucn)
      .subscribe(data => {
        this.applications = data;
        console.log(data)
      }, err => {console.log(err)});
  }

  get filteredApplications(): ApplicantByUser[] {
    if (!this.searchTerm) return this.applications;
    return this.applications.filter(app =>
      app.job_ad_title.toLowerCase().includes(this.searchTerm.toLowerCase())
    );
  }

  onDeleteApplication(applicantId: string | undefined) {
    this.http.delete(`http://localhost:8000/applicant/` + applicantId)
      .subscribe(() => {
        window.location.reload()
      }, err => {
        console.log(err)
      });
  }
}
