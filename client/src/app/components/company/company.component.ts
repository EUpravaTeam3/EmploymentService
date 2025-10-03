import { HttpClient } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Company } from 'src/app/model/company';
import { ReviewOfCompany } from 'src/app/model/reviewOfCompany';
import { CompanyService } from 'src/app/services/company.service';

@Component({
  selector: 'app-company',
  templateUrl: './company.component.html',
  styleUrls: ['./company.component.css']
})
export class CompanyComponent implements OnInit{

    company?: Company | null;
    applicants = []
    reviews: ReviewOfCompany[] = []

  constructor(
    private route: ActivatedRoute,
    private companyService: CompanyService,
    private http: HttpClient
  ) {}

  ngOnInit(){
    
    var companyId = localStorage.getItem("companyId")

    if (companyId){

      localStorage.removeItem("companyId")
      this.companyService.getCompanyById(companyId).subscribe({
        next: (data) => {(this.company = data);
            this.http.get<ReviewOfCompany[]>("http://localhost:8000/company/reviews/" + companyId)
      .subscribe({
        next: (res) => {
          console.log("Reviews received:", res);
          this.reviews = res;
        },
        error: (err) => {
          console.error("Error fetching reviews:", err);
        }
      });}
        ,
        error: (err) => console.error(err)
      });

    } else {
       var ucn = localStorage.getItem("eupravaUcn")
      if (ucn){
        this.companyService.getCompanyByOwner(ucn).subscribe({
        next: (data) => {(this.company = data)
          this.http.get<ReviewOfCompany[]>("http://localhost:8000/company/reviews/" + this.company._id)
      .subscribe({
        next: (res) => {
          console.log("Reviews received:", res);
          this.reviews = res;
        },
        error: (err) => {
          console.error("Error fetching reviews:", err);
        }
      });
        },
        error: (err) => console.error(err)
      });
      }
    }
  }
}
