import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { IonicModule, ModalController } from '@ionic/angular';
import { ForgotPasswordPage } from '../forgot-password/forgot-password.page';

@Component({
  selector: 'app-login',
  templateUrl: './login.page.html',
  styleUrls: ['./login.page.scss'],
  standalone: true,
  imports: [IonicModule, CommonModule, FormsModule],
})
export class LoginPage {
  constructor(protected modalCtrl: ModalController) {}

  async openForgotPasswordModal() {
    const modal = await this.modalCtrl.create({
      component: ForgotPasswordPage,
    });
    modal.present();

    const { data, role } = await modal.onWillDismiss();
    if (role !== 'cancel') {
      // something to do here?
    }
  }
}
