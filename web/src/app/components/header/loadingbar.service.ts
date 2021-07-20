import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root',
})
export class LoadingBarService {
  private status = 0;

  public get isLoading(): boolean {
    return this.status > 0;
  }

  public init() {
    this.status++;
  }

  public finished() {
    this.status--;
  }

  public reset() {
    this.status = 0;
  }
}
