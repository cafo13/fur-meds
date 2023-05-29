import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { AlertController, IonicModule, ModalController } from '@ionic/angular';

import { AnimalSpecies, Pet, PetFood, PetMedicine } from '../../types/types';
import { SetMedicinePage } from '../set-medicine/set-medicine.page';
import { SetFoodPage } from '../set-food/set-food.page';
import { TranslocoModule, TranslocoService } from '@ngneat/transloco';
import { catchError, finalize, of, tap } from 'rxjs';
import { ApiService } from 'src/app/services/api.service';
import { PhotoService } from 'src/app/services/photo.service';
import { FileStorageService } from 'src/app/services/file-storage.service';

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

  animalSpecies = AnimalSpecies;

  constructor(
    private modalCtrl: ModalController,
    private alertCtrl: AlertController,
    private transloco: TranslocoService,
    private api: ApiService,
    private photoService: PhotoService,
    private fileStorage: FileStorageService
  ) {
    this.pet = this.input;
  }

  handleNameChange(event: any) {
    if (this.pet) {
      this.pet.name = event.detail.value;
    }
  }

  handleSpeciesChange(event: any) {
    if (this.pet) {
      this.pet.species = event.detail.value;
    }
  }

  async addPicture() {
    const imageDataUrl = await this.photoService.getPhoto();

    if (imageDataUrl && this.pet?.userUid) {
      this.pet.image = await this.fileStorage.uploadFile(
        this.pet.userUid,
        this.pet.uuid,
        imageDataUrl
      );
    }
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

  async inviteUserToSharedPet(event: any) {
    if (!this.pet) {
      return;
    }

    const alert = await this.alertCtrl.create({
      header: this.transloco.translate('pages.pet.share_pet.header'),
      message: this.transloco.translate('pages.pet.share_pet.confirm_text'),
      inputs: [
        {
          name: 'email',
          type: 'email',
        },
      ],
      buttons: [
        {
          text: this.transloco.translate('global.cancel_button'),
          role: 'cancel',
          handler: () => {
            console.log('Cancelled inviting to share pet ' + this.pet?.uuid);
          },
        },
        {
          text: this.transloco.translate('global.share_button'),
          role: 'confirm',
          handler: (values: any) => {
            if (this.pet && this.pet.uuid) {
              this.api
                .inviteUserToSharedPet(this.pet.uuid, values['email'])
                .pipe(
                  tap(() => console.log('Action performed before any other')),
                  catchError((err) => {
                    this.alertCtrl
                      .create({
                        header: this.transloco.translate('global.error'),
                        subHeader: this.transloco.translate(
                          'pages.pet.share_pet_invite.request_failed_alert_subheader'
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
                .subscribe();
            }
          },
        },
      ],
    });
    await alert.present();
  }

  async deletePet() {
    if (!this.pet) {
      return;
    }

    const alert = await this.alertCtrl.create({
      header: this.transloco.translate('pages.pet.delete_pet.header'),
      message: this.transloco.translate('pages.pet.delete_pet.confirm_text', {
        name: this.pet.name,
      }),
      buttons: [
        {
          text: this.transloco.translate('global.cancel_button'),
          role: 'cancel',
          handler: () => {
            console.log('Cancelled deleting pet ' + this.pet?.uuid);
          },
        },
        {
          text: this.transloco.translate('global.delete_button'),
          role: 'confirm',
          handler: () => {
            if (this.pet && this.pet.uuid) {
              this.api
                .deletePet(this.pet.uuid)
                .pipe(
                  tap(() => console.log('Action performed before any other')),
                  catchError((err) => {
                    this.alertCtrl
                      .create({
                        header: this.transloco.translate('global.error'),
                        subHeader: this.transloco.translate(
                          'pages.pet.delete_pet.request_failed_alert_subheader'
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
                .subscribe(() => {
                  console.log('Successfully deleted pet ' + this.pet?.uuid);
                  this.modalCtrl.dismiss(null, 'delete');
                });
            }
          },
        },
      ],
    });
    await alert.present();
  }

  cancel() {
    return this.modalCtrl.dismiss(null, 'cancel');
  }

  save() {
    return this.modalCtrl.dismiss(this.pet, 'save');
  }
}
