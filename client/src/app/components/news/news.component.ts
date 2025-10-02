import { Component, OnInit } from '@angular/core';
import { NewsService } from 'src/app/services/news.service';
import { News } from 'src/app/model/news';
import { Router } from '@angular/router';
import { HttpClient } from '@angular/common/http';
import { CheckedUser } from 'src/app/model/roleResponse';
import { CompanyService } from 'src/app/services/company.service';
import { Company } from 'src/app/model/company';

@Component({
  selector: 'app-news',
  templateUrl: './news.component.html',
  styleUrls: ['./news.component.css']
})

export class NewsComponent implements OnInit {
  news: News[] = [];
  newNews: News = {
    employer_id: '',
    title: '',
    description: ''
  };
  searchTerm: string = '';
  role: string = ''
  company!: Company;

  constructor(private newsService: NewsService, private router: Router, private http: HttpClient,
    private companyService: CompanyService) {}

  ngOnInit(): void {

    this.company._id = ""

    var ucn = localStorage.getItem("eupravaUcn")

      this.companyService.getCompanyByOwner(ucn!).subscribe({
        next: (data) => (this.company = data),
        error: (err) => (console.log(err))
      });

                this.http.get<CheckedUser>(`http://localhost:9090/user/employment`, {withCredentials: true})
                  .subscribe(data => {
                      this.role = data.role
                  }, err => {
                    window.location.href = "http://localhost:4200/login"
                  });
    this.loadNews();
  }

  loadNews(): void {
    this.newsService.getAllNews().subscribe(data => {
      this.news = data
    });
  }

  onDeleteNews(id: string): void {
    this.newsService.deleteNews(id).subscribe(() => {
      this.news = this.news.filter(news => news._id !== id);
    });
  }

  onCreateNews() {
    this.router.navigateByUrl("/create-news")
  }

      get filteredNews() {
    if (!this.searchTerm) return this.news;
    const term = this.searchTerm.toLowerCase();
    return this.news.filter(newsObj =>
      newsObj.title.toLowerCase().includes(this.searchTerm) ||
      newsObj.description.toLowerCase().includes(this.searchTerm)
    );
  }
}