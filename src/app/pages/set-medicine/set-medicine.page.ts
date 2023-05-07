import { CommonModule } from '@angular/common';
import { Component, Input, OnInit } from '@angular/core';
import { AlertController, IonicModule, ModalController } from '@ionic/angular';

import { v4 as uuidv4 } from 'uuid';

import { PetMedicine, PetMedicineFrequency } from '../../types/types';

@Component({
  selector: 'app-set-medicine',
  templateUrl: 'set-medicine.page.html',
  styleUrls: ['set-medicine.page.scss'],
  standalone: true,
  imports: [IonicModule, CommonModule],
})
export class SetMedicinePage implements OnInit {
  medicine: PetMedicine = {
    uuid: uuidv4(),
    name: '',
    dosage: '',
    frequencies: [],
  };
  mode: 'Add' | 'Update' = 'Add';
  @Input() dataForUpdate: PetMedicine | undefined = undefined;

  constructor(
    private modalCtrl: ModalController,
    private alertCtrl: AlertController
  ) {}

  ngOnInit(): void {
    if (this.dataForUpdate) {
      this.mode = 'Update';
      this.medicine = this.dataForUpdate;
    }
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

  handleDosageChange(event: any) {
    this.medicine.dosage = event.detail.value;
  }

  async setFrequency(
    dataForUpdate: PetMedicineFrequency | undefined = undefined
  ) {
    const mode = dataForUpdate ? 'Update' : 'Add';
    console.log('opening set frequency popup with mode ' + mode);

    const alert = await this.alertCtrl.create({
      inputs: [
        {
          label: 'Time',
          type: 'number',
        },
        {
          label: 'Every X Days',
          type: 'number',
        },
      ],
      buttons: [
        {
          text: 'Cancel',
          role: 'cancel',
          handler: () => {
            console.log('Cancelled setting frequency (mode: ' + mode + ')');
          },
        },
        {
          text: 'Save',
          role: 'confirm',
          handler: (value: any) => {
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

            console.log('Successfully set frequency (mode: ' + mode + ')');
          },
        },
      ],
    });
    alert.present();
  }
}
