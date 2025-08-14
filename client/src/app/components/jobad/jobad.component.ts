import { Component, OnInit } from '@angular/core';
import { JobAdService } from 'src/app/services/jobad.service';
import { JobAd } from 'src/app/model/jobad';
import { Router } from '@angular/router';

@Component({
  selector: 'app-jobad',
  templateUrl: './jobad.component.html',
  styleUrls: ['./jobad.component.css']
})
export class JobadComponent implements OnInit {
  jobAds: JobAd[] = [];

  constructor(private jobAdService: JobAdService, private router: Router) {}

  ngOnInit(): void {
    this.loadJobAds();
  }

  loadJobAds(): void {
    this.jobAdService.getJobAds().subscribe((data) => {
      this.jobAds = data;
    });
  }

  deleteJobAd(id: string): void {
    if (confirm('Are you sure you want to delete this job ad?')) {
      this.jobAdService.deleteJobAd(id).subscribe(() => {
        this.loadJobAds();
      });
    }
  }

    onCreateJobAd() {
    this.router.navigate(['/jobads/create']);
  }
}
