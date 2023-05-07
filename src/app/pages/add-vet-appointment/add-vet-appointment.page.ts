import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { IonicModule, ModalController } from '@ionic/angular';

import { v4 as uuidv4 } from 'uuid';

import { PetVetAppointment } from '../../types/types';

@Component({
  selector: 'app-add-vet-appointment',
  templateUrl: 'add-vet-appointment.page.html',
  styleUrls: ['add-vet-appointment.page.scss'],
  standalone: true,
  imports: [IonicModule, CommonModule],
})
export class AddVetAppointmentPage {
  vetAppointment: PetVetAppointment;

  constructor(private modalCtrl: ModalController) {
    this.vetAppointment = {
      uuid: uuidv4(),
      name: '',
    };
  }

  cancel() {
    return this.modalCtrl.dismiss(null, 'cancel');
  }

  save() {
    return this.modalCtrl.dismiss(this.vetAppointment, 'save');
  }

  handleNameChange(event: any) {
    this.vetAppointment.name = event.detail.value;
  }
}
