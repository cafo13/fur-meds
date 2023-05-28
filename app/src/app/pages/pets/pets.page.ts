import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { AlertController, IonicModule, ModalController } from '@ionic/angular';

import { PetPage } from '../pet/pet.page';
import { Pet, AnimalSpecies } from '../../types/types';
import { AddPetPage } from '../add-pet/add-pet.page';

import { ApiService } from 'src/app/services/api.service';
import { Observable, catchError, finalize, of, tap } from 'rxjs';
import { AuthService } from 'src/app/services/auth.service';
import { TranslocoModule, TranslocoService } from '@ngneat/transloco';
import { Router } from '@angular/router';

@Component({
  selector: 'app-pets',
  templateUrl: 'pets.page.html',
  styleUrls: ['pets.page.scss'],
  standalone: true,
  imports: [IonicModule, CommonModule, TranslocoModule],
})
export class PetsPage implements OnInit {
  myPets$: Observable<Pet[]> | undefined;

  constructor(
    private modalCtrl: ModalController,
    private alertCtrl: AlertController,
    private api: ApiService,
    private transloco: TranslocoService,
    protected router: Router,
    protected auth: AuthService
  ) {}

  ngOnInit() {
    if (this.auth.isLoggedIn) {
      this.loadPets();
    }
  }

  loadPets() {
    this.myPets$ = this.api.getPets().pipe(
      tap(() => console.log('Action performed before any other')),
      catchError((err) => {
        this.alertCtrl
          .create({
            header: this.transloco.translate('global.error'),
            subHeader: this.transloco.translate(
              'pages.mypets.load_pets.request_failed_alert_subheader'
            ),
            message: err.message,
            buttons: [this.transloco.translate('global.ok')],
          })
          .then((alert) => alert.present());
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
    if (role === 'delete') {
      this.loadPets();
    } else if (role !== 'cancel') {
      this.myPets$ = this.api.updatePet(data).pipe(
        tap(() => console.log('Action performed before any other')),
        catchError((err) => {
          this.alertCtrl
            .create({
              header: this.transloco.translate('global.error'),
              subHeader: this.transloco.translate(
                'pages.mypets.update_pet.request_failed_alert_subheader'
              ),
              message: err.message,
              buttons: [this.transloco.translate('global.ok')],
            })
            .then((alert) => alert.present());
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
          this.alertCtrl
            .create({
              header: this.transloco.translate('global.error'),
              subHeader: this.transloco.translate(
                'pages.mypets.add_pet.request_failed_alert_subheader'
              ),
              message: err.message,
              buttons: [this.transloco.translate('global.ok')],
            })
            .then((alert) => alert.present());
          console.error('Error emitted');
          return of([]);
        }),
        finalize(() => console.log('Action to be executed always'))
      );
    }
  }

  getSpeciesText(species: AnimalSpecies): string {
    const indexOfSpecies = Object.keys(AnimalSpecies).indexOf(species);
    return Object.values(AnimalSpecies)[indexOfSpecies];
  }
}
