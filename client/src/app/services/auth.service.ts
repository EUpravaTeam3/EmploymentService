import {Injectable} from '@angular/core';
import {map} from 'rxjs/operators';
import { HttpHeaders } from '@angular/common/http';
import { User } from '../model/user';
import { ApiService } from './api.service';
import { ConfigService } from './config.service';
import { Router } from '@angular/router';

@Injectable({
  providedIn: 'root'
})

export class AuthService {

    constructor(
        private apiService: ApiService,
        private config: ConfigService,
        private router: Router
      ) {
      }

    postJobAd(jobAd: any, code: any){

      /*var jobAdDTO = {
        jobName: job.jobName,
        ...
      }
      
      /*const postHeaders = new HttpHeaders({
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      });
      return this.apiService.post(this.config.jobad_url, JSON.stringify(jobAdDTO), postHeaders)
      .pipe(map(() => {
      }));*/
    }
	
    logout(){
      localStorage.removeItem("eupravaToken")
      localStorage.removeItem("eupravaRole")

      //this.router.navigate(['...../login-page'])
    }

    /* details about the Logged user are kept in the local storage.
       This method should be implemented on pages that require any Role */
    userIsLoggedIn(): boolean{
      var role = localStorage.getItem("eupravaRole")
      if (role == "" || role == null || role == undefined){
        return false
      }
      return true
    }

    // This method should be implemented on pages that require a certain Role   
    userHasRole(requiredRole: string): boolean{
      var role = localStorage.getItem("eupravaRole")
      if (role == requiredRole){
        return true
      }
      return false
    }
  }
