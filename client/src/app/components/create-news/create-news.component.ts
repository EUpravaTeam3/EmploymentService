import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { NewsService } from 'src/app/services/news.service';

@Component({
  selector: 'app-create-news',
  templateUrl: './create-news.component.html',
  styleUrls: ['./create-news.component.css']
})
export class CreateNewsComponent {
  companyName = '';
  title = '';
  description = '';

  constructor(private newsService: NewsService, private router: Router){}

  onSubmit() {
    const newNews = {
      employer_id: this.companyName,
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
