import { CommonModule } from '@angular/common';
import { Component, ContentChild, OnInit } from '@angular/core';
import { IonInput, IonicModule } from '@ionic/angular';

@Component({
  selector: 'app-show-password',
  templateUrl: './show-password.component.html',
  styleUrls: ['./show-password.component.scss'],
  standalone: true,
  imports: [IonicModule, CommonModule],
})
export class ShowPasswordComponent {
  @ContentChild(IonInput) input: IonInput | undefined;
  showPassword = false;

  constructor() {}

  toggleShow() {
    this.showPassword = !this.showPassword;
    if (this.input) {
      this.input.type = this.showPassword ? 'text' : 'password';
    }
  }
}
