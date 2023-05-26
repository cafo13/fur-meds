import { CommonModule } from '@angular/common';
import { Component, Input, OnInit } from '@angular/core';
import { IonicModule, ModalController } from '@ionic/angular';
import { TranslocoModule } from '@ngneat/transloco';

import { v4 as uuidv4 } from 'uuid';

import { PetFood } from '../../types/types';

@Component({
  selector: 'app-set-food',
  templateUrl: 'set-food.page.html',
  styleUrls: ['set-food.page.scss'],
  standalone: true,
  imports: [IonicModule, CommonModule, TranslocoModule],
})
export class SetFoodPage implements OnInit {
  food: PetFood = {
    uuid: uuidv4(),
    name: '',
    dosage: '',
    frequencies: [],
  };
  mode: 'Add' | 'Update' = 'Add';
  @Input() dataForUpdate: PetFood | undefined = undefined;

  constructor(private modalCtrl: ModalController) {}

  ngOnInit(): void {
    if (this.dataForUpdate) {
      this.mode = 'Update';
      this.food = this.dataForUpdate;
    }
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
