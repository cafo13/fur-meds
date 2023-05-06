import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { IonicModule, ModalController } from '@ionic/angular';

import { Pet } from '../../types/types';
import { AddMedicinePage } from '../add-medicine/add-medicine.page';
import { AddFoodPage } from '../add-food/add-food.page';
import { AddVetAppointmentPage } from '../add-vet-appointment/add-vet-appointment.page';

@Component({
  selector: 'app-pet',
  templateUrl: 'pet.page.html',
  styleUrls: ['pet.page.scss'],
  standalone: true,
  imports: [IonicModule, CommonModule],
})
export class PetPage {
  @Input() pet: Pet | undefined = undefined;

  constructor(private modalCtrl: ModalController) {}

  async addMedicine() {
    console.log('opening add add medicine modal');

    const modal = await this.modalCtrl.create({
      component: AddMedicinePage,
    });
    modal.present();

    const { data, role } = await modal.onWillDismiss();

    if (role !== 'cancel') {
      if (this.pet && !this.pet?.medicines) {
        this.pet.medicines = [];
      }
      this.pet?.medicines?.push(data);
    }
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

  async addVetAppointment() {
    console.log('opening add add medicine modal');

    const modal = await this.modalCtrl.create({
      component: AddVetAppointmentPage,
    });
    modal.present();

    const { data, role } = await modal.onWillDismiss();

    if (role !== 'cancel') {
      if (this.pet && !this.pet?.vetAppointments) {
        this.pet.vetAppointments = [];
      }
      this.pet?.vetAppointments?.push(data);
    }
  }

  cancel() {
    return this.modalCtrl.dismiss(null, 'cancel');
  }

  save() {
    return this.modalCtrl.dismiss(this.pet, 'save');
  }
}
