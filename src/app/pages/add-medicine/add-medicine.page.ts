import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { IonicModule, ModalController } from '@ionic/angular';

import { PetMedicine } from '../../types/types';

@Component({
  selector: 'app-add-medicine',
  templateUrl: 'add-medicine.page.html',
  styleUrls: ['add-medicine.page.scss'],
  standalone: true,
  imports: [IonicModule, CommonModule],
})
export class AddMedicinePage {
  medicine: PetMedicine;

  constructor(private modalCtrl: ModalController) {
    this.medicine = {
      name: '',
    };
  }

  cancel() {
    return this.modalCtrl.dismiss(null, 'cancel');
  }

  save() {
    return this.modalCtrl.dismiss(this.medicine, 'save');
  }

  handleNameChange(event: any) {
    this.medicine.name = event.detail.value;
  }
}
