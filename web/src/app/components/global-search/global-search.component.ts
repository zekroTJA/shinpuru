import { Component, EventEmitter, OnInit, Output } from '@angular/core';
import { Router } from '@angular/router';
import { debounceTime } from 'rxjs/operators';
import { Guild, Member, SearchResult } from 'src/app/api/api.models';
import { APIService } from 'src/app/api/api.service';

type Element = Guild | Member;

@Component({
  selector: 'app-global-search',
  templateUrl: './global-search.component.html',
  styleUrls: ['./global-search.component.scss'],
})
export class GlobalSearchComponent implements OnInit {
  searchInput = new EventEmitter<string>();
  lastResult: SearchResult;
  elements: Element[] = [];
  selected = -1;

  @Output() navigate = new EventEmitter<string[]>();

  constructor(private api: APIService) {}

  ngOnInit(): void {
    this.searchInput.pipe(debounceTime(500)).subscribe(async (query) => {
      if (!query) {
        this.reset();
        return;
      }
      this.lastResult = await this.api.getSearch(query, 10).toPromise();
      this.elements = this.elements
        .concat(this.lastResult.guilds)
        .concat(this.lastResult.members);
      console.log(this.lastResult);
    });
    document.getElementById('searchbar').focus();
  }

  nav(e: Element) {
    this.reset();
    const m = e as Member;
    if (m.user) this.navigate.emit(['guilds', m.guild_id, m.user.id]);
    else this.navigate.emit(['guilds', (e as Guild).id]);
  }

  onKeyUp(idx: number) {
    if (idx === 0) this.nav(this.elements[this.selected]);
    if (idx === -1 && this.selected === 0) return;
    if (idx === 1 && this.selected + 1 === this.elements.length) return;
    this.selected += idx;
  }

  isSelected(v: Element) {
    return this.elements[this.selected] === v;
  }

  private reset() {
    this.lastResult = null;
    this.elements = [];
  }
}
