import { Component, OnInit } from '@angular/core';
import { NewsService } from 'src/app/services/news.service';
import { News } from 'src/app/model/news';

@Component({
  selector: 'app-news',
  templateUrl: './news.component.html'
})
export class NewsComponent implements OnInit {
  newsList: News[] = [];
  newNews: News = {
    employer_id: '',
    title: '',
    description: ''
  };

  constructor(private newsService: NewsService) {}

  ngOnInit(): void {
    this.loadNews();
  }

  loadNews(): void {
    this.newsService.getAllNews().subscribe(data => {
      this.newsList = data
    });
  }

  createNews(): void {
    if (!this.newNews.employer_id || !this.newNews.title || !this.newNews.description) {
      alert('All fields are required.');
      return;
    }

    this.newsService.createNews(this.newNews).subscribe({
      next: () => {
        this.newNews = { employer_id: '', title: '', description: '' };
        this.loadNews(); 
      },
      error: err => console.error('Error creating news', err)
    });
  }

  deleteNews(id: string): void {
    this.newsService.deleteNews(id).subscribe(() => {
      this.newsList = this.newsList.filter(news => news._id !== id);
    });
  }
}