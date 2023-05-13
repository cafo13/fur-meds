import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { IonicModule, ModalController } from '@ionic/angular';

import { PetPage } from '../pet/pet.page';
import { Pet } from '../../types/types';
import { AddPetPage } from '../add-pet/add-pet.page';

import { ApiService } from 'src/app/services/api.service';
import { Observable } from 'rxjs';
import { AuthService } from 'src/app/services/auth.service';

@Component({
  selector: 'app-pets',
  templateUrl: 'pets.page.html',
  styleUrls: ['pets.page.scss'],
  standalone: true,
  imports: [IonicModule, CommonModule],
})
export class PetsPage {
  myPets$: Observable<Pet[]>;

  constructor(
    private modalCtrl: ModalController,
    private api: ApiService,
    protected auth: AuthService
  ) {
    this.myPets$ = api.getPets();
  }

  async openPetModal(pet: Pet) {
    console.log('opening page for ' + pet.name);

    const modal = await this.modalCtrl.create({
      component: PetPage,
      componentProps: { pet },
    });
    modal.present();

    const { data, role } = await modal.onWillDismiss();
    if (role !== 'cancel') {
      this.myPets$ = this.api.updatePet(data);
    }
  }

  async addPet() {
    console.log('opening new pet modal');

    const modal = await this.modalCtrl.create({
      component: AddPetPage,
    });
    modal.present();

    const { data, role } = await modal.onWillDismiss();

    if (role !== 'cancel') {
      this.myPets$ = this.api.addPet(data);
    }
  }
}
