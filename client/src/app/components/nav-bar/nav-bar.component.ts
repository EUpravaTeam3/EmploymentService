import { HttpClient } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { ToastrService } from 'ngx-toastr';
import { CheckedUser } from 'src/app/model/roleResponse';
import { AuthService } from 'src/app/services/auth.service';

@Component({
  selector: 'nav-bar',
  templateUrl: './nav-bar.component.html',
  styleUrls: ['./nav-bar.component.css']
})
export class NavBarComponent implements OnInit {

  constructor(private authService: AuthService,
    private toastr: ToastrService, private http: HttpClient,
  private router: Router){}

  role: string = ""

  ngOnInit(): void {

    //this.role = localStorage.getItem("eupravaUcn")!
            this.http.get<CheckedUser>(`http://localhost:9090/user/employment`, {withCredentials: true})
              .subscribe(data => {
                console.log(data)
                  localStorage.setItem("eupravaUcn", data.ucn)
                  localStorage.setItem("eupravaEmail", data.email)
                  localStorage.setItem("eupravaName", data.name)
                  localStorage.setItem("eupravaSurname", data.surname)
                  localStorage.setItem("eupravaRole", data.role)
                  this.role = data.role
              }, err => {
                window.location.href = "http://localhost:4200/login"
              });
  }

  logout(){
      this.http.post('http://localhost:9090/user/logout', {
    }, { withCredentials: true }).subscribe(res => {
      console.log('status:', res);
      localStorage.removeItem("eupravaUcn")
      localStorage.removeItem("eupravaEmail")
      localStorage.removeItem("eupravaName")
      localStorage.removeItem("eupravaSurname")
      localStorage.removeItem("eupravaRole")
      
      window.location.href = "http://localhost:4200/login"
    });
  }

}
