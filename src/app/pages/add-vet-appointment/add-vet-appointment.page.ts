import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { IonicModule, ModalController } from '@ionic/angular';

import { VetAppointments } from '../../types/types';

@Component({
  selector: 'app-add-vet-appointment',
  templateUrl: 'add-vet-appointment.page.html',
  styleUrls: ['add-vet-appointment.page.scss'],
  standalone: true,
  imports: [IonicModule, CommonModule],
})
export class AddVetAppointmentPage {
  vetAppointment: VetAppointments;

  constructor(private modalCtrl: ModalController) {
    this.vetAppointment = {
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
