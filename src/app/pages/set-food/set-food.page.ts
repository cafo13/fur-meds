import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { IonicModule, ModalController } from '@ionic/angular';

import { v4 as uuidv4 } from 'uuid';

import { PetFood } from '../../types/types';

@Component({
  selector: 'app-set-food',
  templateUrl: 'set-food.page.html',
  styleUrls: ['set-food.page.scss'],
  standalone: true,
  imports: [IonicModule, CommonModule],
})
export class SetFoodPage {
  food: PetFood;

  constructor(private modalCtrl: ModalController) {
    this.food = {
      uuid: uuidv4(),
      name: '',
    };
  }

  cancel() {
    return this.modalCtrl.dismiss(null, 'cancel');
  }

  save() {
    return this.modalCtrl.dismiss(this.food, 'save');
  }

  handleNameChange(event: any) {
    this.food.name = event.detail.value;
  }
}
