import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { AlertController, IonicModule, ModalController } from '@ionic/angular';

import { Pet, PetFood, PetMedicine } from '../../types/types';
import { SetMedicinePage } from '../set-medicine/set-medicine.page';
import { AddFoodPage } from '../add-food/add-food.page';

@Component({
  selector: 'app-pet',
  templateUrl: 'pet.page.html',
  styleUrls: ['pet.page.scss'],
  standalone: true,
  imports: [IonicModule, CommonModule],
})
export class PetPage {
  @Input() pet: Pet | undefined = undefined;

  constructor(
    private modalCtrl: ModalController,
    private alertCtrl: AlertController
  ) {}

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
      message:
        "Are you sure you want to delete the medicine '" + medicine.name + "'?",
      buttons: [
        {
          text: 'Cancel',
          role: 'cancel',
          handler: () => {
            console.log('Cancelled deleting medicine ' + medicine.uuid);
          },
        },
        {
          text: 'Delete',
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
    alert.present();
  }

  async addFood() {
    console.log('opening add add medicine modal');

    const modal = await this.modalCtrl.create({
      component: AddFoodPage,
    });
    modal.present();

    const { data, role } = await modal.onWillDismiss();

    if (role !== 'cancel') {
      if (this.pet && !this.pet?.foods) {
        this.pet.foods = [];
      }
      this.pet?.foods?.push(data);
    }
  }

  async updateFood(food: PetFood) {
    console.log('updating food ' + food.uuid);

    const alert = await this.alertCtrl.create({
      inputs: [
        {
          placeholder: 'Name',
          value: food.name,
        },
      ],
      buttons: [
        {
          text: 'Cancel',
          role: 'cancel',
          handler: () => {
            console.log('Cancelled updating food ' + food.uuid);
          },
        },
        {
          text: 'Save',
          role: 'confirm',
          handler: (value: any) => {
            console.log('Successfully updated food ' + food.uuid);
            food.name = value[0];
            this.updatePetFoodItem(food);
          },
        },
      ],
    });
    alert.present();
  }

  updatePetFoodItem(updatedFood: PetFood): void {
    if (!this.pet || !this.pet.foods) {
      return;
    }
    for (let food of this.pet.foods) {
      if (food.uuid === updatedFood.uuid) {
        food.name = updatedFood.name;
      }
    }
  }

  async deleteFood(food: PetFood) {
    console.log('deleting food ' + food.uuid);

    const alert = await this.alertCtrl.create({
      message: "Are you sure you want to delete the food '" + food.name + "'?",
      buttons: [
        {
          text: 'Cancel',
          role: 'cancel',
          handler: () => {
            console.log('Cancelled deleting food ' + food.uuid);
          },
        },
        {
          text: 'Delete',
          role: 'confirm',
          handler: () => {
            console.log('Successfully deleted food ' + food.uuid);
            this.deletePetFoodItem(food.uuid);
          },
        },
      ],
    });
    alert.present();
  }

  deletePetFoodItem(uuid: string) {
    if (!this.pet || !this.pet.foods) {
      return;
    }
    this.pet.foods.forEach((food, index) => {
      if (food.uuid === uuid) {
        this.pet?.foods?.splice(index, 1);
      }
    });
  }

  cancel() {
    return this.modalCtrl.dismiss(null, 'cancel');
  }

  save() {
    return this.modalCtrl.dismiss(this.pet, 'save');
  }
}