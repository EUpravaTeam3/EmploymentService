import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { JobAdService } from 'src/app/services/jobad.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-create-jobad',
  templateUrl: './create-jobad.component.html'
})
export class CreateJobAdComponent {
  jobAdForm: FormGroup;

  constructor(
    private fb: FormBuilder,
    private jobAdService: JobAdService,
    private router: Router
  ) {
    this.jobAdForm = this.fb.group({
      ad_title: ['', Validators.required],
      job_description: ['', Validators.required],
      qualification: ['', Validators.required],
      job_type: ['', Validators.required]
    });
  }
  ngOnInit(): void {
    throw new Error('Method not implemented.');
  }

  onSubmit(): void {
    if (this.jobAdForm.valid) {
      this.jobAdService.postJobAd(this.jobAdForm.value).subscribe(() => {
        this.router.navigate(['/jobads']); // Navigate back to list
      });
    }
  }
}
