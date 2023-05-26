import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { IonicModule } from '@ionic/angular';
import { TranslocoRootModule } from 'src/app/transloco-root.module';
import { environment } from 'src/environments/environment';
import { TranslocoService } from '@ngneat/transloco';

type Language = Record<'code' | 'name' | 'shorthand', string>;
@Component({
  selector: 'app-settings',
  templateUrl: 'settings.page.html',
  styleUrls: ['settings.page.scss'],
  standalone: true,
  imports: [IonicModule, CommonModule, TranslocoRootModule],
})
export class SettingsPage {
  usingSystemDarkTheme: boolean;
  selectedLanguageCode: string;
  languagesList: Array<Language> = [
    {
      code: 'en',
      name: 'English',
      shorthand: 'ENG',
    },
    {
      code: 'de',
      name: 'German',
      shorthand: 'GER',
    },
  ];

  constructor(protected transloco: TranslocoService) {
    this.usingSystemDarkTheme = window.matchMedia(
      '(prefers-color-scheme: dark)'
    ).matches;
    this.selectedLanguageCode = this.transloco.getActiveLang();
  }

  toggleDarkTheme(event: any): void {
    document.body.classList.toggle('dark', event.detail.checked);
    this.usingSystemDarkTheme = event.detail.checked;
    console.log('to be implemented: toggle theme here');
  }

  handleLanguageChange(event: any): void {
    const languageCode = event.detail.value;
    this.transloco.setActiveLang(languageCode);
    this.selectedLanguageCode = languageCode;
    localStorage.setItem('language', languageCode);
  }

  getAppVersion() {
    return environment.version;
  }
}
