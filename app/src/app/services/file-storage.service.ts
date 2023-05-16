import { Injectable } from '@angular/core';
import { AngularFireStorage } from '@angular/fire/compat/storage';

import { finalize } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class FileStorageService {
  constructor(private fireStorage: AngularFireStorage) {}

  uploadFile(
    userUid: string,
    petUuid: string,
    dataUrl: string
  ): Promise<string> {
    const filePath = `${userUid}/${petUuid}`;
    const fileRef = this.fireStorage.ref(filePath);
    const task = fileRef.putString(dataUrl, 'data_url');
    return new Promise((resolve) => {
      task
        .snapshotChanges()
        .pipe(
          finalize(() => {
            const downloadURL = fileRef.getDownloadURL();
            downloadURL.subscribe((url: string) => {
              if (url) {
                resolve(url);
              }
            });
          })
        )
        .subscribe(() => {});
    });
  }
}
