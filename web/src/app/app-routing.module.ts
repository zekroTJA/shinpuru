/** @format */

import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { LoginComponent } from './routes/login/login.component';
import { HomeComponent } from './routes/home/home.component';
import { GuildComponent } from './routes/guild/guild.component';
import { MemberRouteComponent } from './routes/member/member.component';
import { SettingsComponent } from './routes/settings/settings.component';
import { InfoComponent } from './routes/info/info.component';
import { APITokenComponent } from './routes/apitoken/apitoken.component';
import { ScoreboardComponent } from './routes/scoreboard/scoreboard.component';
import { CommandsComponent } from './routes/commands/commands.component';

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
    path: 'guilds/:guildid/scoreboard',
    component: ScoreboardComponent,
  },
  {
    path: 'guilds/:guildid/:memberid',
    component: MemberRouteComponent,
  },
  {
    path: 'settings',
    component: SettingsComponent,
  },
  {
    path: 'info',
    component: InfoComponent,
  },
  {
    path: 'apitoken',
    component: APITokenComponent,
  },
  {
    path: 'commands',
    component: CommandsComponent,
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
