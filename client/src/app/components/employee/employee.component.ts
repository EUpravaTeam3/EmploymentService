import { HttpClient } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';

@Component({
  selector: 'app-employee',
  templateUrl: './employee.component.html',
  styleUrls: ['./employee.component.css']
})
export class EmployeeComponent implements OnInit {

  constructor(private http: HttpClient, private router: Router){}

  employee?: ReceivedEmployee
  name: string = ''

  ngOnInit(): void {

    var ucn = localStorage.getItem("eupravaUcn")
    this.name = localStorage.getItem("eupravaName") + " " + localStorage.getItem("eupravaSurname")
    
    this.http.get<ReceivedEmployee[]>("http://localhost:8000/employee/" + ucn)
    .subscribe(res => { console.log(res); this.employee = res[0]}
  ,err => {console.log(err)})
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