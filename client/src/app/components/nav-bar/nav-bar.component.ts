import { HttpClient } from '@angular/common/http';
import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { ToastrService } from 'ngx-toastr';
import { AuthService } from 'src/app/services/auth.service';

@Component({
  selector: 'nav-bar',
  templateUrl: './nav-bar.component.html',
  styleUrls: ['./nav-bar.component.css']
})
export class NavBarComponent {

  constructor(private authService: AuthService,
    private toastr: ToastrService, private http: HttpClient,
  private router: Router){}

  logout(){
      this.http.post('http://localhost:9090/user/logout', {
    }, { withCredentials: true }).subscribe(res => {
      console.log('status:', res);
      localStorage.removeItem("eupravaUcn")
      localStorage.removeItem("eupravaEmail")
      localStorage.removeItem("eupravaName")
      localStorage.removeItem("eupravaSurname")
      localStorage.removeItem("eupravaAddress")
      
      this.router.navigateByUrl("/sign-in")
    });
  }

  //username = localStorage.getItem("eupravaUsername")

}
