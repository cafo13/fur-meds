import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { AlertController, IonicModule, ModalController } from '@ionic/angular';

import { AnimalSpecies, PetShareInvite } from '../../types/types';
import { TranslocoModule, TranslocoService } from '@ngneat/transloco';
import { ApiService } from 'src/app/services/api.service';
import { catchError, finalize, of, tap } from 'rxjs';

@Component({
  selector: 'app-pet-invites',
  templateUrl: 'pet-invites.page.html',
  styleUrls: ['pet-invites.page.scss'],
  standalone: true,
  imports: [IonicModule, CommonModule, TranslocoModule],
})
export class PetInvitesPage {
  @Input() input: PetShareInvite[] | undefined = undefined;
  invites: PetShareInvite[] | undefined;

  constructor(
    private modalCtrl: ModalController,
    protected transloco: TranslocoService,
    private api: ApiService,
    private alertCtrl: AlertController
  ) {
    this.invites = this.input;
  }

  getSpeciesText(species: AnimalSpecies): string {
    const indexOfSpecies = Object.keys(AnimalSpecies).indexOf(species);
    return Object.values(AnimalSpecies)[indexOfSpecies];
  }

  acceptInvite(invite: PetShareInvite, inviteIndex: number) {
    this.api
      .acceptInviteToSharedPet(invite.pet.uuid)
      .pipe(
        tap(() => console.log('Action performed before any other')),
        catchError((err) => {
          this.alertCtrl
            .create({
              header: this.transloco.translate('global.error'),
              subHeader: this.transloco.translate(
                'pages.pet_invites.accept_shared_pet.request_failed_alert_subheader'
              ),
              message: err.message,
              buttons: [this.transloco.translate('global.ok')],
            })
            .then((alert) => alert.present());
          console.error('Error emitted');
          return of([]);
        }),
        finalize(() => console.log('Action to be executed always'))
      )
      .subscribe(() => this.modalCtrl.dismiss());
  }

  declineInvite(invite: PetShareInvite, inviteIndex: number) {
    this.api
      .denyInviteToSharedPet(invite.pet.uuid)
      .pipe(
        tap(() => console.log('Action performed before any other')),
        catchError((err) => {
          this.alertCtrl
            .create({
              header: this.transloco.translate('global.error'),
              subHeader: this.transloco.translate(
                'pages.pet_invites.deny_shared_pet.request_failed_alert_subheader'
              ),
              message: err.message,
              buttons: [this.transloco.translate('global.ok')],
            })
            .then((alert) => alert.present());
          console.error('Error emitted');
          return of([]);
        }),
        finalize(() => console.log('Action to be executed always'))
      )
      .subscribe(() => {
        if (this.invites) {
          delete this.invites[inviteIndex];
          if (this.invites.length == 0) {
            this.modalCtrl.dismiss();
          }
        }
      });
  }
}
