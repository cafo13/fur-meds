import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { IonicModule, ModalController } from '@ionic/angular';

import { v4 as uuidv4 } from 'uuid';

import { Pet, PetSpecies } from '../../types/types';
import { PhotoService } from '../../services/photo.service';

@Component({
  selector: 'app-add-pet',
  templateUrl: 'add-pet.page.html',
  styleUrls: ['add-pet.page.scss'],
  standalone: true,
  imports: [IonicModule, CommonModule],
})
export class AddPetPage {
  petSpecies = PetSpecies;
  pet: Pet;

  constructor(
    private modalCtrl: ModalController,
    private photoService: PhotoService
  ) {
    this.pet = {
      uuid: uuidv4(),
      name: '',
    };
  }

  cancel() {
    return this.modalCtrl.dismiss(null, 'cancel');
  }

  save() {
    return this.modalCtrl.dismiss(this.pet, 'save');
  }

  handleNameChange(event: any) {
    this.pet.name = event.detail.value;
  }

  handleSpeciesChange(event: any) {
    this.pet.species = event.detail.value;
  }

  async addPicture() {
    this.pet.image = await this.photoService.getPhoto();
  }
}
