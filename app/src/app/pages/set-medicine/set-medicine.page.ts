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
          name: 'time',
          value: dataForUpdate?.time,
          type: 'time',
        },
        {
          name: 'everyDays',
          value: dataForUpdate?.everyDays,
          type: 'number',
          min: 0,
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
          handler: (values: any) => {
            const frequency: PetMedicineFrequency = {
              uuid: dataForUpdate ? dataForUpdate.uuid : uuidv4(),
              time: values['time'],
              everyDays: +values['everyDays'],
            };
            if (mode === 'Update') {
              for (let existingFrequency of this.medicine.frequencies) {
                if (existingFrequency.uuid === dataForUpdate?.uuid) {
                  existingFrequency.time = frequency.time;
                  existingFrequency.everyDays = +frequency.everyDays;
                }
              }
            } else {
              this.medicine.frequencies.push(frequency);
            }

            console.log('Successfully set frequency (mode: ' + mode + ')');
          },
        },
      ],
    });
    await alert.present();
  }

  async deleteFrequency(frequency: PetMedicineFrequency) {
    console.log('deleting frequency ' + frequency.uuid);

    const alert = await this.alertCtrl.create({
      message:
        "Are you sure you want to delete the frequency '" +
        frequency.time +
        "' every " +
        frequency.everyDays +
        'days?',
      buttons: [
        {
          text: 'Cancel',
          role: 'cancel',
          handler: () => {
            console.log('Cancelled deleting frequency ' + frequency.uuid);
          },
        },
        {
          text: 'Delete',
          role: 'confirm',
          handler: () => {
            this.medicine.frequencies.forEach((existingFrequency, index) => {
              if (frequency.uuid === frequency.uuid) {
                this.medicine.frequencies.splice(index, 1);
              }
            });
            console.log('Successfully deleted frequency ' + frequency.uuid);
          },
        },
      ],
    });
    await alert.present();
  }
}
