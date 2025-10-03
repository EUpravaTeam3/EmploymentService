import { HttpClient } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { ReviewOfCompany } from 'src/app/model/reviewOfCompany';

@Component({
  selector: 'app-employee',
  templateUrl: './employee.component.html',
  styleUrls: ['./employee.component.css']
})
export class EmployeeComponent implements OnInit {

  constructor(private http: HttpClient, private router: Router){}

  employee?: ReceivedEmployee
  name: string = ''
  review: ReviewOfCompany = {
    employee_id : '',
    employer_id: '000000000000000000000000',
    rating: 1,
    description: ''
  }

  ngOnInit(): void {

    var ucn = localStorage.getItem("eupravaUcn")
    this.name = localStorage.getItem("eupravaName") + " " + localStorage.getItem("eupravaSurname")
    
    this.http.get<ReceivedEmployee[]>("http://localhost:8000/employee/" + ucn)
    .subscribe(res => { console.log(res); this.employee = res[0]}
  ,err => {console.log(err)})
  }

  onQuit(){
    const result = confirm("Are you sure you want to quit?");

    if (result) {

      console.log(this.employee)
  
      this.http.post("http://localhost:8000/employee/quit", {"ucn": this.employee?.citizen_ucn},
         { withCredentials: true }).subscribe(res => window.location.reload(), err => console.log(err))
    }
  }

  onSubmitReview() {

  this.review.employee_id = this.employee?._id!
  console.log("Submitting review:", this.review);

  this.http.post("http://localhost:8000/company/review", this.review, { withCredentials: true }).subscribe(
    res => {
      console.log("Review submitted", res);
      alert("Review submitted successfully!");
      this.review = { rating: 0, description: '', employee_id: this.employee?._id!,
        employer_id: '000000000000000000000000'
       }; // reset form
    },
    err => console.error("Error submitting review", err)
  );
}
}

export interface ReceivedEmployee {
  _id: string;
  citizen_ucn: string;
  start_date: string;      
  end_date?: string;       
  employer_review?: string;
  position_name: string;
  pay: number;
  company_name: string;
}