import { Component, OnInit } from '@angular/core';
import { ToastrService } from 'ngx-toastr';
import { User } from 'src/app/model/user';
import { AuthService } from 'src/app/services/auth.service';


@Component({
  selector: 'app-profile-page',
  templateUrl: './profile-page.component.html',
  styleUrls: ['./profile-page.component.css']
})
export class ProfilePageComponent implements OnInit {

  constructor(private authService: AuthService,
    private toastr: ToastrService,) { }

  // booleans for the navbar to check the users role and restrict access to pages
  isUserLoggedIn = this.authService.userIsLoggedIn()
  //isUserHEmployer = this.authService.userHasRole("EMPLOYER")
  //isUserCitizen = this.authService.userHasRole("CITIZEN")

  user: User = new User()
  userToEdit: User = new User()

  ngOnInit(): void {

    /*this.authService.getProfile(localStorage.getItem("eupravaId")!).subscribe(data => {
      this.user = data
    }, err => {
      console.log(err)
    })*/
}
}

//reload page after 3 seconds
function reloadTimeOut(){
  setTimeout(() => window.location.reload(), 3000)
}
