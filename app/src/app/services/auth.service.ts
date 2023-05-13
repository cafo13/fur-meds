import { Injectable } from '@angular/core';
import { AngularFireAuth } from '@angular/fire/compat/auth';
import { Router } from '@angular/router';
import { ToastController } from '@ionic/angular';

import * as auth from 'firebase/auth';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  userData: any;

  constructor(
    private fireAuth: AngularFireAuth,
    private router: Router,
    private toastCtrl: ToastController
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
      .catch((error) => {
        this.toastCtrl.create({
          message: error.message,
          position: 'bottom',
          color: 'danger',
        });
      });
  }

  async SignOut() {
    return this.fireAuth.signOut().then(() => {
      localStorage.removeItem('user');
      this.router.navigate(['/tabs/account']);
    });
  }

  get isLoggedIn(): boolean {
    const user = JSON.parse(localStorage.getItem('user')!);
    return user !== null ? true : false;
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
      .catch((error) => {
        this.toastCtrl.create({
          message: error.message,
          position: 'bottom',
          color: 'danger',
        });
      });
  }

  async SendVerificationMail() {
    return this.fireAuth.currentUser
      .then((u: any) => u.sendEmailVerification())
      .then(() => {
        this.router.navigate(['verify-email-address']);
      });
  }

  async ForgotPassword(passwordResetEmail: string) {
    return this.fireAuth
      .sendPasswordResetEmail(passwordResetEmail)
      .then(() => {
        this.toastCtrl.create({
          message: 'Password reset email sent, check your inbox.',
          position: 'bottom',
          color: 'success',
        });
      })
      .catch((error) => {
        this.toastCtrl.create({
          message: error.message,
          position: 'bottom',
          color: 'danger',
        });
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
      .catch((error) => {
        this.toastCtrl.create({
          message: error.message,
          position: 'bottom',
          color: 'danger',
        });
      });
  }
}
