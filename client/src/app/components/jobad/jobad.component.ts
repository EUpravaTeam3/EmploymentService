import { Component, OnInit } from '@angular/core';
import { JobAdService } from 'src/app/services/jobad.service';
import { JobAd } from 'src/app/model/jobad';
import { Router } from '@angular/router';
import { JobAdDTO } from 'src/app/model/jobadDTO';

@Component({
  selector: 'app-jobad',
  templateUrl: './jobad.component.html',
  styleUrls: ['./jobad.component.css']
})
export class JobadComponent implements OnInit {
  jobAds: JobAdDTO[] = [];

  searchTerm: string = '';

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
    this.router.navigate(['/jobads']);
  }

      seeCompany(company: string): void {
    this.router.navigate(['/company/' + company]);
  }

    get filteredJobAds() {
    if (!this.searchTerm) return this.jobAds;
    const term = this.searchTerm.toLowerCase();
    return this.jobAds.filter(job =>
      job.ad_title.toLowerCase().includes(term) ||
      job.job_description.toLowerCase().includes(term) ||
      job.company_name.toLowerCase().includes(term) ||
      job.qualification.toLowerCase().includes(term) ||
      job.job_type.toLowerCase().includes(term)
    );
  }
}
