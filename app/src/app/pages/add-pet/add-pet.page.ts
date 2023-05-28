import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { AlertController, IonicModule, ModalController } from '@ionic/angular';

import { v4 as uuidv4 } from 'uuid';

import { Pet, AnimalSpecies } from '../../types/types';
import { PhotoService } from '../../services/photo.service';
import { FileStorageService } from 'src/app/services/file-storage.service';
import { AuthService } from 'src/app/services/auth.service';
import { TranslocoModule, TranslocoService } from '@ngneat/transloco';

@Component({
  selector: 'app-add-pet',
  templateUrl: 'add-pet.page.html',
  styleUrls: ['add-pet.page.scss'],
  standalone: true,
  imports: [IonicModule, CommonModule, TranslocoModule],
})
export class AddPetPage {
  animalSpecies = AnimalSpecies;
  pet: Pet;

  constructor(
    private modalCtrl: ModalController,
    private photoService: PhotoService,
    private auth: AuthService,
    private fileStorage: FileStorageService,
    private alertCtrl: AlertController
  ) {
    this.pet = {
      uuid: uuidv4(),
      userUid: '',
      name: '',
      species: undefined,
      image: '',
    };
  }

  cancel() {
    return this.modalCtrl.dismiss(null, 'cancel');
  }

  save() {
    if (!this.pet.name) {
      this.alertCtrl
        .create({
          header: 'Error',
          subHeader: 'Not all required fields are filled',
          message: 'Please fill in a name',
          buttons: ['OK'],
        })
        .then((alert) => alert.present());
      return;
    }

    if (!this.pet.species) {
      this.alertCtrl
        .create({
          header: 'Error',
          subHeader: 'Not all required fields are filled',
          message: 'Please fill in a species',
          buttons: ['OK'],
        })
        .then((alert) => alert.present());
      return;
    }

    if (!this.pet.image) {
      this.alertCtrl
        .create({
          header: 'Error',
          subHeader: 'Not all required fields are filled',
          message: 'Please upload an image',
          buttons: ['OK'],
        })
        .then((alert) => alert.present());
      return;
    }

    return this.modalCtrl.dismiss(this.pet, 'save');
  }

  handleNameChange(event: any) {
    this.pet.name = event.detail.value;
  }

  handleSpeciesChange(event: any) {
    this.pet.species = event.detail.value;
  }

  async addPicture() {
    const imageDataUrl = await this.photoService.getPhoto();
    const userUid = this.auth.currentUserUid;

    if (imageDataUrl && userUid) {
      this.pet.image = await this.fileStorage.uploadFile(
        userUid,
        this.pet.uuid,
        imageDataUrl
      );
    }
  }
}
