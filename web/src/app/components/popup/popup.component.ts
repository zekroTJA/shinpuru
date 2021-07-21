/** @format */

import {
  Component,
  Input,
  OnInit,
  EventEmitter,
  Output,
  TemplateRef,
} from '@angular/core';

export interface PopupElement {
  el: string | TemplateRef<any>;
  action: () => void;
}

@Component({
  selector: 'app-popup',
  templateUrl: './popup.component.html',
  styleUrls: ['./popup.component.scss'],
})
export class PopupComponent implements OnInit {
  @Input() public elements: PopupElement[];
  @Output() public closing: EventEmitter<any> = new EventEmitter();

  constructor() {
    window.onclick = this.windowClickHandler.bind(this);
  }

  ngOnInit() {}

  private windowClickHandler(e: any) {
    this.closing.emit(e);
  }

  public isTemplate(element) {
    return element.el instanceof TemplateRef;
  }
}
