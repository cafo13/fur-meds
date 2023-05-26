import { Component } from '@angular/core';
import { IonicModule, ModalController } from '@ionic/angular';
import { AuthService } from 'src/app/services/auth.service';
import { LoginPage } from '../login/login.page';
import { RegistrationPage } from '../registration/registration.page';
import { CommonModule } from '@angular/common';
import { TranslocoModule } from '@ngneat/transloco';

@Component({
  selector: 'app-account',
  templateUrl: 'account.page.html',
  styleUrls: ['account.page.scss'],
  standalone: true,
  imports: [IonicModule, CommonModule, TranslocoModule],
})
export class AccountPage {
  constructor(
    protected modalCtrl: ModalController,
    protected auth: AuthService
  ) {}

  async openLoginModal() {
    const modal = await this.modalCtrl.create({
      component: LoginPage,
    });
    modal.present();

    const { data, role } = await modal.onWillDismiss();
    if (role !== 'cancel') {
      // something to do here?
    }
  }

  async openRegistrationModal() {
    const modal = await this.modalCtrl.create({
      component: RegistrationPage,
    });
    modal.present();

    const { data, role } = await modal.onWillDismiss();
    if (role !== 'cancel') {
      // something to do here?
    }
  }

  async signInWithGoogle() {
    await this.auth.SignInWithGoogle();
  }
}
