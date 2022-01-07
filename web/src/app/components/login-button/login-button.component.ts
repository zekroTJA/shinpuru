import { Component, EventEmitter, OnInit, Output } from '@angular/core';

@Component({
  selector: 'app-login-button',
  templateUrl: './login-button.component.html',
  styleUrls: ['./login-button.component.scss'],
})
export class LoginButtonComponent implements OnInit {
  @Output() login = new EventEmitter<any>();

  constructor() {}

  ngOnInit() {}

  onLogin() {
    this.login.emit();
    window.location.assign('/api/auth/login');
  }
}
