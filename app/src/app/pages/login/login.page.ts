import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { IonicModule, ModalController } from '@ionic/angular';
import { ForgotPasswordPage } from '../forgot-password/forgot-password.page';
import { AuthService } from 'src/app/services/auth.service';
import { TranslocoModule } from '@ngneat/transloco';

@Component({
  selector: 'app-login',
  templateUrl: './login.page.html',
  styleUrls: ['./login.page.scss'],
  standalone: true,
  imports: [IonicModule, CommonModule, FormsModule, TranslocoModule],
})
export class LoginPage {
  constructor(
    protected modalCtrl: ModalController,
    private auth: AuthService
  ) {}

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

  async dismissModal() {
    await this.modalCtrl.dismiss();
  }

  async signIn(userMail: any, userPassword: any) {
    const succeeded = await this.auth.SignIn(userMail, userPassword);
    if (succeeded) {
      await this.dismissModal();
    }
  }
}
