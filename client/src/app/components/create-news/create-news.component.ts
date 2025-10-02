import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { Company } from 'src/app/model/company';
import { CompanyService } from 'src/app/services/company.service';
import { NewsService } from 'src/app/services/news.service';

@Component({
  selector: 'app-create-news',
  templateUrl: './create-news.component.html',
  styleUrls: ['./create-news.component.css']
})
export class CreateNewsComponent implements OnInit{
  companyName = '';
  title = '';
  description = '';
  role = '';
  company?: Company;

  constructor(private newsService: NewsService, private router: Router, private companyService: CompanyService){}

  ngOnInit(): void {
    var ucn = localStorage.getItem("eupravaUcn")
      if (ucn){
        this.companyService.getCompanyByOwner(ucn).subscribe({
        next: (data) => (this.company = data),
        error: (err) => console.error(err)
      });
      }
  }

  onSubmit() {
    const newNews = {
      employer_id: this.company?._id!,
      title: this.title,
      description: this.description
    };

        if (!newNews.title || !newNews.description) {
      alert('All fields are required.');
      return;
    }

    this.newsService.createNews(newNews).subscribe({
      next: () => {
        this.router.navigateByUrl("/news")
      },
      error: err => console.error('Error creating news', err)
    });
  }
}
