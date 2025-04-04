import { Component } from '@angular/core';
import { ToastrService } from 'ngx-toastr';
import { AuthService } from 'src/app/services/auth.service';

@Component({
  selector: 'nav-bar',
  templateUrl: './nav-bar.component.html',
  styleUrls: ['./nav-bar.component.css']
})
export class NavBarComponent {

  constructor(private authService: AuthService,
    private toastr: ToastrService){}

  // booleans for the navbar to check the users role and restrict access to pages
  //isUserLoggedIn = this.authService.userIsLoggedIn()
  //isUserEmployer = this.authService.userHasRole("EMPLOYER")
  //isUserCitizen = this.authService.userHasRole("CITIZEN")

  logout(){
    //this.authService.logout()
    //this.toastr.success("Successfully Loged out!", "Success");
  }

  //username = localStorage.getItem("eupravaUsername")

}
