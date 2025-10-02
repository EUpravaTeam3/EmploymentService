import { Component } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Company } from 'src/app/model/company';
import { CompanyService } from 'src/app/services/company.service';

@Component({
  selector: 'app-company',
  templateUrl: './company.component.html',
  styleUrls: ['./company.component.css']
})
export class CompanyComponent {

    company?: Company | null;
    applicants = []

  constructor(
    private route: ActivatedRoute,
    private companyService: CompanyService
  ) {}

  ngOnInit(): void {
    
      var ucn = localStorage.getItem("eupravaUcn")
      if (ucn){
        this.companyService.getCompanyByOwner(ucn).subscribe({
        next: (data) => (this.company = data),
        error: (err) => console.error(err)
      });
      }
    
  }
}
