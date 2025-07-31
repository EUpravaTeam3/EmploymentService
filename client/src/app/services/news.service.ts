import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { News } from '../model/news';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class NewsService {
  private baseUrl = 'http://localhost:8000/news';

  constructor(private http: HttpClient) {}

  getAllNews(): Observable<News[]> {
    return this.http.get<News[]>(`${this.baseUrl}`);
  }

  getNewsById(id: string): Observable<News> {
    return this.http.get<News>(`${this.baseUrl}/${id}`);
  }

  createNews(news: News): Observable<any> {
    return this.http.post(`${this.baseUrl}`, news);
  }

  updateNews(id: string, news: Partial<News>): Observable<any> {
    return this.http.put(`${this.baseUrl}/${id}`, news);
  }

  deleteNews(id: string): Observable<any> {
    return this.http.delete(`${this.baseUrl}/${id}`);
  }
}