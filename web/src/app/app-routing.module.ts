/** @format */

import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { LoginComponent } from './routes/login/login.component';
import { HomeComponent } from './routes/home/home.component';
import { GuildComponent } from './routes/guild/guild.component';
import { MemberRouteComponent } from './routes/member/member.component';
const routes: Routes = [
  {
    path: '',
    redirectTo: '/guilds',
    pathMatch: 'full',
  },
  {
    path: 'guilds',
    component: HomeComponent,
  },
  {
    path: 'login',
    component: LoginComponent,
  },
  {
    path: 'guilds/:id',
    component: GuildComponent,
  },
  {
    path: 'guilds/:guildid/:memberid',
    component: MemberRouteComponent,
  },
  {
    path: '**',
    redirectTo: '/guilds',
    pathMatch: 'full',
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
