import { Component, OnInit } from '@angular/core';
import { NewsService } from 'src/app/services/news.service';
import { News } from 'src/app/model/news';
import { Router } from '@angular/router';

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

  constructor(private newsService: NewsService, private router: Router) {}

  ngOnInit(): void {
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