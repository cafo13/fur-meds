import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { IonicModule, ModalController } from '@ionic/angular';
import { AuthService } from 'src/app/services/auth.service';

@Component({
  selector: 'app-forgot-password',
  templateUrl: './forgot-password.page.html',
  styleUrls: ['./forgot-password.page.scss'],
  standalone: true,
  imports: [IonicModule, CommonModule, FormsModule],
})
export class ForgotPasswordPage {
  constructor(private modalCtrl: ModalController, private auth: AuthService) {}

  async dismissModal() {
    await this.modalCtrl.dismiss();
  }

  async forgotPassword(userMail: any) {
    await this.auth.ForgotPassword(userMail);
    await this.dismissModal();
  }
}
