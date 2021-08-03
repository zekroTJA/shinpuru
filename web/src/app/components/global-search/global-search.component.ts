import { Component, EventEmitter, OnInit, Output } from '@angular/core';
import { Router } from '@angular/router';
import { debounceTime } from 'rxjs/operators';
import { Member, SearchResult } from 'src/app/api/api.models';
import { APIService } from 'src/app/api/api.service';

@Component({
  selector: 'app-global-search',
  templateUrl: './global-search.component.html',
  styleUrls: ['./global-search.component.scss'],
})
export class GlobalSearchComponent implements OnInit {
  searchInput = new EventEmitter<string>();
  lastResult: SearchResult;

  @Output() navigate = new EventEmitter<string[]>();

  constructor(private api: APIService) {}

  ngOnInit(): void {
    this.searchInput.pipe(debounceTime(500)).subscribe(async (query) => {
      if (!query) {
        this.lastResult = null;
        return;
      }
      this.lastResult = await this.api.getSearch(query, 10).toPromise();
    });
    document.getElementById('searchbar').focus();
  }

  nav(...routes: string[]) {
    this.searchInput = null;
    this.navigate.emit(routes);
  }
}
