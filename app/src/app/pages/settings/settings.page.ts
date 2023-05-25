import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { IonicModule } from '@ionic/angular';
import { TranslocoRootModule } from 'src/app/transloco-root.module';
import { environment } from 'src/environments/environment';
import { TranslocoService } from '@ngneat/transloco';

type Language = Record<'imgUrl' | 'code' | 'name' | 'shorthand', string>;
@Component({
  selector: 'app-settings',
  templateUrl: 'settings.page.html',
  styleUrls: ['settings.page.scss'],
  standalone: true,
  imports: [IonicModule, CommonModule, TranslocoRootModule],
})
export class SettingsPage {
  usingSystemDarkTheme: boolean;
  selectedLanguage: Language;
  languagesList: Array<Language> = [
    {
      imgUrl: '/assets/images/English.png',
      code: 'en',
      name: 'English',
      shorthand: 'ENG',
    },
    {
      imgUrl: '/assets/images/Deutsch.png',
      code: 'de',
      name: 'German',
      shorthand: 'GER',
    },
  ];
  constructor(protected transloco: TranslocoService) {
    this.usingSystemDarkTheme = window.matchMedia(
      '(prefers-color-scheme: dark)'
    ).matches;
    this.selectedLanguage = {
      imgUrl: '/assets/images/Deutsch.png',
      code: 'de',
      name: 'German',
      shorthand: 'GER',
    };
  }

  toggleDarkTheme(event: any): void {
    document.body.classList.toggle('dark', event.detail.checked);
    this.usingSystemDarkTheme = event.detail.checked;
    console.log('to be implemented: toggle theme here');
  }

  handleLanguageChange(event: any): void {
    this.transloco.setActiveLang(languageCode);
    languageCode === 'fa'
      ? (document.body.style.direction = 'rtl')
      : (document.body.style.direction = 'ltr');
  }

  getAppVersion() {
    return environment.version;
  }
}
