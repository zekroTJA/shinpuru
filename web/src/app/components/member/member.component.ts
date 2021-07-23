/** @format */

import { Component, Input } from '@angular/core';
import { Member, Role } from 'src/app/api/api.models';
import { toHexClr } from 'src/app/utils/utils';

@Component({
  selector: 'app-member',
  templateUrl: './member.component.html',
  styleUrls: ['./member.component.scss'],
})
export class MemberComponent {
  @Input() public member: Member;
  @Input() public roles: Role[];

  public toHexClr = toHexClr;

  public getMemberColor(op: number): string {
    if (!this.roles) {
      return null;
    }

    const role = this.roles
      .filter((r) => r.color !== 0 && this.member.roles.includes(r.id))
      .sort((a, b) => b.position - a.position)[0];

    if (!role) {
      return null;
    }

    return this.toHexClr(role.color, op);
  }
}
