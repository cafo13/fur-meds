import { Component, HostListener, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { IonicModule, ModalController } from '@ionic/angular';
import { AuthService } from 'src/app/services/auth.service';
import { TranslocoModule } from '@ngneat/transloco';
import { ShowPasswordComponent } from 'src/app/components/show-password/show-password.component';

@Component({
  selector: 'app-registration',
  templateUrl: './registration.page.html',
  styleUrls: ['./registration.page.scss'],
  standalone: true,
  imports: [
    IonicModule,
    CommonModule,
    FormsModule,
    TranslocoModule,
    ShowPasswordComponent,
  ],
})
export class RegistrationPage {
  constructor(private modalCtrl: ModalController, private auth: AuthService) {}

  async dismissModal() {
    await this.modalCtrl.dismiss();
  }

  async signUp(userMail: any, userPassword: any) {
    await this.auth.SignUp(userMail, userPassword);
    await this.dismissModal();
  }
}
