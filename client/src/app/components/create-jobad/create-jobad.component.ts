import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { JobAdService } from 'src/app/services/jobad.service';
import { Router } from '@angular/router';
import { JobAd } from 'src/app/model/jobad';

@Component({
  selector: 'app-create-jobad',
  templateUrl: './create-jobad.component.html',
  styleUrls: ['./create-jobad.component.css']
})
export class CreateJobAdComponent {

  ad_title = '';
  job_description = '';
  qualification = '';
  job_type = '';

  constructor(private router: Router, private jobAdService: JobAdService) {}


  onSubmit(): void {
      const newJob: JobAd = {
      ad_title: this.ad_title,
      job_description: this.job_description,
      qualification: this.qualification,
      job_type: this.job_type,
      job_id: ''
    };
      this.jobAdService.postJobAd(newJob).subscribe(() => {
        this.router.navigate(['/jobads']); // Navigate back to list
      });
  }
}
