import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { IonicModule, ModalController, ToastController } from '@ionic/angular';

import { PetPage } from '../pet/pet.page';
import { Pet } from '../../types/types';
import { AddPetPage } from '../add-pet/add-pet.page';

import { ApiService } from 'src/app/services/api.service';
import { Observable, catchError, finalize, of, tap } from 'rxjs';
import { AuthService } from 'src/app/services/auth.service';

@Component({
  selector: 'app-pets',
  templateUrl: 'pets.page.html',
  styleUrls: ['pets.page.scss'],
  standalone: true,
  imports: [IonicModule, CommonModule],
})
export class PetsPage implements OnInit {
  myPets$: Observable<Pet[]> | undefined;

  constructor(
    private modalCtrl: ModalController,
    private toastCtrl: ToastController,
    private api: ApiService,
    protected auth: AuthService
  ) {}

  ngOnInit() {
    this.loadPets();
  }

  loadPets() {
    this.myPets$ = this.api.getPets().pipe(
      tap(() => console.log('Action performed before any other')),
      catchError((err) => {
        this.toastCtrl
          .create({
            message: err.message,
            position: 'middle',
            color: 'danger',
            duration: 10000,
          })
          .then((toast) => toast.present());
        console.error('Error emitted');
        return of([]);
      }),
      finalize(() => console.log('Action to be executed always'))
    );
  }

  async openPetModal(pet: Pet) {
    console.log('opening page for ' + pet.name);

    const modal = await this.modalCtrl.create({
      component: PetPage,
      componentProps: { pet },
    });
    modal.present();

    const { data, role } = await modal.onWillDismiss();
    if (role !== 'cancel') {
      this.myPets$ = this.api.updatePet(data).pipe(
        tap(() => console.log('Action performed before any other')),
        catchError((err) => {
          this.toastCtrl
            .create({
              message: err.message,
              position: 'middle',
              color: 'danger',
              duration: 10000,
            })
            .then((toast) => toast.present());
          console.error('Error emitted');
          return of([]);
        }),
        finalize(() => console.log('Action to be executed always'))
      );
    }
  }

  async addPet() {
    console.log('opening new pet modal');

    const modal = await this.modalCtrl.create({
      component: AddPetPage,
    });
    modal.present();

    const { data, role } = await modal.onWillDismiss();

    if (role !== 'cancel') {
      this.myPets$ = this.api.addPet(data).pipe(
        tap(() => console.log('Action performed before any other')),
        catchError((err) => {
          this.toastCtrl
            .create({
              message: err.message,
              position: 'middle',
              color: 'danger',
              duration: 10000,
            })
            .then((toast) => toast.present());
          console.error('Error emitted');
          return of([]);
        }),
        finalize(() => console.log('Action to be executed always'))
      );
    }
  }
}
