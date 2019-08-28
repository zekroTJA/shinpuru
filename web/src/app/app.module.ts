/** @format */

import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { LoginComponent } from './routes/login/login.component';
import { HeaderComponent } from './components/header/header.component';
import { ToastComponent } from './components/toast/toast.component';
import { HomeComponent } from './routes/home/home.component';
import { SpinnerComponent } from './components/spinner/spinner.component';
import { GuildComponent } from './routes/guild/guild.component';
import { PopupComponent } from './components/popup/popup.component';
import { RoleBadgeComponent } from './components/rolebadge/rolebadge.component';
import { MemberComponent } from './components/member/member.component';
import { MemberRouteComponent } from './routes/member/member.component';
import { ReportComponent } from './components/report/report.component';
import { SpoilerComponent } from './components/spoiler/spoiler.component';
import { FormsModule } from '@angular/forms';
import { TagsInputComponent } from './components/tagsinput/tagsinput.component';
import { SettingsComponent } from './routes/settings/settings.component';
import { SpinnerButtonComponent } from './components/spinnerButton/spinnerButton.component';
import { SysInfoComponent } from './routes/sysinfo/sysinfo.component';
import { SliderComponent } from './components/slider/slider.component';

@NgModule({
  declarations: [
    AppComponent,
    LoginComponent,
    HeaderComponent,
    ToastComponent,
    HomeComponent,
    SpinnerComponent,
    GuildComponent,
    PopupComponent,
    RoleBadgeComponent,
    MemberComponent,
    MemberRouteComponent,
    ReportComponent,
    SpoilerComponent,
    TagsInputComponent,
    SettingsComponent,
    SpinnerButtonComponent,
    SysInfoComponent,
    SliderComponent,
  ],
  imports: [
    NgbModule,
    BrowserModule,
    AppRoutingModule,
    HttpClientModule,
    FormsModule,
  ],
  providers: [],
  bootstrap: [AppComponent],
})
export class AppModule {}
