import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { JobAdService } from 'src/app/services/jobad.service';
import { Router } from '@angular/router';
import { JobAd } from 'src/app/model/jobad';
import { Job } from 'src/app/model/job';
import { HttpClient } from '@angular/common/http';
import { CheckedUser } from 'src/app/model/roleResponse';
import { Company } from 'src/app/model/company';

@Component({
  selector: 'app-create-jobad',
  templateUrl: './create-jobad.component.html',
  styleUrls: ['./create-jobad.component.css']
})

export class CreateJobAdComponent implements OnInit {

  ad_title = '';
  job_description = '';
  qualification = '';
  job_type = '';
  jobs: Job[] = []
  selectedJobId = ''

  constructor(private router: Router, private jobAdService: JobAdService, private http: HttpClient) {}

  ngOnInit(): void {
        this.http.get<CheckedUser>(`http://localhost:9090/user/employment`, {withCredentials: true})
          .subscribe(data => {
            console.log(data)
            if (data.role != "employer") {
              this.router.navigate(["/"])
            }
          });

          var ucn = localStorage.getItem("eupravaUcn")

          this.http.get<Company>(`http://localhost:8000/company/owner/` + ucn)
          .subscribe(company => {
            console.log(company)
            this.http.get<Job[]>("http://localhost:8000/job/company/" + company._id)
            .subscribe(data => {
              this.jobs = data
            })
          })
  }

  onSubmit(): void {
      const newJob: JobAd = {
      ad_title: this.ad_title,
      job_description: this.job_description,
      qualification: this.qualification,
      job_type: this.job_type,
      job_id: this.selectedJobId
    };
      this.jobAdService.postJobAd(newJob).subscribe(() => {
        this.router.navigate(['/jobads']); // Navigate back to list
      });
  }
}
