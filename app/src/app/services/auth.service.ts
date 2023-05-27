import { Injectable } from '@angular/core';
import { AngularFireAuth } from '@angular/fire/compat/auth';
import { Router } from '@angular/router';
import { AlertController } from '@ionic/angular';
import { TranslocoService } from '@ngneat/transloco';

import * as auth from 'firebase/auth';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  userData: any;

  constructor(
    private fireAuth: AngularFireAuth,
    private router: Router,
    private alertCtrl: AlertController,
    private transloco: TranslocoService
  ) {
    this.fireAuth.authState.subscribe((user) => {
      if (user) {
        this.userData = user;
        localStorage.setItem('user', JSON.stringify(this.userData));
        JSON.parse(localStorage.getItem('user')!);
      } else {
        localStorage.setItem('user', 'null');
        JSON.parse(localStorage.getItem('user')!);
      }
    });
  }

  get isLoggedIn(): boolean {
    const user = JSON.parse(localStorage.getItem('user')!);
    return user !== null ? true : false;
  }

  get currentUserUid(): string | undefined {
    const user = JSON.parse(localStorage.getItem('user')!);
    if (user !== null && user.uid) {
      return user.uid;
    }

    return undefined;
  }

  async SignIn(email: string, password: string): Promise<boolean> {
    return await this.fireAuth
      .signInWithEmailAndPassword(email, password)
      .then(async (user) => {
        if (user && user.user) {
          this.router.navigate(['/tabs/pets']);
          return true;
        }
        return false;
      })
      .catch(async (error) => {
        const alert = await this.alertCtrl.create({
          header: this.transloco.translate('global.error'),
          subHeader: this.transloco.translate('services.auth.signin_error'),
          message: error.message,
          buttons: [this.transloco.translate('global.ok')],
        });
        await alert.present();
        console.error(error.message);
        return false;
      });
  }

  async SignOut() {
    return this.fireAuth.signOut().then(() => {
      localStorage.removeItem('user');
      this.router.navigate(['/tabs/account']);
    });
  }

  async SignUp(email: string, password: string) {
    return this.fireAuth
      .createUserWithEmailAndPassword(email, password)
      .then((result) => {
        /* Call the SendVerificaitonMail() function when new user sign
        up and returns promise */
        this.SendVerificationMail();
        this.userData = result.user;
        localStorage.setItem('user', JSON.stringify(result.user));
      })
      .catch(async (error) => {
        const alert = await this.alertCtrl.create({
          header: this.transloco.translate('global.error'),
          subHeader: this.transloco.translate('services.auth.signup_error'),
          message: error.message,
          buttons: [this.transloco.translate('global.ok')],
        });
        await alert.present();
        console.error(error.message);
      });
  }

  async SendVerificationMail() {
    return this.fireAuth.currentUser
      .then((u: any) => u.sendEmailVerification())
      .then(async () => {
        const alert = await this.alertCtrl.create({
          header: this.transloco.translate(
            'services.auth.registration_successful_header'
          ),
          message: this.transloco.translate(
            'services.auth.registration_successful_message'
          ),
          buttons: [this.transloco.translate('global.ok')],
        });
        await alert.present();
      });
  }

  async ForgotPassword(passwordResetEmail: string) {
    return this.fireAuth
      .sendPasswordResetEmail(passwordResetEmail)
      .then(async () => {
        const alert = await this.alertCtrl.create({
          header: this.transloco.translate(
            'services.auth.password_reset_header'
          ),
          subHeader: this.transloco.translate(
            'services.auth.password_reset_subheader'
          ),
          message: this.transloco.translate(
            'services.auth.password_reset_message'
          ),
          buttons: [this.transloco.translate('global.ok')],
        });
        await alert.present();
      })
      .catch(async (error) => {
        const alert = await this.alertCtrl.create({
          header: this.transloco.translate('global.error'),
          subHeader: this.transloco.translate(
            'services.auth.password_reset_error'
          ),
          message: error.message,
          buttons: [this.transloco.translate('global.ok')],
        });
        await alert.present();
        console.error(error.message);
      });
  }

  async SignInWithGoogle() {
    return this.fireAuth
      .signInWithPopup(new auth.GoogleAuthProvider())
      .then((result) => {
        this.router.navigate(['/tabs/pets']);
        this.userData = result.user;
        localStorage.setItem('user', JSON.stringify(result.user));
      })
      .catch(async (error) => {
        const alert = await this.alertCtrl.create({
          header: this.transloco.translate('global.error'),
          subHeader: this.transloco.translate(
            'services.auth.signin_with_google_error'
          ),
          message: error.message,
          buttons: [this.transloco.translate('global.ok')],
        });
        await alert.present();
        console.error(error.message);
      });
  }
}
