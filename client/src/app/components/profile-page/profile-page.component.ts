import { HttpClient } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { ToastrService } from 'ngx-toastr';
import { CheckedUser } from 'src/app/model/roleResponse';
import { User } from 'src/app/model/user';
import { AuthService } from 'src/app/services/auth.service';


@Component({
  selector: 'app-profile-page',
  templateUrl: './profile-page.component.html',
  styleUrls: ['./profile-page.component.css']
})
export class ProfilePageComponent implements OnInit {

  constructor(private authService: AuthService,
    private toastr: ToastrService,
  private http: HttpClient) { }

  name: string = ''
  ucn: string = ''
  email: string = ''
  isUserLoggedIn = this.authService.userIsLoggedIn()
  user: User = new User()
  userToEdit: User = new User()

  ngOnInit(): void {

    this.http.get<CheckedUser>(`http://localhost:9090/user/employment`, {withCredentials: true})
          .subscribe(data => {
            console.log(data)
            this.name = data.name + " " + data.surname
            this.ucn = data.ucn
            this.email = data.email
      }, err => {
          window.location.href = "http://localhost:4200/login"
        });
  }
}

