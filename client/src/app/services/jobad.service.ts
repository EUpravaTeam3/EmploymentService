import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { JobAd } from '../model/jobad';
import { JobAdDTO } from '../model/jobadDTO';
import { Applicant } from '../model/applicant';

@Injectable({
  providedIn: 'root'
})
export class JobAdService {
  private apiUrl = 'http://localhost:8000/jobad';

  constructor(private http: HttpClient) {}

  getJobAds(): Observable<JobAdDTO[]> {
    return this.http.get<JobAdDTO[]>(this.apiUrl);
  }

  getJobAdById(id: string): Observable<JobAd> {
    return this.http.get<JobAd>(`${this.apiUrl}/${id}`);
  }

  postJobAd(jobAd: JobAd): Observable<any> {
    return this.http.post(this.apiUrl, jobAd);
  }

  updateJobAd(id: string, jobAd: JobAd): Observable<any> {
    return this.http.put(`${this.apiUrl}/${id}`, jobAd);
  }

  deleteJobAd(id: string): Observable<any> {
    return this.http.delete(`${this.apiUrl}/${id}`);
  }
}