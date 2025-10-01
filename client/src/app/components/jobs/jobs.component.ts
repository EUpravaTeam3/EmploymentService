import { HttpClient } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { Company } from 'src/app/model/company';
import { Job } from 'src/app/model/job';
import { CheckedUser } from 'src/app/model/roleResponse';

@Component({
  selector: 'app-jobs',
  templateUrl: './jobs.component.html',
  styleUrls: ['./jobs.component.css']
})
export class JobsComponent implements OnInit{

  constructor(private router: Router, private http: HttpClient) {}

  jobs: Job[] = []
  newJob: Job = {
    position_name: "",
    pay: 0,
    employee_capacity: 0,
    num_of_employees: 0,
    company_id: ""
  }
  companyId: string | undefined = ""
  creating = false;

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
                this.companyId = company._id
                this.http.get<Job[]>("http://localhost:8000/job/company/" + company._id)
                .subscribe(data => {
                  this.jobs = data
                })
              })
  }

    onCreateJob() {
    if (!this.newJob.position_name || this.newJob.pay == null || this.newJob.employee_capacity == null) return;

    this.creating = true;
    const payload: Job = {
      position_name: this.newJob.position_name.trim(),
      pay: this.newJob.pay,
      employee_capacity: this.newJob.employee_capacity,
      company_id: this.companyId!,
      num_of_employees: 0
    };

    this.http.post<Job>(`http://localhost:8000/job`, payload)
      .subscribe({
        next: created => {
          this.jobs.unshift(created || payload);
          this.newJob = { position_name: '', pay: 0, employee_capacity: 1, company_id: '', num_of_employees: 0 };
          this.creating = false;
        },
        error: err => {
          console.error('Failed to create job', err);
          this.creating = false;
        }
      });
  }


  onDeleteJob(jobId: string | undefined, index: number) {
    if (!jobId) {
      this.jobs.splice(index, 1);
      return;
    }

    const removed = this.jobs.splice(index, 1)[0];

    this.http.delete(`http://localhost:8000/job/${jobId}`)
      .subscribe({
        next: () => {
        },
        error: err => {
          console.error('Failed to delete job', err);
          this.jobs.splice(index, 0, removed);
        }
      });
  }
}
