import { Component, OnInit } from '@angular/core';
import { JobAdService } from 'src/app/services/jobad.service';
import { JobAd } from 'src/app/model/jobad';
import { Router } from '@angular/router';
import { JobAdDTO } from 'src/app/model/jobadDTO';
import { Applicant } from 'src/app/model/applicant';
import { HttpClient } from '@angular/common/http';
import { CompanyService } from 'src/app/services/company.service';
import { Company } from 'src/app/model/company';

@Component({
  selector: 'app-jobad',
  templateUrl: './jobad.component.html',
  styleUrls: ['./jobad.component.css']
})
export class JobadComponent implements OnInit {
  jobAds: JobAdDTO[] = [];

  searchTerm: string = '';
  company: Company = {
    company_name: '',
    status: '',
    registration_date: '',
    id_number: 0,
    tax_id_number: 0,
    owner_ucn: '',
    address: undefined,
    work_field: undefined
  }
  role: string = ''

  constructor(private jobAdService: JobAdService, private router: Router, private http: HttpClient,
    private companyService: CompanyService) {}

  ngOnInit(): void {
    var ownerUcn = localStorage.getItem("eupravaUcn")!
    this.role = localStorage.getItem("eupravaRole")!

    this.companyService.getCompanyByOwner(ownerUcn).subscribe({
        next: (data) => (this.company = data),
        error: (err) => console.error(err)
      });

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
    this.router.navigate(['/jobads/create'])
      .then(ok => console.log('nav success', ok))
  .catch(err => console.error('nav error', err));;
  }

      seeCompany(company: string): void {
    this.router.navigate(['/company/' + company]);
  }

    get filteredJobAds() {
      if (this.jobAds.length > 0){
    if (!this.searchTerm) return this.jobAds;
    const term = this.searchTerm.toLowerCase();
    return this.jobAds.filter(job =>
      job.ad_title.toLowerCase().includes(term) ||
      job.job_description.toLowerCase().includes(term) ||
      job.company_name.toLowerCase().includes(term) ||
      job.qualification.toLowerCase().includes(term) ||
      job.job_type.toLowerCase().includes(term)
    );
  } return null
  }

  applyForJobAd(jobAdId: string) {
    const storedUserUcn = localStorage.getItem("eupravaUcn")
    if (storedUserUcn) {
      const userUcn = JSON.parse(storedUserUcn);
      var application: Applicant = {
        job_ad_id: jobAdId
      }
      this.http.post("http://localhost:8000/applicant/" + storedUserUcn, application)
      .subscribe(data => {console.log(data)
        window.location.reload()
      }, err => {console.log(err); alert(err.Error())})
    }
  }
}
