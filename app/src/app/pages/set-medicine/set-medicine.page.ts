import { CommonModule } from '@angular/common';
import { Component, Input, OnInit } from '@angular/core';
import { AlertController, IonicModule, ModalController } from '@ionic/angular';
import { TranslocoModule, TranslocoService } from '@ngneat/transloco';

import { v4 as uuidv4 } from 'uuid';

import { PetMedicine, PetMedicineFrequency } from '../../types/types';

@Component({
  selector: 'app-set-medicine',
  templateUrl: 'set-medicine.page.html',
  styleUrls: ['set-medicine.page.scss'],
  standalone: true,
  imports: [IonicModule, CommonModule, TranslocoModule],
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
    private alertCtrl: AlertController,
    protected transloco: TranslocoService
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
      header:
        mode === 'Update'
          ? this.transloco.translate(
              'pages.set_medicine.frequencies.update_title'
            )
          : this.transloco.translate(
              'pages.set_medicine.frequencies.add_title'
            ),
      message: this.transloco.translate(
        'pages.set_medicine.frequencies.explanation'
      ),
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
          text: this.transloco.translate('global.cancel_button'),
          role: 'cancel',
          handler: () => {
            console.log('Cancelled setting frequency (mode: ' + mode + ')');
          },
        },
        {
          text: this.transloco.translate('global.save_button'),
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
        frequency.everyDays === 1
          ? this.transloco.translate(
              'pages.set_medicine.frequencies.delete_frequency_confirm_text_one',
              { time: frequency.time }
            )
          : this.transloco.translate(
              'pages.set_medicine.frequencies.delete_frequency_confirm_text',
              { time: frequency.time, everyDays: frequency.everyDays }
            ),
      buttons: [
        {
          text: this.transloco.translate('global.cancel_button'),
          role: 'cancel',
          handler: () => {
            console.log('Cancelled deleting frequency ' + frequency.uuid);
          },
        },
        {
          text: this.transloco.translate('global.delete_button'),
          role: 'confirm',
          handler: () => {
            this.medicine.frequencies.forEach((existingFrequency, index) => {
              if (existingFrequency.uuid === frequency.uuid) {
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
