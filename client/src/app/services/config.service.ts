import {Injectable} from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class ConfigService {

  private _api_url = 'http://localhost:8000';
  private _jobad_url = this._api_url + '/jobad';

  get jobad_url(): string {
    return this._jobad_url;
  }
}