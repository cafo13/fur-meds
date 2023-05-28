import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { AlertController, IonicModule, ModalController } from '@ionic/angular';

import { Pet, PetFood, PetMedicine } from '../../types/types';
import { SetMedicinePage } from '../set-medicine/set-medicine.page';
import { SetFoodPage } from '../set-food/set-food.page';
import { TranslocoModule, TranslocoService } from '@ngneat/transloco';
import { catchError, finalize, of, tap } from 'rxjs';
import { ApiService } from 'src/app/services/api.service';

@Component({
  selector: 'app-pet',
  templateUrl: 'pet.page.html',
  styleUrls: ['pet.page.scss'],
  standalone: true,
  imports: [IonicModule, CommonModule, TranslocoModule],
})
export class PetPage {
  @Input() input: Pet | undefined = undefined;
  pet: Pet | undefined;

  constructor(
    private modalCtrl: ModalController,
    private alertCtrl: AlertController,
    private transloco: TranslocoService,
    private api: ApiService
  ) {
    this.pet = this.input;
  }

  async setMedicine(dataForUpdate: PetMedicine | undefined = undefined) {
    const mode = dataForUpdate ? 'Update' : 'Add';
    console.log('opening set medicine modal with mode ' + mode);

    const modal = await this.modalCtrl.create({
      component: SetMedicinePage,
      componentProps: { dataForUpdate },
    });
    modal.present();

    const { data, role } = await modal.onWillDismiss();

    if (role === 'save') {
      if (this.pet && !this.pet?.medicines) {
        this.pet.medicines = [];
      }
      if (mode === 'Update' && this.pet && this.pet.medicines) {
        for (let medicine of this.pet.medicines) {
          if (medicine.uuid === data.uuid) {
            medicine = data;
          }
        }
      } else {
        this.pet?.medicines?.push(data);
      }
    }
  }

  async deleteMedicine(medicine: PetMedicine) {
    console.log('deleting medicine ' + medicine.uuid);

    const alert = await this.alertCtrl.create({
      message: this.transloco.translate(
        'pages.pet.delete_medicine_confirm_text',
        { medicine: medicine.name }
      ),
      buttons: [
        {
          text: this.transloco.translate('global.cancel_button'),
          role: 'cancel',
          handler: () => {
            console.log('Cancelled deleting medicine ' + medicine.uuid);
          },
        },
        {
          text: this.transloco.translate('global.delete_button'),
          role: 'confirm',
          handler: () => {
            console.log('Successfully deleted medicine ' + medicine.uuid);
            if (!this.pet || !this.pet.medicines) {
              return;
            }
            this.pet.medicines.forEach((existingMedicine, index) => {
              if (existingMedicine.uuid === medicine.uuid) {
                this.pet?.medicines?.splice(index, 1);
              }
            });
          },
        },
      ],
    });
    await alert.present();
  }

  async setFood(dataForUpdate: PetFood | undefined = undefined) {
    const mode = dataForUpdate ? 'Update' : 'Add';
    console.log('opening set food modal with mode ' + mode);

    const modal = await this.modalCtrl.create({
      component: SetFoodPage,
      componentProps: { dataForUpdate },
    });
    modal.present();

    const { data, role } = await modal.onWillDismiss();

    if (role === 'save') {
      if (this.pet && !this.pet?.foods) {
        this.pet.foods = [];
      }
      if (mode === 'Update' && this.pet && this.pet.foods) {
        for (let food of this.pet.foods) {
          if (food.uuid === data.uuid) {
            food = data;
          }
        }
      } else {
        this.pet?.foods?.push(data);
      }
    }
  }

  async deleteFood(food: PetFood) {
    console.log('deleting food ' + food.uuid);

    const alert = await this.alertCtrl.create({
      message: this.transloco.translate('pages.pet.delete_food_confirm_text', {
        food: food.name,
      }),
      buttons: [
        {
          text: this.transloco.translate('global.cancel_button'),
          role: 'cancel',
          handler: () => {
            console.log('Cancelled deleting food ' + food.uuid);
          },
        },
        {
          text: this.transloco.translate('global.delete_button'),
          role: 'confirm',
          handler: () => {
            console.log('Successfully deleted food ' + food.uuid);
            if (!this.pet || !this.pet.foods) {
              return;
            }
            this.pet.foods.forEach((existingFood, index) => {
              if (existingFood.uuid === food.uuid) {
                this.pet?.foods?.splice(index, 1);
              }
            });
          },
        },
      ],
    });
    await alert.present();
  }

  acceptInviteToSharedPet(event: any) {
    if (!this.pet) {
      return;
    }

    this.api
      .acceptInviteToSharedPet('petUuid')
      .pipe(
        tap(() => console.log('Action performed before any other')),
        catchError((err) => {
          this.alertCtrl
            .create({
              header: this.transloco.translate('global.error'),
              subHeader: this.transloco.translate(
                'pages.pets.accept_shared_pet.request_failed_alert_subheader'
              ),
              message: err.message,
              buttons: [this.transloco.translate('global.ok')],
            })
            .then((alert) => alert.present());
          console.error('Error emitted');
          return of([]);
        }),
        finalize(() => console.log('Action to be executed always'))
      )
      .subscribe()
      .unsubscribe();
  }

  inviteUserToSharedPet(event: any) {
    if (!this.pet) {
      return;
    }

    this.api
      .inviteUserToSharedPet(this.pet.uuid, 'levke97@gmail.com')
      .pipe(
        tap(() => console.log('Action performed before any other')),
        catchError((err) => {
          this.alertCtrl
            .create({
              header: this.transloco.translate('global.error'),
              subHeader: this.transloco.translate(
                'pages.pets.share_pet_invite.request_failed_alert_subheader'
              ),
              message: err.message,
              buttons: [this.transloco.translate('global.ok')],
            })
            .then((alert) => alert.present());
          console.error('Error emitted');
          return of([]);
        }),
        finalize(() => console.log('Action to be executed always'))
      )
      .subscribe()
      .unsubscribe();
  }

  cancel() {
    return this.modalCtrl.dismiss(null, 'cancel');
  }

  save() {
    return this.modalCtrl.dismiss(this.pet, 'save');
  }
}
