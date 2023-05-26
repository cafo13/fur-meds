import { Injectable } from '@angular/core';
import { Camera, CameraResultType, CameraSource } from '@capacitor/camera';
import { AlertController } from '@ionic/angular';
import { TranslocoService } from '@ngneat/transloco';

@Injectable({
  providedIn: 'root',
})
export class PhotoService {
  constructor(
    private alertCtrl: AlertController,
    private transloco: TranslocoService
  ) {}

  public async getPhoto(): Promise<string | undefined> {
    const permissions = await Camera.checkPermissions();
    if (permissions.photos === 'denied') {
      await Camera.requestPermissions({
        permissions: ['photos'],
      });
    }
    if (permissions.camera === 'denied') {
      await Camera.requestPermissions({
        permissions: ['camera'],
      });
    }

    const newPermissions = await Camera.checkPermissions();

    if (
      newPermissions.camera !== 'denied' &&
      newPermissions.photos !== 'denied'
    ) {
      const photo = await Camera.getPhoto({
        resultType: CameraResultType.DataUrl,
        quality: 50,
        allowEditing: true,
        height: 200,
        promptLabelCancel: this.transloco.translate('global.cancel_button'),
        promptLabelPhoto: this.transloco.translate(
          'services.photo.select_from_gallery'
        ),
        promptLabelPicture: this.transloco.translate(
          'services.photo.take_picture'
        ),
      });
      return photo.dataUrl;
    } else {
      const alert = await this.alertCtrl.create({
        header: this.transloco.translate('global.error'),
        subHeader: this.transloco.translate(
          'services.photo.no_permissions_title'
        ),
        message: this.transloco.translate('services.photo.no_permissions_hint'),
        buttons: [this.transloco.translate('global.ok')],
      });
      await alert.present();
    }

    return;
  }
}
