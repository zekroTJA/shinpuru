import { Component, OnInit } from '@angular/core';
import { APIService } from 'src/app/api/api.service';
import { getCryptoRandomString } from 'src/app/utils/crypto';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss'],
})
export class LoginComponent implements OnInit {
  code: string;

  constructor(private api: APIService) {}

  async ngOnInit() {
    while (await this.pushCode());
  }

  private generatePushCode() {
    this.code = getCryptoRandomString(16);
  }

  private async pushCode(): Promise<boolean> {
    this.generatePushCode();
    try {
      await this.api.postPushCode(this.code).toPromise();
      window.location.assign('/guilds');
      return false;
    } catch {}
    return true;
  }
}
