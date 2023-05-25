import { Injectable } from '@angular/core';
import { Camera, CameraResultType, CameraSource } from '@capacitor/camera';
import { AlertController } from '@ionic/angular';

@Injectable({
  providedIn: 'root',
})
export class PhotoService {
  constructor(private alertCtrl: AlertController) {}

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
        source: CameraSource.Camera,
        quality: 100,
      });
      return photo.dataUrl;
    } else {
      const alert = await this.alertCtrl.create({
        header: 'Error',
        subHeader: 'No permissions',
        message:
          'You need to grant permissions for photo album and camera, may check device settings, if you denied them permantly by mistake.',
        buttons: ['OK'],
      });
      await alert.present();
    }

    return;
  }
}
