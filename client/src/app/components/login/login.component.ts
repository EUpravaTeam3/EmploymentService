import { HttpClient } from '@angular/common/http';
import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { UserDTO } from 'src/app/model/userDTO';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent {
    username = '';
  password = '';

  constructor(private http: HttpClient, private router: Router) {}

  onLogin() {
      this.http.post('http://localhost:9090/user', {
      ucn: this.username,
      password: this.password
    }, { withCredentials: true }).subscribe(res => {
      console.log('Logged in:', res);
      var user: any
      user = res
      localStorage.setItem("eupravaUcn", user.ucn)
      localStorage.setItem("eupravaEmail", user.email)
      localStorage.setItem("eupravaName", user.name)
      localStorage.setItem("eupravaSurname", user.surname)
      localStorage.setItem("eupravaAddress", user.address)
      this.router.navigateByUrl("/")
    });
  }
  }


