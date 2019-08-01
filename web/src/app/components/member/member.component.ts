/** @format */

import { Component, Input } from '@angular/core';
import { Member } from 'src/app/api/api.models';
import { toHexClr } from 'src/app/ts/utils';

@Component({
  selector: 'app-member',
  templateUrl: './member.component.html',
  styleUrls: ['./member.component.sass'],
})
export class MemberComponent {
  @Input() public member: Member;

  public toHexClr = toHexClr;
}
