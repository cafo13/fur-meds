import { Injectable } from '@angular/core';
import { AngularFireAuth } from '@angular/fire/compat/auth';
import { Router } from '@angular/router';
import { AlertController } from '@ionic/angular';

import * as auth from 'firebase/auth';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  userData: any;

  constructor(
    private fireAuth: AngularFireAuth,
    private router: Router,
    private alertCtrl: AlertController
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

  async SignIn(email: string, password: string) {
    return this.fireAuth
      .signInWithEmailAndPassword(email, password)
      .then((_user) => {
        this.fireAuth.authState.subscribe((user) => {
          if (user) {
            this.router.navigate(['']);
          }
        });
      })
      .catch(async (error) => {
        const alert = await this.alertCtrl.create({
          header: 'Error',
          subHeader: 'Authentication error on signin',
          message: error.message,
          buttons: ['OK'],
        });
        await alert.present();
        console.error(error.message);
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
          header: 'Error',
          subHeader: 'Authentication error on signup',
          message: error.message,
          buttons: ['OK'],
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
          header: 'Registration succesful',
          message: 'Email verification mail sent, check your inbox.',
          buttons: ['OK'],
        });
        await alert.present();
      });
  }

  async ForgotPassword(passwordResetEmail: string) {
    return this.fireAuth
      .sendPasswordResetEmail(passwordResetEmail)
      .then(async () => {
        const alert = await this.alertCtrl.create({
          header: 'Password reset',
          subHeader: 'Request for resetting password was succesful',
          message: 'Password reset email sent, check your inbox.',
          buttons: ['OK'],
        });
        await alert.present();
      })
      .catch(async (error) => {
        const alert = await this.alertCtrl.create({
          header: 'Error',
          subHeader: 'Authentication error at resetting password',
          message: error.message,
          buttons: ['OK'],
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
          header: 'Error',
          subHeader: 'Authentication error at Google sign in',
          message: error.message,
          buttons: ['OK'],
        });
        await alert.present();
        console.error(error.message);
      });
  }
}
